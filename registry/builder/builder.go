package builder

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"eywa/registry/db"
	"eywa/registry/types"
)

var (
	templateLocations = map[string]string{
		"go":      "./templates/go",
		"node14":  "./templates/node14",
		"python3": "./templates/python3",
		"ruby":    "./templates/ruby",
		"csharp":  "./templates/csharp",
		"custom":  "./templates/custom",
	}
	requiredFiles = map[string][]string{
		"go":      {"handler.go"},
		"node14":  {"handler.js", "package.json"},
		"python3": {"handler.py", "requirements.txt"},
		"ruby":    {"handler.rb", "Gemfile"},
		"csharp":  {"Function.csproj", "FunctionHandler.cs"},
	}
)

// Config represents builder config
type Config struct {
	Registry         string
	RegistryUser     string
	RegistryPassword string
	DB               *db.Client
	NumWorkers       int
}

// Client represents the builder client
type Client struct {
	registry         string
	db               *db.Client
	incomingBuilds   chan BuildRequest
	numWorkers       int
	inProgressBuilds map[string]BuildRequest
}

// New creates a new builder client
func New(conf *Config) (*Client, error) {
	registryURL := fmt.Sprintf("https://%s", conf.Registry)
	if err := authenticate(registryURL, conf.RegistryUser, conf.RegistryPassword); err != nil {
		return nil, err
	}

	return &Client{
		registry:         conf.Registry,
		db:               conf.DB,
		incomingBuilds:   make(chan BuildRequest),
		numWorkers:       conf.NumWorkers,
		inProgressBuilds: make(map[string]BuildRequest),
	}, nil
}

// Start starts the builder
func (c *Client) Start() {
	for i := 0; i < c.numWorkers; i++ {
		go func() {
			for br := range c.incomingBuilds {
				logs := []string{}

				state := StateSuccess
				if err := c.db.UpdateImageState(br.ImageID, StateBuilding); err != nil {
					log.Errorf("Failed to update image state: %s", err)
					logs = append(logs, BuildSystemErrorMessage(err.Error()))
					state = StateFailed
				} else {
					buildErr := c.build(br)

					logFile, err := os.OpenFile(br.LogFile, os.O_APPEND|os.O_RDWR, 0666)
					if err != nil {
						logs = append(logs, BuildSystemErrorMessage(err.Error()))
					} else {
						if buildErr != nil {
							log.WithFields(log.Fields{
								"Image ID": br.ImageID,
								"Runtime":  br.Runtime,
								"Version":  br.Version,
							}).Errorf("Failed to build container image", buildErr)

							state = StateFailed
							errMsg := ""
							if buildErr.Type == ErrTypeSystemError {
								errMsg = BuildSystemErrorMessage(ErrInternalError.String())
							} else if buildErr.Type == ErrTypeUserError {
								errMsg = BuildUserErrorMessage(buildErr.String())
							} else if buildErr.Type == ErrTypeBuild {
								errMsg = BuildErrorMessage(buildErr.String())
							}

							_, err = logFile.WriteString(errMsg + BuildFailedMessage())
							if err != nil {
								logs = append(logs, BuildSystemErrorMessage(err.Error()))
							}
						}

						logFile.Seek(0, io.SeekStart)

						scanner := bufio.NewScanner(logFile)
						scanner.Split(bufio.ScanLines)
						for scanner.Scan() {
							logs = append(logs, scanner.Text())
						}
					}
				}
				if err := c.db.UpdateBuild(br.ImageID, state, logs); err != nil {
					log.Errorf("Failed to update build: %s", err)
				}

				if err := c.db.UpdateImageState(br.ImageID, state); err != nil {
					log.Errorf("Failed to update image state: %s", err)
				}

				delete(c.inProgressBuilds, br.ImageID+br.UserID)

				if br.tmpDir != "" {
					os.RemoveAll(br.tmpDir)
				}
			}
		}()
	}
}

// GetBuild returns a build in progress
func (c *Client) GetBuild(buildID, userID string) *BuildRequest {
	if build, exists := c.inProgressBuilds[buildID+userID]; exists {
		return &build
	}
	return nil
}

// Enqueue queues up a new build request
func (c *Client) Enqueue(br BuildRequest) *Error {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return SystemError(err.Error())
	}

	logFilePath := filepath.Join(tmpDir, "logs.txt")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		return SystemError(err.Error())
	}
	defer logFile.Close()

	_, err = logFile.WriteString(BuildQueuedMessage(br.ImageID, br.Runtime, br.Version))
	if err != nil {
		return SystemError(err.Error())
	}

	br.tmpDir = tmpDir
	br.LogFile = logFilePath

	if br.Runtime == "custom" {
		if br.ExecutablePath == nil {
			return UserError("Executable path is required when using custom runner")
		}
		br.requiredFiles = []string{*br.ExecutablePath}
	} else {
		br.requiredFiles = requiredFiles[br.Runtime]
	}

	c.inProgressBuilds[br.ImageID+br.UserID] = br

	build := &types.Build{
		ImageID:   br.ImageID,
		UserID:    br.UserID,
		State:     StateQueued,
		Logs:      pq.StringArray{},
		CreatedAt: time.Now(),
	}

	if err := c.db.CreateBuild(build); err != nil {
		return SystemError(err.Error())
	}

	taggedRegistry := fmt.Sprintf("%s/%s:%s", c.registry, br.ImageID, br.Version)

	image := &types.Image{
		ID:             br.ImageID,
		UserID:         br.UserID,
		Runtime:        br.Runtime,
		Name:           br.Name,
		Version:        br.Version,
		State:          StateQueued,
		TaggedRegistry: taggedRegistry,
		CreatedAt:      time.Now(),
		Size:           len(br.ZippedSource),
		Source:         base64.StdEncoding.EncodeToString(br.ZippedSource), // Only store if successfull ?
	}

	if err := c.db.CreateImage(image); err != nil {
		log.Errorf("Failed to create image in db: %s", err)
		return SystemError(err.Error())
	}

	c.incomingBuilds <- br

	return nil
}

// Build builds the container based on the language
func (c *Client) build(br BuildRequest) *Error {
	logFile, err := os.OpenFile(br.LogFile, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return SystemError(err.Error())
	}
	defer logFile.Close()

	_, err = logFile.WriteString(BuildStartMessage())
	if err != nil {
		return SystemError(err.Error())
	}

	if _, ok := templateLocations[br.Runtime]; !ok {
		return ErrUnsupportedLanguage
	}

	if err := extractSource(br.tmpDir, br.ZippedSource); err != nil {
		return err
	}

	buildErr := prepareLanguage(br.tmpDir, br.Runtime, br.requiredFiles...)
	if buildErr != nil {
		return buildErr
	}

	image := fmt.Sprintf("%s/%s:%s", c.registry, br.ImageID, br.Version)

	if br.Runtime == "custom" {
		err = buildImage(br.tmpDir, image, logFile, "FPROCESS="+*br.ExecutablePath)
	} else {
		err = buildImage(br.tmpDir, image, logFile)
	}
	if err != nil {
		return ImageProcessError(err.Error())
	}

	if err := pushImage(image, logFile); err != nil {
		return ImageProcessError(err.Error())
	}

	if err := cleanup(); err != nil {
		return SystemError(err.Error())
	}

	_, err = logFile.WriteString(BuildSuccessMessage())
	if err != nil {
		return SystemError(err.Error())
	}

	return nil
}

func extractSource(baseDir string, source []byte) *Error {
	if err := os.Mkdir(filepath.Join(baseDir, "source"), os.ModePerm); err != nil {
		return SystemError(err.Error())
	}

	r, err := zip.NewReader(bytes.NewReader(source), int64(len(source)))
	if err != nil {
		if err == zip.ErrAlgorithm || err == zip.ErrFormat || err == zip.ErrChecksum {
			return UserError(err.Error())
		}

		return SystemError(err.Error())
	}

	for _, file := range r.File {
		// Since we are using MkdirAll, there is no need to worry about dirs
		// Empty directories are not needed anyways
		if file.FileInfo().IsDir() {
			continue
		}

		fPath := filepath.Join(baseDir, "source", file.Name)

		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return SystemError(err.Error())
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return SystemError(err.Error())
		}

		rc, err := file.Open()
		if err != nil {
			return SystemError(err.Error())
		}

		_, err = io.Copy(outFile, rc)
		if err := outFile.Close(); err != nil {
			return SystemError(err.Error())
		}

		if err := rc.Close(); err != nil {
			return SystemError(err.Error())
		}

		if err != nil {
			return SystemError(err.Error())
		}
	}

	return nil
}

func prepareLanguage(buildDir, language string, files ...string) *Error {
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(buildDir, "source", file)); os.IsNotExist(err) {
			return UserError(fmt.Sprintf("Missing top level %s entry file", file))
		}
	}

	if err := os.Rename(filepath.Join(buildDir, "source"), filepath.Join(buildDir, "function")); err != nil {
		return SystemError(err.Error())
	}

	templateFiles, err := ioutil.ReadDir(templateLocations[language])
	if err != nil {
		return SystemError(err.Error())
	}

	for _, file := range templateFiles {
		from, err := os.Open(filepath.Join(templateLocations[language], file.Name()))
		if err != nil {
			return SystemError(err.Error())
		}
		defer from.Close()

		to, err := os.OpenFile(filepath.Join(buildDir, file.Name()), os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return SystemError(err.Error())
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			return SystemError(err.Error())
		}
	}

	return nil
}

func authenticate(registryURL, user, password string) error {
	authCMD := exec.Command("img", "login", "-u", user, "-p", password, registryURL)
	return execCommand(authCMD, nil)
}

func buildImage(buildDir, image string, logFile *os.File, args ...string) error {
	buildArgs := []string{"build", "-t", image}
	for _, arg := range args {
		buildArgs = append(buildArgs, []string{"--build-arg", arg}...)
	}
	buildArgs = append(buildArgs, buildDir)
	buildCMD := exec.Command("img", buildArgs...)
	return execCommand(buildCMD, logFile)
}

func pushImage(image string, logFile *os.File) error {
	pushCMD := exec.Command("img", "push", image)
	return execCommand(pushCMD, logFile)
}

func cleanup() error {
	pushCMD := exec.Command("img", "prune")
	return execCommand(pushCMD, nil)
}

func execCommand(command *exec.Cmd, logFile *os.File) error {
	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	if err := command.Start(); err != nil {
		return err
	}

	if logFile != nil {
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			_, err := logFile.WriteString(scanner.Text() + "\n")
			if err != nil {
				return err
			}
		}
	}

	if err := command.Wait(); err != nil {
		return err
	}

	return nil
}
