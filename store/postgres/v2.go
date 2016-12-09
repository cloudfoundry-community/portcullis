package postgres

import "github.com/starkandwayne/goutils/log"

type v2 struct {
}

func (v v2) migrate(p *Postgres) error {

	log.Debugf("Starting v2 Migration...")

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

	_, err = transaction.Exec(`CREATE TABLE mappings (
						 name      TEXT PRIMARY KEY,
						 location  TEXT NOT NULL,
						 config    TEXT NOT NULL DEFAULT '{}'
					 )`)
	if err != nil {
		log.Debugf("Failed perform command: %s", err.Error())
		return err
	}

	_, err = transaction.Exec(`UPDATE schema_info SET version = $1`, v.version())
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

func (v v2) version() int {
	return 2
}
