package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/starkandwayne/goutils/log"

	"gopkg.in/yaml.v2"
)

//Version is the version string of this build of Portcullis
var Version = "(development build)"

//Config is an in-memory representation of the configuration file for Portcullis
type Config struct {
	Store    StoreConfig  `yaml:"store"`
	API      APIConfig    `yaml:"api"`
	Broker   BrokerConfig `yaml:"broker"`
	LogLevel string       `yaml:"log_level"`
}

//Load creates and fills out a Config struct from the configuration file
//found at the given path. An error is returned if the configuration file
//is not parseable as YAML, or if the file cannot be found.
func Load(path string) (c Config, err error) {
	log.Infof("Loading configuration file")
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Unable to open configuration file: %s", err.Error())
		return
	}
	err = yaml.Unmarshal(buffer, &c)
	if err != nil {
		err = fmt.Errorf("Error while parsing configuration file: %s", err.Error())
	}

	c.setDefaults()
	c.LogLevel = strings.ToLower(c.LogLevel)

	err = c.verifyBaseConfig()
	return
}

func (c *Config) setDefaults() {
	const defaultAPIDescription = "Portcullis API"
	const defaultLogLevel = "info"

	if c.API.Description == "" {
		log.Infof("Setting API Description to default: %s", defaultAPIDescription)
		c.API.Description = defaultAPIDescription
	}

	if c.LogLevel == "" {
		log.Infof("Setting Log Level to default: %s", defaultLogLevel)
		c.LogLevel = defaultLogLevel
	}
}

func (c *Config) verifyBaseConfig() error {
	if !(c.LogLevel == "debug" || c.LogLevel == "info" || c.LogLevel == "warn" ||
		c.LogLevel == "error" || c.LogLevel == "crit" || c.LogLevel == "off") {
		return fmt.Errorf("Configured log level was not understood: %s", c.LogLevel)
	}

	if c.LogLevel == "off" {
		//We don't use "emerg" level in the program, so this squelches everything
		c.LogLevel = "emerg"
	}
	return nil
}
