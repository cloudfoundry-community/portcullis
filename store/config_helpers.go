package store

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

//ErrIfMissingKeys returns a nil error if all given keys are present in `config, and
// returns an error with an appropriate message if it encounters missing keys
func ErrIfMissingKeys(config map[string]interface{}, keys ...string) (err error) {
	var messages = []string{}
	for _, k := range keys {
		_, found := config[k]
		if !found {
			messages = append(messages, fmt.Sprintf("Unable to locate key `%s` in store config", k))
		}
	}
	if len(messages) > 0 {
		err = fmt.Errorf(strings.Join(messages, "\n"))
	}
	return err
}

func ParseStoreConfig(conf map[string]interface{}, confStruct interface{}) error {
	yamlConf, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("Could not read store config")
	}
	err = yaml.Unmarshal(yamlConf, confStruct)
	if err != nil {
		return fmt.Errorf("Something nasty happened while parsing store config")
	}
	return nil
}
