package gorm

import (
	"fmt"
	"reflect"
)

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
var supportedTypes = map[reflect.Kind]bool{
	reflect.String:  true,
	reflect.Bool:    true,
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Float32: true,
	reflect.Float64: true,
}

func getSQLValue(obj interface{}, f reflect.StructField) string {
	v := reflect.Indirect(reflect.ValueOf(obj)).FieldByName(f.Name)
	switch f.Type.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.String:
		return fmt.Sprintf("'%s'", v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if _, ok := f.Tag.Lookup("fk"); ok && v.Int() < 1 {
			return "NULL"
		}
		return fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	}
	panic("type not supported")
}
