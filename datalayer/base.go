package datalayer

import (
	"database/sql"
	"errors"

	// import the MySQL Driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/chytilp/links/config"
)

var (
	db *sql.DB

	// ErrNotFound is returned when the no records where matched by the query
	ErrNotFound = errors.New("not found")
)

// custom type so we can convert sql results to easily
type scanner func(dest ...interface{}) error

func getDb() *sql.DB {
	var err error
	db, err = sql.Open("mysql", config.App.Database.GetConnectionString())
	if err != nil {
		// if the DB cannot be accessed -> panic
		//logging.L.Error("Error from sql.open. err: %s", err)
		panic(err.Error())
	}
	return db
}

// NewRecords creates instance of records object.
func newRecords(db *sql.DB) *records {
	if db == nil {
		db = getDb()
	}
	records := &records{
		db: db,
	}
	return records
}

// records is generic object for common methods above one db table.
type records struct {
	db *sql.DB
}

// insert is generic method for insert record to db.
func (r *records) insert(values []interface{}, expression string) (int, error) {
	stmt, err := r.db.Prepare(expression)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastID), nil
}

// update is generic method for update record to db.
func (r *records) update(values []interface{}, expression string) error {
	stmt, err := r.db.Prepare(expression)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	return nil
}

// close db connection.
func (r *records) close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}
