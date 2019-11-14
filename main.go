package main

import (
	"log"
	"os"
	"time"

	"github.com/pivotal-cf/om/commands"
	_ "github.com/pivotal-cf/om/download_clients"
)

var version = "unknown"
var applySleepDurationString = "10s"

func main() {
	applySleepDuration, _ := time.ParseDuration(applySleepDurationString)

	command := commands.NewMain(
		os.Stdout,
		os.Stderr,
		os.Stdin,
		version,
		applySleepDuration,
	)

	err := command.Execute(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
}
