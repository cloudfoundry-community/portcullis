package config

import (
	"fmt"
	"strings"

	"github.com/starkandwayne/goutils/log"

	"gopkg.in/yaml.v2"
)

const (
	//StoreKey means error messages should be tailored to store configs
	StoreKey = "store"
	//AuthKey means error messages should be tailored to auth configs
	AuthKey = "auth"
	//TestKey is the key used for test cases
	TestKey = "test"
)

//ValidateConfigKeys logs to Warnf if there are keys other than those provided
// to `keys`. It logs to Errorf and returns an error if there are keys provided
// to `keys` that were not found in `config`
func ValidateConfigKeys(confkey string, config map[string]interface{}, keys ...string) error {
	if err := errIfExtraKeys(confkey, config, keys...); err != nil {
		log.Warnf(err.Error())
	}
	err := errIfMissingKeys(confkey, config, keys...)
	if err != nil {
		log.Errorf(err.Error())
	}
	return err
}

//errIfMissingKeys returns a nil error if all given keys are present in `config, and
// returns an error with an appropriate message if it encounters missing keys
// confkey is
func errIfMissingKeys(confkey string, config map[string]interface{}, keys ...string) (err error) {
	validateConfKey(confkey)

	var messages = []string{}
	for _, k := range keys {
		_, found := config[k]
		if !found {
			messages = append(messages, fmt.Sprintf("Unable to locate key `%s.config.%s` in config", confkey, k))
		}
	}
	return errorFromMessages(messages)
}

//errIfExtraKeys returns an error with an appropriate message if there are keys
// other than those specified to this function found within config
func errIfExtraKeys(confkey string, config map[string]interface{}, keys ...string) (err error) {
	validateConfKey(confkey)

	var messages = []string{}
	for k := range config {
		for _, k2 := range keys {
			if k == k2 {
				goto notExtra
			}
		}
		messages = append(messages, fmt.Sprintf("Extraneous config key found at `%s.config.%s`", confkey, k))
	notExtra:
	}
	return errorFromMessages(messages)
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
	if key != StoreKey && key != AuthKey && key != TestKey {
		panic("Unrecognized configuration type!!!!!!")
	}
}

func errorFromMessages(messages []string) (err error) {
	if len(messages) > 0 {
		err = fmt.Errorf(strings.Join(messages, "\n"))
	}
	return err
}
