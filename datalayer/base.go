package datalayer

import (
	"database/sql"
	"errors"

	// import the MySQL Driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/chytilp/links/config"
	"github.com/chytilp/links/logging"
)

var (
	db *sql.DB

	// ErrNotFound is returned when the no records where matched by the query
	ErrNotFound = errors.New("not found")
)

// custom type so we can convert sql results to easily
type scanner func(dest ...interface{}) error

func getDb() (*sql.DB, error) {
	if db == nil {
		if config.App == nil {
			err := errors.New("config is not initialized")
			logging.L.Error("Error from getDb. err: %s", err)
			return nil, err
		}

		var err error
		db, err = sql.Open("mysql", config.App.Database.GetConnectionString())
		if err != nil {
			// if the DB cannot be accessed -> panic
			logging.L.Error("Error from sql.open. err: %s", err)
			panic(err.Error())
		}
	}

	return db, nil
}

// only for test purposes
func setDb(database *sql.DB) {
	db = database
}

func init() {
	// Ensure the config is loaded and the db initialized.
	_, _ = getDb()
}
