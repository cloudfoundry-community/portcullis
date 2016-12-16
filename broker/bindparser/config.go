package bindparser

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

//Config is a configuration object used to set up a bindparser Flavor
type Config struct {
	FlavorName string                 `json:"flavor",yaml:"flavor"`
	Config     map[string]interface{} `json:"config",yaml:"config"`
}

var flavorMap = map[string]flavorMaker{
	"dummy": NewDummy,
}

//CreateFlavor creates the implementation of flavor as specified by the config
// object, and unmarshals the map into the struct where applicable
func (c Config) CreateFlavor() (f Flavor, e error) {
	flavorFunc, found := flavorMap[c.FlavorName]
	if !found {
		return nil, fmt.Errorf("There is no known flavor type with name `%s`", c.FlavorName)
	}
	f = flavorFunc()
	err := configIntoFlavor(c.Config, f)
	if err != nil {
		return nil, err
	}
	return
}

//VerifyFlavor is just a shorthand for CreateFlavor, and then calling Verify on
// the produced Flavor instance
func (c Config) VerifyFlavor() error {
	flavor, err := c.CreateFlavor()
	if err != nil {
		return err
	}
	err = flavor.Verify()
	return err
}

func configIntoFlavor(conf map[string]interface{}, flavor Flavor) error {
	j, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("Could not marshal bind config into JSON")
	}
	err = yaml.Unmarshal(j, flavor)
	return err
}
