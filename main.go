package main

import (
	"os"

	"github.com/cloudfoundry-community/portcullis/config"
	"github.com/cloudfoundry-community/portcullis/store"
	"github.com/starkandwayne/goutils/log"

	_ "github.com/cloudfoundry-community/portcullis/store/dummy"
	_ "github.com/cloudfoundry-community/portcullis/store/postgres"
)

func main() {
	log.SetupLogging(log.LogConfig{
		Type:  "console",
		Level: "debug",
	})
	configPath := os.Getenv("PORTCULLIS_CONFIG")

	if configPath == "" {
		log.Warnf("I'm a little %s", "teapot")
		bailWith("Please define the configuration file location with the environment variable `PORTCULLIS_CONFIG`")
	}

	_, err := config.Load(configPath)
	if err != nil {
		bailWith("Error while loading config:\n%s", err)
	}

	store.SetStoreType("dummy")
}

func bailWith(mess string, args ...interface{}) {
	log.Critf(mess, args...)
	os.Exit(1)
}
