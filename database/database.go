package database

import (
	"database/sql"
	"time"
)

// DB contains the database connection pool pointer.
// It is used so that methods are declared on it, in order
// to mock the database easily, through the help of the Datastore interface
type DB struct {
	*sql.DB
}

// NewDB creates a database connection pool in DB and pings the database
func NewDB(host, user, password, database string) (*DB, error) {
	connStr := "postgres://" + user + ":" + password + "@" + host + "/" + database + "?sslmode=disable&connect_timeout=1"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	fails := 0
	for {
		err = db.Ping()
		if err == nil {
			break
		}
		fails++
		if fails == 3 {
			return nil, err
		}
		time.Sleep(200 * time.Millisecond)
	}
	return &DB{db}, nil
}

// CheckConnectivity pings the database and runs failCallback on failure
func (db *DB) CheckConnectivity(failCallback func(err error)) {
	if err := db.Ping(); err != nil {
		failCallback(err)
	}
}
