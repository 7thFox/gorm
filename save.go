package gorm

import (
	"fmt"
	"reflect"
	"strings"
)

// Save creates or updates a DB entry for the given DataObject.
// Create/Update logic is based on the PK being > 0 for updates.
func (db *DatabaseConnection) Save(obj DataObject) error {
	pk := db.getPrimaryKeyValue(obj)
	if pk > 0 {
		return db.update(obj)
	}
	return db.insert(obj)
}

func (db *DatabaseConnection) update(obj DataObject) error {
	var query strings.Builder
	t := reflect.ValueOf(obj).Elem().Type()
	pkname := db.getPrimaryKey(obj)
	pkvalue := fmt.Sprintf("%d", db.getPrimaryKeyValue(obj))

	db.log("Updating ", t.Name(), " with ID ", pkvalue, "...")

	query.WriteString("UPDATE ")
	query.WriteString(db.Schema)
	query.WriteRune('.')
	query.WriteString(t.Name())
	query.WriteString("\nSET \n\t")

	everInserted := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Name == pkname {
			continue
		}

		everInserted = getUpdateCol(obj, &query, f, everInserted)
	}

	query.WriteString("\nWHERE ")
	query.WriteString(pkname)
	query.WriteString(" = ")
	query.WriteString(pkvalue)
	query.WriteRune(';')

	_, err := db.connection.Exec(query.String())

	return err
}

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
		return err
	}

	db.setPrimaryKeyValue(obj, id)
	return nil
}

func getSQLValue(obj interface{}, f reflect.StructField) string {
	v := reflect.Indirect(reflect.ValueOf(obj)).FieldByName(f.Name)

	switch f.Type.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.String:
		return fmt.Sprintf("'%s'", v.String())
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Float32:
	case reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	}
	return f.Name
}

func getInsertCol(obj interface{}, query, values *strings.Builder, f reflect.StructField, everInserted bool) bool {
	if f.Tag.Get("exclude") == "true" {
		return everInserted
	}
	if f.Type.Kind() == reflect.Struct {
		innerObj := reflect.Indirect(reflect.ValueOf(obj)).FieldByName(f.Name).Interface()
		for i := 0; i < f.Type.NumField(); i++ {
			everInserted = getInsertCol(innerObj, query, values, f.Type.Field(i), everInserted)
		}
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

func getUpdateCol(obj interface{}, query *strings.Builder, f reflect.StructField, everInserted bool) bool {
	if f.Tag.Get("exclude") == "true" {
		return everInserted
	}
	if f.Type.Kind() == reflect.Struct {
		innerObj := reflect.Indirect(reflect.ValueOf(obj)).FieldByName(f.Name).Interface()
		for i := 0; i < f.Type.NumField(); i++ {
			everInserted = getUpdateCol(innerObj, query, f.Type.Field(i), everInserted)
		}
		return everInserted
	}
	if everInserted {
		query.WriteString(",\n\t")
	}
	query.WriteString(f.Name)
	query.WriteString(" = ")
	query.WriteString(getSQLValue(obj, f))
	return true
}
