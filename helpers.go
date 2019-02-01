package gorm

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func (db *DatabaseConnection) log(s ...string) {
	if db.debug {
		fmt.Println(strings.Join(s, ""))
	}
}

func errorWithQuery(err error, query string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s\n\nQuery:\n%s", err.Error(), query)
}

func isFlattenableStruct(f reflect.StructField) bool {
	return f.Type.Kind() == reflect.Struct && f.Tag.Get("flatten") == "true"
}

func isSupported(f reflect.StructField) bool {
	return supportedTypes[f.Type.Kind()]
}

func asSettableUnexported(val reflect.Value) reflect.Value {
	return reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()
}
