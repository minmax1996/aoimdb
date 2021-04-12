package main

import (
	"fmt"
	"reflect"

	"github.com/minmax1996/aoimdb/internal/aoimdb"
)

func main() {
	aoimdb.InitDatabaseController()
	aoimdb.SelectDatabase("default")
	aoimdb.CreateTable("default", "newTable",
		[]string{"id", "col2"},
		[]reflect.Kind{reflect.Int, reflect.String})
	table, _ := aoimdb.GetTable("default", "newTable")
	table.Insert([]string{"id", "col2"}, []interface{}{1, "Name1"})
	table.Insert([]string{"id", "col2"}, []interface{}{2, "Name2"})
	table.Insert([]string{"id", "col2"}, []interface{}{3, "Name3"})
	//returns copy of table
	fmt.Println(table.Select([]string{"col2", "id"}).Select([]string{"id"}))
}
