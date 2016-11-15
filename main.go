package main

import (
	"os"

	"github.com/cloudfoundry-community/portcullis/config"
	"github.com/cloudfoundry-community/portcullis/log"
)

func main() {
	configPath := os.Getenv("PORTCULLIS_CONFIG")

	if configPath == "" {
		bailWithF("Please define the configuration file location with the environment variable `PORTCULLIS_CONFIG`")
	}

	_, err := config.Load(configPath)
	if err != nil {
		bailWithF("Error while loading config:\n%s", err)
	}
}

func bailWithF(mess string, fmtargs ...interface{}) {
	log.Errorf(mess, fmtargs...)
	os.Exit(1)
}
