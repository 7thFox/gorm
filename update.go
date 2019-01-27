package gorm

import (
	"fmt"
	"reflect"
	"strings"
)

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

	return errorWithQuery(err, query.String())
}

func getUpdateCol(obj interface{}, query *strings.Builder, f reflect.StructField, everInserted bool) bool {
	if f.Tag.Get("exclude") == "true" {
		return everInserted
	}
	if isFlattenableStruct(f) {
		innerObj := reflect.Indirect(reflect.ValueOf(obj)).FieldByName(f.Name).Interface()
		for i := 0; i < f.Type.NumField(); i++ {
			everInserted = getUpdateCol(innerObj, query, f.Type.Field(i), everInserted)
		}
		return everInserted
	}
	if !isSupported(f) {
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
