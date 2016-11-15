package store

var mystore store

type store interface {
	connect(port int, location, username, password, dbname string) error
	getMappings() []Mapping
	addMapping(Mapping)
	deleteMapping(name string)
}
