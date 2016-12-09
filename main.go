package main

import (
	"os"

	"github.com/cloudfoundry-community/portcullis/api"
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
		bailWith("Please define the configuration file location with the environment variable `PORTCULLIS_CONFIG`")
	}

	conf, err := config.Load(configPath)
	if err != nil {
		bailWith("Error while loading config: %s", err)
	}

	err = store.SetStoreType(conf.Store.Type)
	if err != nil {
		bailWith("Error while setting store type: %s", err)
	}
	err = store.Initialize(conf.Store.Config)
	if err != nil {
		bailWith("Error while initializing store: %s", err)
	}
	err = api.Initialize(conf.API)
	if err != nil {
		bailWith("Error while initializing API server: %s", err)
	}
	apiChan := make(chan error)
	go api.Launch(apiChan)
	select {
	case err := <-apiChan:
		bailWith("API Server closed with error: %s", err)
	}
}

func bailWith(mess string, args ...interface{}) {
	log.Critf(mess, args...)
	os.Exit(1)
}
