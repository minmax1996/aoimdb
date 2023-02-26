package main

import (
	"fmt"
	"strconv"

	aoimdb "github.com/minmax1996/aoimdb/cmd/aoimd/database"
	"github.com/minmax1996/aoimdb/cmd/aoimd/datatypes"
	"github.com/minmax1996/aoimdb/cmd/aoimd/table"
)

func main() {
	aoimdb.InitDatabaseController()
	aoimdb.SelectDatabase("default")
	aoimdb.CreateTable("default", "newTable",
		[]string{"id", "name"},
		[]datatypes.Datatype{datatypes.Int, datatypes.String})

	ttable, _ := aoimdb.GetTable("default", "newTable")

	for i := 0; i < 1000; i++ {
		ttable.Insert([]string{"id", "name"}, []interface{}{i, "Name" + strconv.Itoa(i)})
	}

	fmt.Println("Inserted")
	filtered1 := ttable.Filter(func(m table.Row) bool {
		return m[0].(int)%21 == 0
	})
	fmt.Println("filtered1------------------------------")
	fmt.Println(filtered1)

	filtered2 := filtered1.Filter(func(m table.Row) bool {
		return m[0].(int)%5 == 0
	})
	fmt.Println("filtered2------------------------------")
	fmt.Println(filtered2)
	filtered2.Delete(func(m table.Row) bool {
		return m[0].(int)%2 == 0
	})
	fmt.Println("afterDelete------------------------------")
	fmt.Println(ttable)
}
