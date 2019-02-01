package gorm

import (
	"fmt"
	"reflect"
	"strings"
)

// GetRelationshipIDs will return a list of FK values of a relationship, as well as fill the given list with
// empty structs of that object to the same size as the returned ID list
func (db *DatabaseConnection) GetRelationshipIDs(fkName string, fkValue int, relList interface{}) ([]int, error) {
	l := reflect.ValueOf(relList).Elem()
	t := l.Type().Elem()

	pkname := db.getPrimaryKey(reflect.New(t.Elem()).Interface())
	query := strings.Builder{}
	query.WriteString("SELECT ")
	query.WriteString(pkname)
	query.WriteString(" FROM ")
	query.WriteString(db.Schema)
	query.WriteRune('.')
	query.WriteString(t.Elem().Name())
	query.WriteString("\nWHERE ")
	query.WriteString(fkName)
	query.WriteString(" = ")
	query.WriteString(fmt.Sprintf("%d", fkValue))
	query.WriteRune(';')

	// return errors.New(query.String())

	rows, err := db.connection.Query(query.String())
	if err != nil {
		return nil, err
	}

	ids := []int{}
	for rows.Next() {
		id := 0
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
		v := reflect.New(t.Elem())
		// if err := db.Populate(v.Interface().(DataObject), id); err != nil {
		// 	return err
		// }

		l.Set(reflect.Append(l, v))
	}

	return ids, nil
}
