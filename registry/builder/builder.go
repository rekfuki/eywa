package builder

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"eywa/registry/clients/mongo"
	"eywa/registry/types"
)

var (
	// ErrUnsupportedLanguage ...
	ErrUnsupportedLanguage = fmt.Errorf("Unsupported language")
	templateLocations      = map[string]string{"go": "./templates/go"}
	supportedLanguages     = []string{"go"}
)

// Config represents builder config
type Config struct {
	Registry         string
	RegistryUser     string
	RegistryPassword string
	Mongo            *mongo.Client
	NumWorkers       int
}

// Client represents the builder client
type Client struct {
	registry       string
	mc             *mongo.Client
	incomingBuilds chan types.BuildRequest
	numWorkers     int
}

// New creates a new builder client
func New(conf *Config) (*Client, error) {
	registryURL := fmt.Sprintf("https://%s", conf.Registry)
	if err := authenticate(registryURL, conf.RegistryUser, conf.RegistryPassword); err != nil {
		return nil, err
	}

	return &Client{
		registry:       conf.Registry,
		mc:             conf.Mongo,
		incomingBuilds: make(chan types.BuildRequest),
		numWorkers:     conf.NumWorkers,
	}, nil
}

// Start starts the builder
func (c *Client) Start() {
	for i := 0; i < c.numWorkers; i++ {
		go func() {
			for br := range c.incomingBuilds {
				state := types.StateSuccess
				if err := c.build(br); err != nil {
					state = types.StateFailed
					log.WithFields(log.Fields{
						"ID":       br.ID,
						"Language": br.Language,
						"Version":  br.Version,
					}).Error("Failed to build container image")
				}

				if err := c.mc.UpdateImageState(br.ID, state); err != nil {
					log.Errorf("Failed to update source state: %s", err)
				}
			}
		}()
	}
}

// Enqueue queues up a new build request
func (c *Client) Enqueue(br types.BuildRequest) string {
	c.incomingBuilds <- br
	return fmt.Sprintf("%s/%s:%s", c.registry, br.ID, br.Version)
}

// Build builds the container based on the language
func (c *Client) build(br types.BuildRequest) error {
	if _, ok := templateLocations[br.Language]; !ok {
		return ErrUnsupportedLanguage
	}

	baseDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(baseDir)

	if err := extractSource(baseDir, br.ZipReader); err != nil {
		return err
	}

	switch br.Language {
	case "go":
		prepareGo(baseDir)
	}

	image := fmt.Sprintf("%s/%s:%s", c.registry, br.ID, br.Version)
	if err := buildImage(baseDir, image); err != nil {
		return err
	}

	if err := pushImage(image); err != nil {
		return err
	}

	return nil
}

func extractSource(baseDir string, r *zip.Reader) error {
	if err := os.Mkdir(filepath.Join(baseDir, "source"), os.ModePerm); err != nil {
		return err
	}

	for _, file := range r.File {
		// Since we are using MkdirAll, there is no need to worry about dirs
		// Empty directories are not needed anyways
		if file.FileInfo().IsDir() {
			continue
		}

		fPath := filepath.Join(baseDir, "source", file.Name)

		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func authenticate(registryURL, user, password string) error {
	authCMD := exec.Command("img", "login", "-u", user, "-p", password, registryURL)
	return execCommand(authCMD)
}

func buildImage(buildDir, image string) error {
	buildCMD := exec.Command("img", "build", "-t", image, buildDir)
	return execCommand(buildCMD)
}

func pushImage(image string) error {
	pushCMD := exec.Command("img", "push", image)
	return execCommand(pushCMD)
}

func execCommand(command *exec.Cmd) error {

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	if err := command.Start(); err != nil {
		return err
	}

	// We probably want to log something here but no the whole process
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := command.Wait(); err != nil {
		return err
	}

	return nil
}
