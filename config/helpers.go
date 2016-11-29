package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	//StoreKey means error messages should be tailored to store configs
	StoreKey = "store"
	//AuthKey means error messages should be tailored to auth configs
	AuthKey = "auth"
)

//ErrIfMissingKeys returns a nil error if all given keys are present in `config, and
// returns an error with an appropriate message if it encounters missing keys
// confkey is
func ErrIfMissingKeys(confkey string, config map[string]interface{}, keys ...string) (err error) {
	validateConfKey(confkey)

	var messages = []string{}
	for _, k := range keys {
		_, found := config[k]
		if !found {
			messages = append(messages, fmt.Sprintf("Unable to locate key `%s.config.%s` in store config", confkey, k))
		}
	}
	if len(messages) > 0 {
		err = fmt.Errorf(strings.Join(messages, "\n"))
	}
	return err
}

//ParseMapConfig takes the config for a store and unmarshals it into the
// given interface. Useful for taking the arbitrary map in the default
// configuration object and making it into a store-specific struct.
func ParseMapConfig(confkey string, conf map[string]interface{}, confStruct interface{}) error {
	validateConfKey(confkey)

	yamlConf, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("Could not read %s config", confkey)
	}
	err = yaml.Unmarshal(yamlConf, confStruct)
	if err != nil {
		return fmt.Errorf("Something nasty happened while parsing %s config", confkey)
	}
	return nil
}

func validateConfKey(key string) {
	if key != StoreKey && key != AuthKey {
		panic("Unrecognized configuration type!!!!!!")
	}
}
