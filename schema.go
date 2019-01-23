package gorm

import "fmt"

// CreateSchema will create the schema of the DatabaseConnection
func (db *DatabaseConnection) CreateSchema() error {
	db.log("Creating schema ", db.Schema, "...")
	query := fmt.Sprintf("CREATE SCHEMA %s", db.Schema)
	_, err := db.connection.Exec(query)
	return err
}

// DropSchema will drop the schema of the DatabaseConnection
func (db *DatabaseConnection) DropSchema() error {
	db.log("Dropping schema ", db.Schema, "...")
	query := fmt.Sprintf("DROP SCHEMA %s CASCADE", db.Schema)
	_, err := db.connection.Exec(query)
	return err
}
