package gorm

import (
	"fmt"
	"reflect"

	"database/sql"
	// postgresql provider
	_ "github.com/lib/pq"
)

// DataObject is defines the requirement of structs to be used with this package.
// Currently this interface is empty, but may be used in the future.
type DataObject interface {
}

// DatabaseConnection represents an open DB connection
type DatabaseConnection struct {
	Schema string

	connection *sql.DB
	pkCache    map[reflect.Type]string
	debug      bool
}

// Open opens a DB connection using the same user/db/schema
func Open(uds string) (*DatabaseConnection, error) {
	return OpenVerbose(uds, uds, uds, false)
}

// OpenVerbose opens a DB allowing to specify different user/db/schema and enable debug log output
func OpenVerbose(user, database, schema string, debug bool) (*DatabaseConnection, error) {
	db := DatabaseConnection{Schema: schema, pkCache: map[reflect.Type]string{}, debug: debug}
	db.log("Opening DB connection...")
	connSettings := fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, database)
	c, err := sql.Open("postgres", connSettings)
	if err != nil {
		return nil, err
	}
	db.connection = c
	return &db, nil
}

// Close closes the db connection
func (db *DatabaseConnection) Close() error {
	db.log("Closing DB connnection...")
	err := db.connection.Close()
	return err
}
