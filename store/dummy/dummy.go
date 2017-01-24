package dummy

//This package will register this store type with the store package when it is
// imported (use an underscore import)
//CONFIG:
//  confirm: (boolean) true if this dummy should actually be initialized.

import (
	"fmt"

	"github.com/cloudfoundry-community/portcullis/config"
	"github.com/cloudfoundry-community/portcullis/store"
)

//Dummy is an in-memory store that is used for testing and probably shouldn't
// be used in any real circumstances, unless you really don't care about your
// data or its persistence.
//TODO: Might want to mutex these calls at some point. Stretchiest of goals.
type Dummy struct {
	storage     map[string]store.Mapping
	secgroups   map[string]store.SecGroupInfo
	initialized bool
}

type dummyConfig struct {
	Confirm bool `yaml:"confirm"`
}

func init() {
	store.RegisterStoreType("dummy", &Dummy{})
}

//Initialize sets up the Dummy store to be useable
func (d *Dummy) Initialize(conf map[string]interface{}) error {
	if conf == nil {
		return fmt.Errorf("Dummy store config is nil")
	}
	if err := config.ValidateConfigKeys(config.StoreKey, conf, "confirm"); err != nil {
		return err
	}

	dummyConf := dummyConfig{}
	config.ParseMapConfig(config.StoreKey, conf, &dummyConf)

	if !dummyConf.Confirm {
		return fmt.Errorf("Dummy store config key `confirm` not set to true")
	}

	d.storage = map[string]store.Mapping{}
	d.initialized = true
	return nil
}

//ListMappings returns all of the Mappings in the Dummy store
func (d *Dummy) ListMappings() ([]store.Mapping, error) {
	if !d.initialized {
		return nil, fmt.Errorf("Dummy not initialized")
	}

	ret := []store.Mapping{}
	for _, m := range d.storage {
		ret = append(ret, m)
	}
	return ret, nil
}

//GetMapping returns the mapping with the given name in the Dummy store
func (d *Dummy) GetMapping(name string) (ret store.Mapping, err error) {
	if !d.initialized {
		return ret, fmt.Errorf("Dummy not initialized")
	}

	var found bool
	ret, found = d.storage[name]
	if !found {
		err = store.ErrNotFound
	}
	return
}

//AddMapping adds a new mapping with a unique name to the Dummy store.
// Returns an error if a mapping with that name already exists
func (d *Dummy) AddMapping(m store.Mapping) error {
	if !d.initialized {
		return fmt.Errorf("Dummy not initialized")
	}

	if _, err := d.GetMapping(m.Name); err == nil {
		return store.ErrDuplicate
	}

	d.storage[m.Name] = m
	return nil
}

//EditMapping edits an existing mapping with the same name as the one in the
// provided store.Mapping object. The resulting mapping in the store will have
// all the same values as the one in the provided store.Mapping
func (d *Dummy) EditMapping(name string, m store.Mapping) error {
	if !d.initialized {
		return fmt.Errorf("Dummy not initialized")
	}

	if _, err := d.GetMapping(name); err != nil {
		return err
	}

	//Check if the name to edit to already exists in the store
	if name != m.Name {
		//dummy won't return error other than NotFound for GetMapping (unless not
		// initialized, but we would've already caught that). Other stores need to
		// handle connection errors and the like
		if _, err := d.GetMapping(m.Name); err == nil {
			//if there's no error, that means it was found and this is a duplicate
			return store.ErrDuplicate
		}
	}

	delete(d.storage, name)
	d.storage[m.Name] = m
	return nil
}

//DeleteMapping removes a mapping from the Dummy store if it exists, and
// returns an error otherwise
func (d *Dummy) DeleteMapping(name string) error {
	if !d.initialized {
		return fmt.Errorf("Dummy not initialized")
	}

	if _, err := d.GetMapping(name); err != nil {
		return store.ErrNotFound
	}

	delete(d.storage, name)
	return nil
}

//Size returns the length of the storage map
func (d *Dummy) Size() (int, error) {
	if !d.initialized {
		return -1, fmt.Errorf("Dummy not initialized")
	}

	return len(d.storage), nil
}

//ClearMappings makes the internal storage a new memory map, wiping out all
//preexisting data
func (d *Dummy) ClearMappings() error {
	if !d.initialized {
		return fmt.Errorf("Dummy not initialized")
	}
	d.storage = map[string]store.Mapping{}
	return nil
}

//GetSecGroupInfoByName returns the SecGroupInfo in the map with that name if it
// exists and returns ErrNotFound otherwise
// The map is indexed by ServiceInstanceGUID, not this, so this access takes O(n)
func (d *Dummy) GetSecGroupInfoByName(name string) (result store.SecGroupInfo, err error) {
	if !d.initialized {
		return result, fmt.Errorf("Dummy not initialized")
	}

	for _, secgroup := range d.secgroups {
		if secgroup.SecGroupName == name {
			return secgroup, nil
		}
	}

	return result, store.ErrNotFound
}

//GetSecGroupInfoByInstance returns the SecGroupInfo in the map with the given
// ServiceInstanceGUID value if it exists, and returns ErrNotFound otherwise.
func (d *Dummy) GetSecGroupInfoByInstance(GUID string) (result store.SecGroupInfo, err error) {
	if !d.initialized {
		return result, fmt.Errorf("Dummy not initialized")
	}

	if secgroup, ok := d.secgroups[GUID]; ok {
		return secgroup, nil
	}

	return result, store.ErrNotFound
}

//AddSecGroupInfo puts a copy of the given SecGroupInfo object into the map.
// ErrDuplicate is thrown if a SecGroupInfo with that ServiceInstanceGUID already
// exists.
func (d *Dummy) AddSecGroupInfo(toAdd store.SecGroupInfo) error {
	if !d.initialized {
		return fmt.Errorf("Dummy not initialized")
	}

	if _, exists := d.secgroups[toAdd.ServiceInstanceGUID]; exists {
		return store.ErrDuplicate
	}

	d.secgroups[toAdd.ServiceInstanceGUID] = toAdd

	return nil
}

//DeleteSecGroupInfoByInstance finds the SecGroupInfo object in the map with the
// given ServiceInstanceGUID and then, if it exists, it removes it from the map.
// Otherwise, it returns ErrNotFound
func (d *Dummy) DeleteSecGroupInfoByInstance(GUID string) error {
	if !d.initialized {
		return fmt.Errorf("Dummy not initialized")
	}

	if _, found := d.secgroups[GUID]; !found {
		return store.ErrNotFound
	}

	delete(d.secgroups, GUID)
	return nil
}

//DeleteSecGroupInfoByName finds the SecGroupInfo object in the map with the
// given SecGroupName and then, if it exists, it removes it from the map.
// Otherwise, it returns err not found.
// The search for the item to delete takes O(n)
func (d *Dummy) DeleteSecGroupInfoByName(name string) error {
	if !d.initialized {
		return fmt.Errorf("Dummy not initialized")
	}

	secgroup, err := d.GetSecGroupInfoByName(name)
	if err != nil {
		return err
	}

	delete(d.secgroups, secgroup.ServiceInstanceGUID)
	return nil
}
