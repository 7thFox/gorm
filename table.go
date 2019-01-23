package gorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// CreateTable will create a new table based upon the DataObject given
func (db *DatabaseConnection) CreateTable(obj DataObject) error {
	var query strings.Builder
	t := reflect.TypeOf(obj)
	db.log("Creating table ", t.Name(), "...")

	query.WriteString("CREATE TABLE ")
	query.WriteString(db.Schema)
	query.WriteByte('.')
	query.WriteString(t.Name())
	query.WriteString(" (\n\t")

	hasPk := false
	everInserted := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		isPrimary := f.Tag.Get("primary") == "true"
		if isPrimary {
			if hasPk {
				return errors.New("Multiple PKs defined")
			}
			db.pkCache[t] = f.Name
			hasPk = true
		}

		everInserted = getColumnDef(&query, f, isPrimary, everInserted)

	}
	if !hasPk {
		return errors.New("No PK defined")
	}
	query.WriteString("\n);")

	db.log(query.String())

	_, err := db.connection.Exec(query.String())
	return err
}

// DropTable will drop the table related the given DataObject
func (db *DatabaseConnection) DropTable(obj DataObject) error {
	table := reflect.TypeOf(obj).Name()
	db.log("Dropping table ", table, "...")
	query := fmt.Sprintf("DROP TABLE %s.%s", db.Schema, table)
	_, err := db.connection.Exec(query)
	return err
}

var golangToPostgresTypes = map[reflect.Kind]string{
	reflect.String:  "VARCHAR",
	reflect.Bool:    "BOOLEAN",
	reflect.Int:     "BIGINT",
	reflect.Int8:    "SMALLINT",
	reflect.Int16:   "SMALLINT",
	reflect.Int32:   "INTEGER",
	reflect.Int64:   "BIGINT",
	reflect.Float32: "FLOAT(4)",
	reflect.Float64: "FLOAT(8)",
}

func getColumnDef(s *strings.Builder, f reflect.StructField, isPrimary, everInserted bool) bool {
	if f.Tag.Get("exclude") == "true" {
		return everInserted
	}
	if f.Type.Kind() == reflect.Struct {
		for i := 0; i < f.Type.NumField(); i++ {
			everInserted = getColumnDef(s, f.Type.Field(i), false, everInserted)
		}
		return everInserted
	}
	if everInserted {
		s.WriteString(",\n\t")
	}

	isUnique := f.Tag.Get("unique") == "true"
	size, sizeok := f.Tag.Lookup("size")

	s.WriteString(f.Name)
	s.WriteByte('\t')
	if isPrimary {

		s.WriteString("SERIAL PRIMARY KEY")
		return true
	}

	s.WriteString(golangToPostgresTypes[f.Type.Kind()])
	if sizeok {
		s.WriteByte('(')
		s.WriteString(size)
		s.WriteByte(')')
	}

	if isUnique {
		s.WriteString(" UNIQUE")
	}
	return true
}
