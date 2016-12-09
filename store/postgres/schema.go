package postgres

type schema interface {
	//Migrate is responsible for making the changes for this schema relative
	// to the last one. Migrate should not write the schema version to the
	// database - that is handled in the main postgres.go file to make sure it
	// remains consistent across schemas
	// Migration should be done as a single transaction if possible so that it
	// rolls back if there is a failure.
	//Schema version migrate() is responsible for writing it's version to the
	// table
	migrate(*Postgres) error
	version() int
}
