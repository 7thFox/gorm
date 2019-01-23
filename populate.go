package gorm

import (
	"fmt"
	"reflect"
)

// Populate fills the given DataObject with the info from the DB with PK of index
func (db *DatabaseConnection) Populate(obj DataObject, index int) error {
	e := reflect.ValueOf(obj).Elem()
	t := e.Type()
	pkname := db.getPrimaryKey(obj)

	db.log(fmt.Sprintf("Selecting %s with ID %d...", t.Name(), index))

	values := make([]interface{}, 0)
	for i := 0; i < t.NumField(); i++ {
		values = addFieldToValueList(values, e, i)
	}

	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE %s = %d;", db.Schema, t.Name(), pkname, index)
	err := db.connection.QueryRow(query).Scan(values...)

	return err
}

func addFieldToValueList(values []interface{}, e reflect.Value, i int) []interface{} {
	if e.Type().Field(i).Tag.Get("exclude") == "true" {
		return values
	}
	f := e.Field(i)
	ft := f.Type()
	if ft.Kind() == reflect.Struct {
		for j := 0; j < ft.NumField(); j++ {
			values = addFieldToValueList(values, f, j)
		}
		return values
	}

	return append(values, f.Addr().Interface())
}
