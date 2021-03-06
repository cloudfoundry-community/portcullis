package store

import "encoding/json"
import "github.com/cloudfoundry-community/portcullis/broker/bindparser"

//Mapping represents a mapping between a service broker name and a service
//broker backend, as well as the configuration details of how to work with it
type Mapping struct {
	Name       string            `json:"name"`
	Location   string            `json:"location"`
	BindConfig bindparser.Config `json:"bind_config"`
}

//MappingFields is an array of all the top-level fields in a JSON object
// representing a Mapping that are understood by the program
var MappingFields = [3]string{"name", "location", "bind_config"}

//RequiredMappingFields is an array of all the top-level fields in a JSON object
// representing a Mapping in order for there to be enough information for the
// program to use the Mapping
var RequiredMappingFields = [3]string{"name", "location", "bind_config"}

//WithName generates a new Mapping with all the properties of the target Mapping,
// except with the given name
func (m Mapping) WithName(name string) Mapping {
	ret := m
	ret.Name = name
	return ret
}

//WithConfig generates a new Mapping with all the properties of the target Mapping,
// except with the given config
func (m Mapping) WithConfig(config bindparser.Config) Mapping {
	ret := m
	ret.BindConfig = config
	return ret
}

//ToMap converts this Mapping object into a map[string]interface{}, with keys as
// specified by the JSON interface
func (m Mapping) ToMap() (ret map[string]interface{}, err error) {
	var j []byte
	j, err = json.Marshal(m)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &ret)
	return
}

//MappingFromMap creates a Mapping object from a map[string]interface{} to a Mapping object,
// using the keys as expected by the JSON interface
func MappingFromMap(m map[string]interface{}) (ret Mapping, err error) {
	var j []byte
	j, err = json.Marshal(m)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &ret)
	return
}

//MappingList is an array of Mapping objects, named so that it may implement sort.Interface
type MappingList []Mapping

//Len returns the length of the array
func (m MappingList) Len() int { return len(m) }

//Swap reverses the places of the objects at indexes i and j
func (m MappingList) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

//Less returns true if the Name of Mapping at index i is lexically earlier than
// the name at j
func (m MappingList) Less(i, j int) bool { return m[i].Name < m[j].Name }
