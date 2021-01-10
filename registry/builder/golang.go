package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func prepareGo(buildDir string) error {
	// Check for top level handler.go entry file
	if _, err := os.Stat(filepath.Join(buildDir, "source", "handler.go")); os.IsNotExist(err) {
		return fmt.Errorf("Missing top level handler.go entry file")
	}

	if err := os.Rename(filepath.Join(buildDir, "source"), filepath.Join(buildDir, "function")); err != nil {
		return err
	}

	templateFiles, err := ioutil.ReadDir(templateLocations["go"])
	if err != nil {
		return err
	}

	for _, file := range templateFiles {
		err := os.Link(filepath.Join(templateLocations["go"], file.Name()), filepath.Join(buildDir, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
