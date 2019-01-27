package gorm

import (
	"reflect"
	"strings"
)

func (db *DatabaseConnection) insert(obj DataObject) error {
	var query strings.Builder
	t := reflect.ValueOf(obj).Elem().Type()
	pkname := db.getPrimaryKey(obj)

	db.log("Inserting record into ", t.Name(), "...")

	query.WriteString("INSERT INTO ")
	query.WriteString(db.Schema)
	query.WriteRune('.')
	query.WriteString(t.Name())
	query.WriteString(" (\n\t")

	var values strings.Builder
	everInserted := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Name == pkname {
			continue
		}

		everInserted = getInsertCol(obj, &query, &values, f, everInserted)
	}

	query.WriteString("\n)\nVALUES (")
	query.WriteString(values.String())
	query.WriteString(")\nRETURNING ")
	query.WriteString(db.getPrimaryKey(obj))
	query.WriteRune(';')

	var id int
	if err := db.connection.QueryRow(query.String()).Scan(&id); err != nil {
		return errorWithQuery(err, query.String())
	}

	db.setPrimaryKeyValue(obj, id)
	return nil
}

func getInsertCol(obj interface{}, query, values *strings.Builder, f reflect.StructField, everInserted bool) bool {
	if f.Tag.Get("exclude") == "true" {
		return everInserted
	}
	if isFlattenableStruct(f) {
		innerObj := reflect.Indirect(reflect.ValueOf(obj)).FieldByName(f.Name).Interface()
		for i := 0; i < f.Type.NumField(); i++ {
			everInserted = getInsertCol(innerObj, query, values, f.Type.Field(i), everInserted)
		}
		return everInserted
	}
	if !isSupported(f) {
		return everInserted
	}
	if everInserted {
		values.WriteString(", ")
		query.WriteString(",\n\t")
	}
	query.WriteString(f.Name)
	values.WriteString(getSQLValue(obj, f))
	return true
}
