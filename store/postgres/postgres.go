package postgres

import "fmt"
import "github.com/cloudfoundry-community/portcullis/store"

// Postgres is an implementation of a Portcullis store that reads and writes
//   from a Postgres database
type Postgres struct {
}

func init() {
	store.RegisterStoreType("postgres", &Postgres{})
}

//Initialize checks the existing schema in the connected database and sets up
//the tables to store data in the up-to-date schema, updating and migrating as
//necessary
func (p *Postgres) Initialize(map[string]interface{}) error {
	//TODO
	return fmt.Errorf("Not yet implemented")
}

//ListMappings returns the list of all mappings stored in the Postgres database
func (p *Postgres) ListMappings() ([]store.Mapping, error) {
	//TODO
	return nil, fmt.Errorf("Not yet implemented")
}

//GetMapping returns a mapping corresponding to the name given. Errs if no
// mapping with that name exists in the Postgres database
func (p *Postgres) GetMapping(name string) (store.Mapping, error) {
	//TODO
	return store.Mapping{}, fmt.Errorf("Not yet implemented")
}

//AddMapping stores a new mapping in a row in the Postgres database
//Will return an error if a mapping with that name already exists in the db
func (p *Postgres) AddMapping(store.Mapping) error {
	//TODO
	return fmt.Errorf("Not yet implemented")
}

//EditMapping changes an existing entry for a Mapping with the same name as the
// in the Postgres database as the provided Mapping to have the same data as in
// the provided mapping. Errs if no mapping with that name exists in the database
func (p *Postgres) EditMapping(store.Mapping) error {
	//TODO
	return fmt.Errorf("Not yet implemented")
}

//DeleteMapping removes a mapping from the Postgres database, and errs if no
//such mapping exists
func (p *Postgres) DeleteMapping(name string) error {
	//TODO
	return fmt.Errorf("Not yet implemented")
}
