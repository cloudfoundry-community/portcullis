package config

import (
	"fmt"
	"io/ioutil"

	"github.com/starkandwayne/goutils/log"

	"gopkg.in/yaml.v2"
)

//Version is the version string of this build of Portcullis
var Version = "(development build)"

//Config is an in-memory representation of the configuration file for Portcullis
type Config struct {
	Store  StoreConfig  `yaml:"store"`
	API    APIConfig    `yaml:"api"`
	Broker BrokerConfig `yaml:"broker"`
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
	return
}

//^
func (c Config) setDefaults() {
	//TODO
}
