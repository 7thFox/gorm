package main

import (
	"fmt"

	"github.com/7thFox/gorm"
)

type Bar string

type Foobar struct {
	FoobarID int    `primary:"true"`
	Name     string `unique:"true" size:"100"`
	Debug    bool   `exclude:"true"`

	NestedStruct struct {
		Bar Bar
		Baz float64
	}
}

func main() {
	// db and user should be set up beforehand
	db, _ := gorm.Open("myfoobarschema")
	defer db.Close()

	db.DropSchema()
	db.CreateSchema()
	db.CreateTable(Foobar{})

	foo := Foobar{
		Name:  "example1",
		Debug: true,
	}
	foo.NestedStruct.Baz = 5.6

	db.Save(&foo) // record inserted; struct's PK has been updated

	// myfoobarschema=# SELECT * FROM myfoobarschema.Foobar;
	// foobarid |   name   | bar | baz
	// ----------+----------+-----+-----
	//        1 | example1 |     | 5.6
	// (1 row)

	foo.Name = "example2"

	db.Save(&foo) // record updated

	// myfoobarschema=# SELECT * FROM myfoobarschema.Foobar;
	// foobarid |   name   | bar | baz
	// ----------+----------+-----+-----
	//        1 | example2 |     | 5.6
	// (1 row)

	foo2 := Foobar{}
	db.Populate(&foo2, 1) // selects with ID 1
	fmt.Printf("%#v\n", foo2)
	// $ main.Foobar{FoobarID:1, Name:"example2", Debug:false, NestedStruct:struct { Bar main.Bar; Baz float64 }{Bar:"", Baz:5.599999904632568}}
}
