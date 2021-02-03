package executor

import (
	"bufio"
	"io"

	log "github.com/sirupsen/logrus"
)

// bindLoggingPipe spawns a goroutine for passing through logging of the given output pipe.
func bindLoggingPipe(name string, pipe io.Reader, output *[]string) {
	log.Infof("Started logging %s from function.", name)

	scanner := bufio.NewScanner(pipe)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			*output = append(*output, text)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error scanning %s: %s", name, err.Error())
		}
	}()
}
