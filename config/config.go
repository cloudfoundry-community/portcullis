package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//Config is an in-memory representation of the configuration file for Portcullis
type Config struct {
	Store StoreConfig `yaml:"store"`
}

//Load creates and fills out a Config struct from the configuration file
//found at the given path. An error is returned if the configuration file
//is not parseable as YAML, or if the file cannot be found.
func Load(path string) (c Config, err error) {
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
