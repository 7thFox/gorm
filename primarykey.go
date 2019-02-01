package gorm

import (
	"reflect"
)

func (db *DatabaseConnection) getPrimaryKeyValue(obj interface{}) int {
	return int(reflect.ValueOf(obj).Elem().FieldByName(db.getPrimaryKey(obj)).Int())
}

func (db *DatabaseConnection) setPrimaryKeyValue(obj DataObject, id int) {
	pk := reflect.ValueOf(obj).Elem().FieldByName(db.getPrimaryKey(obj))
	asSettableUnexported(pk).SetInt(int64(id))
}

func (db *DatabaseConnection) getPrimaryKey(obj interface{}) string {
	t := reflect.ValueOf(obj).Elem().Type()
	if val, ok := db.pkCache[t]; ok {
		return val
	}

	hasPk := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("primary") == "true" {
			if hasPk {
				panic("Multiple PKs defined")
			}
			db.pkCache[t] = f.Name
			hasPk = true
		}
	}
	if !hasPk {
		panic("No PK defined")
	}

	return db.pkCache[t]
}
