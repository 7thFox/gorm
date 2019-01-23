# gorm
GORM is a simple golang ORM for PostgreSQL built for very simple SQL interactions.
If the usage section below isn't enough for you, I would recommend reading the source. 
It's <500 lines and is basically just a exercise in golang's reflection.

There are a few quirks, such as the PK not being able to reside within a nested Struct.

Also I don't have any feature to delete a record, but I'm sure I'll do that one day...

# usage

```golang

type Foobar struct {
	FoobarID   int    `primary:"true"`
	Name       string `unique:"true" size:"100"`
	Debug      string `exclude:"true"`

	NestedStruct struct {
		Bar       Bar // type Bar string
		Baz       float64
	}
}

func main() {
	db, _ = gorm.Open("myfoobarschema")
	defer db.Close()

	db.DropSchema()
	db.CreateSchema()
	db.CreateTable(Foobar{})

	foo := Foobar{}
  // set values...

	db.Save(&foo)// record inserted; PK has been updated
  // change values...

	db.Save(&foo)// record updated

	foo2 := Foobar{}
	db.Populate(&foo2, 1)// selects with ID 1
	fmt.Printf("%#v\n", foo2)
}

```
