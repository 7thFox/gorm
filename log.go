package gorm

import (
	"fmt"
	"strings"
)

func (db *DatabaseConnection) log(s ...string) {
	if db.debug {
		fmt.Println(strings.Join(s, ""))
	}
}
