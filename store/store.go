package store

import (
	"fmt"

	"github.com/starkandwayne/goutils/log"
)

var mystore Store

//Store is an interface which represents a targetable storage backend
type Store interface {
	//Initialize should set up the database to store data with the current version
	// of Portcullis. This means establishing a connection if necessary, and
	// performing necessary schema migrations.
	//The config is given in the form of a map to maintain a robust interface.
	// Most store implementation will likely expect a "port", "location", "username",
	// "password", and "dbname".
	Initialize(map[string]interface{}) error
	//ListMappings should return all Mappings that are currently in the store
	ListMappings() ([]Mapping, error)
	//GetMapping should return the mapping with the given name, and return
	//ErrNotFound if there is no mapping with that name in the store
	GetMapping(string) (Mapping, error)
	//AddMapping should put a new mapping into the store, and return ErrDuplicate if
	//a mapping with that name already exists in the store
	AddMapping(Mapping) error
	//EditMapping should edit the mapping with the provided name to
	//have all the values in the given Mapping. Should return ErrNotFound if there
	//is no mapping in the store with the name in the given Mapping
	EditMapping(string, Mapping) error
	//DeleteMapping should remove an existing mapping from the store, and return
	//ErrNotFound if the Mapping to remove did not exist in the store
	DeleteMapping(string) error
	//Size should return the number of mappings in the store
	Size() (int, error)
	//ClearMappings should delete all mappings from the database. A user should not need
	//to reinitialize the database, but ListMappings should return an empty list
	//after a call to ClearMappings
	ClearMappings() error
}

var (
	storeTypes  = map[string]Store{}
	activeStore Store
)

//SetStoreType sets the active store of the store library to the variant
// referenced by the given string.
//
//Current types are "dummy", "postgres"
func SetStoreType(variant string) (err error) {
	log.Infof("Setting store type to %s", variant)
	var found bool
	activeStore, found = storeTypes[variant]
	if !found {
		errorString := fmt.Sprintf("No store exists with variant name `%s`", variant)
		log.Errorf(errorString)
		err = fmt.Errorf(errorString)
	}
	return err
}

//Initialize configures the database to store data with the current version
// of Portcullis. This means establishing a connection if necessary, and
// performing necessary schema migrations.
//The config is given in the form of a map to maintain a robust interface.
// Most store implementation will likely expect a "port", "location", "username",
// "password", and "dbname".
func Initialize(config map[string]interface{}) error {
	log.Infof("Initializing store")
	return activeStore.Initialize(config)
}

//ListMappings returns all Mappings that are currently in the store
func ListMappings() (m []Mapping, err error) {
	return activeStore.ListMappings()
}

//GetMapping returns the mapping with the given name, and return ErrNotFound if
// there is no mapping with that name in the store
func GetMapping(name string) (Mapping, error) {
	return activeStore.GetMapping(name)
}

//AddMapping puts a new mapping into the store, and return ErrDuplicate if a
// mapping with that name already exists in the store
func AddMapping(m Mapping) error {
	//TODO: Create and enforce restrictions on mapping fields
	//  Make sure name is proper length/content
	//  Make sure location is parseable as a URL

	err := m.BindConfig.VerifyFlavor()
	if err != nil {
		return NewErrInvalid(err.Error())
	}
	return activeStore.AddMapping(m)
}

//EditMapping edits the mapping with the name in the given Mapping to
//have all the values in the given Mapping. Should return ErrNotFound if there
//is no mapping in the store with the name in the given Mapping. Should return
//ErrDuplicate if the name is being edited, and the name to edit to already
//exists in the store.
func EditMapping(name string, m Mapping) error {
	//TODO: See restriction checking for AddMapping
	err := m.BindConfig.VerifyFlavor()
	if err != nil {
		return NewErrInvalid(err.Error())
	}

	return activeStore.EditMapping(name, m)
}

//DeleteMapping removes an existing mapping from the store, and return
//ErrNotFound if the Mapping to remove did not exist in the store
func DeleteMapping(name string) error {
	return activeStore.DeleteMapping(name)
}

//ClearMappings deletes all existing mappings from the store. Mostly here for
//making testing more reliable
func ClearMappings() error {
	return activeStore.ClearMappings()
}

//Size returns the number of mappings in the store
func Size() (int, error) {
	return activeStore.Size()
}

//RegisterStoreType maps a type of store to a string that names it, such that
// this package can attach to a type of store that is indicated by a user string
func RegisterStoreType(name string, s Store) {
	//validate that this name hasn't already been registered
	if _, found := storeTypes[name]; found {
		panic("duplicate store type registered")
	}
	storeTypes[name] = s
}
