package main

import (
	"os"

	"github.com/cloudfoundry-community/portcullis/api"
	"github.com/cloudfoundry-community/portcullis/broker"
	"github.com/cloudfoundry-community/portcullis/config"
	"github.com/cloudfoundry-community/portcullis/store"
	"github.com/starkandwayne/goutils/log"

	_ "github.com/cloudfoundry-community/portcullis/store/dummy"
	_ "github.com/cloudfoundry-community/portcullis/store/postgres"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cmdLine        = kingpin.New("portcullis", "A server which makes managing your CF service brokers easier").Version("portcullis " + config.Version)
	configPath     = cmdLine.Flag("config", "The path to the configuration file").Short('c').Default(os.Getenv("PORTCULLIS_CONFIG")).PlaceHolder("/path/to/config").String()
	skipBrokerFlag = cmdLine.Flag("test-without-broker", "Skip starting the broker server").Hidden().Bool()
)

func main() {
	cmdLine.HelpFlag.Short('h')
	cmdLine.VersionFlag.Short('v')
	command := kingpin.MustParse(cmdLine.Parse(os.Args[1:]))
	switch command {
	case "":
		initializePortcullis()
	default:
		bailWith("Unrecognized command: %s", command)
	}
}

func initializePortcullis() {
	//Need a default logging endpoint if the program needs to log before the config
	// can be loaded
	log.SetupLogging(log.LogConfig{
		Type:  "console",
		Level: "debug",
	})

	if *configPath == "" {
		bailWith("Please define the configuration file location with the `-c` flag")
	}

	conf, err := config.Load(*configPath)
	if err != nil {
		bailWith("Error while loading config: %s", err)
	}

	log.SetupLogging(log.LogConfig{
		Type:  "console",
		Level: conf.LogLevel,
	})

	log.Debugf("Logging settings configured")

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

	if !*skipBrokerFlag {
		err = broker.Initialize(conf.Broker)
		if err != nil {
			bailWith("Error while initializing Broker server: %s", err)
		}
	} else {
		log.Infof("Skipping broker initialization")
	}

	apiChan := make(chan error)
	go api.Launch(apiChan)

	//Launch the broker if we should
	brokerChan := make(chan error)
	if !*skipBrokerFlag {
		go broker.Launch(brokerChan)
	} else {
		log.Infof("Skipping broker launch")
	}

	select {
	case err := <-apiChan:
		bailWith("API Server closed with error: %s", err)
	case err := <-brokerChan:
		bailWith("Broker Server closed with error: %s", err)
	}
}

func bailWith(mess string, args ...interface{}) {
	log.Critf(mess, args...)
	os.Exit(1)
}
