package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	watchFile   = "/configmap/marker"
	rehupSignal = syscall.SIGHUP

	spacer = "\n\n****************************************************\n\n"
)

func main() {

	log.Println("Starting Watcher...")
	info, err := os.Stat(watchFile)
	if err != nil {
		log.Fatal("ERROR Reading File: ", err)
	}
	ts := info.ModTime()

	for {
		time.Sleep(2 * time.Second)
		info, err = os.Stat(watchFile)
		if err != nil {
			log.Println("ERROR Reading File: ", err)
			continue
		}

		ts2 := info.ModTime()

		if ts2 != ts {
			log.Println(spacer + "Config file updated, reloading Envoy" + spacer)
			ts = ts2

			// Copy over new files
			shell("/reload.sh")
		}
	}
}

func shell(params ...string) {

	cmd := params[0]
	args := params[1:]

	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command(cmd, args...).CombinedOutput(); err != nil {
		log.Println("ERROR: ", err)
	}
	log.Println(string(cmdOut))
}
