package postgres

import "github.com/starkandwayne/goutils/log"

type v1 struct {
}

func (v v1) migrate(p *Postgres) error {

	log.Debugf("Starting v1 Migration...")

	transaction, err := p.connection.Begin()

	defer func() {
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				log.Infof("Failed to roll back transaction: %s", err.Error())
			} else {
				log.Infof("Rolled back transaction for v1")
			}
		}
	}()

  // Creates the schema_info table to hold the version number of the schemas in
	// this repo.  Schema versions should be sequential step 1.  Done in a
	// transaction.
	_, err = transaction.Exec(`CREATE TABLE schema_info (
               version INTEGER
             )`)
	if err != nil {
		log.Debugf("Failed perform command: %s", err.Error())
		return err
	}
	//defer p.connection.Close() //needed for Exec?

  // Inserts the very first version (1) which is simply that the table was created
	// successfully as well as this row was inserted transactionally.
	_, err = transaction.Exec(`INSERT INTO schema_info VALUES ($1)`, v.version())
	if err != nil {
		log.Debugf("Failed perform command: %s", err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	return nil
}

func (v v1) version() int {
	return 1
}
