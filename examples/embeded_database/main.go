package main

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/minmax1996/aoimdb/internal/aoimdb"
	"github.com/minmax1996/aoimdb/internal/aoimdb/datatypes"
)

func main() {
	aoimdb.InitDatabaseController()
	aoimdb.SelectDatabase("default")
	aoimdb.CreateTable("default", "newTable",
		[]string{"id", "name"},
		[]reflect.Type{datatypes.Int, datatypes.String})

	table, _ := aoimdb.GetTable("default", "newTable")

	for i := 0; i < 1000; i++ {
		table.Insert([]string{"id", "name"}, []interface{}{i, "Name" + strconv.Itoa(i)})
	}

	fmt.Println("Inserted")
	filtered1 := table.Filter(func(m datatypes.Row) bool {
		return m[0].(int)%21 == 0
	})
	fmt.Println("filtered1------------------------------")
	fmt.Println(filtered1)

	filtered2 := filtered1.Filter(func(m datatypes.Row) bool {
		return m[0].(int)%5 == 0
	})
	fmt.Println("filtered2------------------------------")
	fmt.Println(filtered2)
	filtered2.Delete(func(m datatypes.Row) bool {
		return m[0].(int)%2 == 0
	})
	fmt.Println("afterDelete------------------------------")
	fmt.Println(table)
}
