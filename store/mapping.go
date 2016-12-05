package store

//Mapping represents a mapping between a service broker name and a service
//broker backend, as well as the configuration details of how to work with it
type Mapping struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	//TODO
}

//WithName generates a new Mapping with all the properties of the target Mapping,
// except with the given name
func (m Mapping) WithName(name string) Mapping {
	ret := m
	ret.Name = name
	return ret
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
