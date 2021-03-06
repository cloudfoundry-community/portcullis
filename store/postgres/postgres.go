package postgres

//This package will register this store type with the store package when it is
// imported (use an underscore import)
//CONFIG:
//  host:     (string) The hostname/ip address of the Postgres server
//  port:			(int)    The port on which the postgres server is listening
//  dbname:   (string) The name of the database to connect to
//  username: (string) The name of the user to connect with
//  password: (string) The password of the user specified with `username`

import (
	"database/sql"
	"fmt"

	"encoding/json"

	"github.com/cloudfoundry-community/portcullis/broker/bindparser"

	"github.com/cloudfoundry-community/portcullis/config"
	"github.com/cloudfoundry-community/portcullis/store"
	"github.com/lib/pq"
	"github.com/starkandwayne/goutils/log"
)

// Postgres is an implementation of a Portcullis store that reads and writes
//   from a Postgres database
type Postgres struct {
	connection *sql.DB
}

type postgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

const (
	schemaTable   = "schema_info"
	mappingsTable = "mappings"
)

//If you're making a new schema, it needs to be added to the end of this array
var schemas = map[int]schema{
	1: v1{},
	2: v2{},
}

func init() {
	store.RegisterStoreType("postgres", &Postgres{})
}

func (p *Postgres) getSchemaVersion() (int, error) {

	r, err := p.connection.Query(`SELECT version FROM schema_info LIMIT 1`)
	if err != nil {
		if err.(*pq.Error).Code == "42P01" { //undefined table...
			return 0, nil
		}
		return -1, err
	}
	defer r.Close()

	// no records = no schema
	if !r.Next() {
		return 0, nil
	}

	var v int
	err = r.Scan(&v)
	// failed unmarshall is an actual error
	if err != nil {
		return 0, err
	}

	// invalid (negative) schema version is an actual error
	if v < 0 || v > len(schemas) {
		return 0, fmt.Errorf("Invalid schema version %d found", v)
	}

	return int(v), nil

}

//Initialize checks the existing schema in the connected database and sets up
//the tables to store data in the up-to-date schema, updating and migrating as
//necessary
func (p *Postgres) Initialize(conf map[string]interface{}) error {
	var err error
	if conf == nil {
		return fmt.Errorf("Postgres store config is nil")
	}
	if err = config.ValidateConfigKeys(config.StoreKey, conf, "host", "port", "dbname", "username", "password"); err != nil {
		return err
	}

	pgConf := postgresConfig{}
	config.ParseMapConfig(config.StoreKey, conf, &pgConf)

	p.connection, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=5", pgConf.Username, pgConf.Password, pgConf.Host, pgConf.Port, pgConf.DBName))
	if err != nil {
		log.Infof(err.Error())
		return err
	}

	if err = p.connection.Ping(); err != nil {
		log.Infof(err.Error())
		return err
	}

	//Entry point into initializing schema...
	var currentVersion int
	currentVersion, err = p.getSchemaVersion()
	if err != nil {
		log.Debugf("schema error")
		return err
	}

	log.Debugf("schema found: %d", currentVersion)

	for currentVersion < len(schemas) {
		if err = schemas[currentVersion+1].migrate(p); err != nil {
			log.Debugf(err.Error())
			return err
		}

		currentVersion, err = p.getSchemaVersion()
		if err != nil {
			log.Debugf(err.Error())
			return err
		}
	}

	return nil
}

//ListMappings returns the list of all mappings stored in the Postgres database
func (p *Postgres) ListMappings() ([]store.Mapping, error) {
	log.Debugf("Attempting to retrieve all rows from mappings table...")
	rows, err := p.connection.Query("SELECT name, location, config FROM mappings")
	if err != nil {
		log.Infof("Scan error attempting to retrieve all rows from mapping")
		return []store.Mapping{}, err
	}
	ret := []store.Mapping{}
	for rows.Next() {
		var name, location, mappingConfig string
		err := rows.Scan(&name, &location, &mappingConfig)
		if err != nil {
			log.Infof("Scan error attempting to retrieve all rows from mapping")
			return []store.Mapping{}, err
		}

		var bc bindparser.Config

		if err := json.Unmarshal([]byte(mappingConfig), &bc); err != nil {
			log.Infof("Scan error attempting to retrieve row with name: %s, cloud not unmarshal config %s", name, err.Error())
		}

		log.Debugf("Found row: %s, %s", name, location)
		ret = append(ret, store.Mapping{Name: name, Location: location, BindConfig: bc})
	}

	return ret, err
}

//GetMapping returns a mapping corresponding to the name given. Errs if no
// mapping with that name exists in the Postgres database
func (p *Postgres) GetMapping(name string) (store.Mapping, error) {
	log.Debugf("Attempting to get a row from the mappings table...")

	ret := store.Mapping{}
	var mappingName, location, mappingConfig string

	err := p.connection.QueryRow("SELECT name, location, config FROM mappings WHERE name = $1", name).Scan(&mappingName, &location, &mappingConfig)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("No rows found while attempting to retrieve row with name: %s", name)
			return ret, store.ErrNotFound

		}
		log.Infof("Scan error attempting to retrieve row with name: %s", name)
		return ret, err
	}

	var bc bindparser.Config

	if err := json.Unmarshal([]byte(mappingConfig), &bc); err != nil {
		log.Infof("Scan error attempting to retrieve row with name: %s, cloud not unmarshal config %s", name, err.Error())
	}

	log.Debugf("Found row with name:%s and location: %s", mappingName, location)
	retMapping := store.Mapping{
		Name:       mappingName,
		Location:   location,
		BindConfig: bc,
	}

	return retMapping, err
}

//AddMapping stores a new mapping in a row in the Postgres database
//Will return an error if a mapping with that name already exists in the db
func (p *Postgres) AddMapping(m store.Mapping) error {

	log.Debugf("Attempting to add a row into mappings table...")

	bc, _ := json.Marshal(m.BindConfig)

	_, err := p.connection.Exec(`INSERT INTO mappings (name, location, config) VALUES ($1, $2, $3)`, m.Name, m.Location, string(bc))
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			log.Infof("Could not insert into %s table, duplicate row: %s", mappingsTable, err.Error())
			return store.ErrDuplicate
		}
		log.Infof("Could not insert into %s table: %s", mappingsTable, err.Error())
	}
	return err
}

//EditMapping changes an existing entry for a Mapping with the same name as the
// in the Postgres database as the provided Mapping to have the same data as in
// the provided mapping. Errs if no mapping with that name exists in the database
func (p *Postgres) EditMapping(name string, m store.Mapping) error {
	log.Debugf("Attempting to update a row in mappings table...")

	var numRows int
	err := p.connection.QueryRow(`SELECT COUNT(name) FROM mappings WHERE name = $1`, name).Scan(&numRows)
	if err != nil {
		log.Infof("Scan error attempting to retrieve row with name: %s", name)
		return err
	}

	if numRows < 1 {
		log.Infof("No mappings found for key value: %s", name)
		return store.ErrNotFound
	}

	bc, _ := json.Marshal(m.BindConfig)

	_, err = p.connection.Exec(`UPDATE mappings SET name = $1, location = $2, config = $3 WHERE name = $4`, m.Name, m.Location, string(bc), name)

	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			log.Infof("Could not insert into %s table, duplicate row: %s", mappingsTable, err.Error())
			return store.ErrDuplicate
		}
		log.Infof("Could not update mappings entry %s to become (%s, %s) : %s", name, m.Name, m.Location, err.Error())
	}
	return err
}

//DeleteMapping removes a mapping from the Postgres database, and errs if no
//such mapping exists
func (p *Postgres) DeleteMapping(name string) error {
	log.Debugf("Attempting to delete a row from mappings table...")

	var numRows int
	err := p.connection.QueryRow(`SELECT COUNT(name) FROM mappings WHERE name = $1`, name).Scan(&numRows)
	if err != nil {
		log.Infof("Scan error attempting to retrieve row with name: %s", name)
		return err
	}

	if numRows < 1 {
		log.Infof("No mappings found for key value: %s", name)
		return store.ErrNotFound
	}

	_, err = p.connection.Exec(`DELETE FROM mappings WHERE name = $1`, name)
	if err != nil {
		log.Infof("Could not delete mappings entry %s: %s", name, err.Error())
	}
	return err
}

//Size returns the number of mapping rows in the Postgres database
func (p *Postgres) Size() (int, error) {
	log.Debugf("Getting the row count in the mappings table...")

	var numRows int
	err := p.connection.QueryRow(`SELECT COUNT(name) FROM mappings`).Scan(&numRows)
	if err != nil {
		log.Infof("Scan error attempting to retrieve row count")
		return 0, err
	}

	return numRows, nil

}

//ClearMappings removes all mappings from the Postgres database by truncating
//the mapping table
func (p *Postgres) ClearMappings() error {
	log.Debugf("Truncating table mappings...")

	_, err := p.connection.Exec(`TRUNCATE TABLE mappings`)

	if err != nil {
		log.Infof("Could not TRUNCATE TABLE mappings: %s", err.Error())
	}
	return err
}

//TODO: Comment
func (p *Postgres) GetSecGroupInfoByName(name string) (result store.SecGroupInfo, err error) {
	return result, fmt.Errorf("Not yet implemented")
}

//TODO: Comment
func (p *Postgres) GetSecGroupInfoByInstance(GUID string) (result store.SecGroupInfo, err error) {
	return result, fmt.Errorf("Not yet implemented")
}

//TODO: Comment
func (p *Postgres) AddSecGroupInfo(toAdd store.SecGroupInfo) error {
	return fmt.Errorf("Not yet implemented")
}

//TODO: Comment
func (p *Postgres) DeleteSecGroupInfoByInstance(GUID string) error {
	return fmt.Errorf("Not yet implemented")
}

//TODO: Comment
func (p *Postgres) DeleteSecGroupInfoByName(name string) error {
	return fmt.Errorf("Not yet implemented")
}

//TODO: Comment
func (p *Postgres) NumSecGroupInfo() (int, error) {
	return -1, fmt.Errorf("Not yet implemented")
}

//TODO: Comment
func (p *Postgres) ClearSecGroupInfo() error {
	return fmt.Errorf("Not yet implemented")
}
