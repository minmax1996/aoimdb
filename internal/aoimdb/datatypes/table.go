package datatypes

import (
	"errors"
	"time"
)

type ColumnType int

const (
	IntType ColumnType = iota
	StringType
	DateType
)

//
func checkType(val interface{}, res ColumnType) bool {
	switch val.(type) {
	case int:
		return res == IntType
	case string:
		return res == StringType
	case time.Time:
		return res == DateType
	default:
		return false
	}
}

type Row []interface{}

type Table struct {
	Name        string
	ColumnNames []string
	ColumnTypes []ColumnType
	DataRows    []Row
}

//NewTableSchema initialize new table scheam with no data
func NewTableSchema(tableName string, names []string, types []ColumnType) *Table {
	if len(names) != len(types) {
		return nil
	}
	return &Table{
		Name:        tableName,
		ColumnNames: names,
		ColumnTypes: types,
		DataRows:    make([]Row, 0),
	}
}

//NewTableWithRows initailexe new table with rows
func NewTableWithRows(tableName string, names []string, types []ColumnType, rows []Row) *Table {
	if len(names) != len(types) {
		return nil
	}
	return &Table{
		Name:        tableName,
		ColumnNames: names,
		ColumnTypes: types,
		DataRows:    rows,
	}
}

//Insert inserts into rows new row with values if values pass typecheck
func (t *Table) Insert(names []string, values []interface{}) error {
	if len(names) != len(values) {
		return errors.New("lens not equal")
	}

	row := make(Row, len(t.ColumnNames))
	for i, val := range values {
		ind := findIndex(t.ColumnNames, names[i])
		if ind == -1 || !checkType(val, t.ColumnTypes[ind]) {
			return errors.New("cant append value")
		}
		row[ind] = val
	}
	t.DataRows = append(t.DataRows, row)
	//TODO insert indexes here
	return nil
}

//Select returns filtered table with only given names in given order
func (t *Table) Select(names []string) *Table {
	indexes := findIndexes(t.ColumnNames, names)
	if len(indexes) == 0 {
		return nil
	}

	//make empty table with fixed lens
	resultTable := &Table{
		ColumnNames: make([]string, len(indexes)),
		ColumnTypes: make([]ColumnType, len(indexes)),
		DataRows:    make([]Row, len(t.DataRows)),
	}

	//initiate all datarows with zero elements and capacity eq to len of all fields
	for j := 0; j < len(t.DataRows); j++ {
		resultTable.DataRows[j] = make(Row, 0, len(indexes))
	}

	for i, index := range indexes {
		resultTable.ColumnNames[i] = t.ColumnNames[index]
		resultTable.ColumnTypes[i] = t.ColumnTypes[index]
		for j := 0; j < len(t.DataRows); j++ { //for each datarow add to result one
			resultTable.DataRows[j] = append(resultTable.DataRows[j], t.DataRows[j][index])
		}
	}

	return resultTable
}

//findIndexes finds all indexes that correspond all columns names in collection
func findIndexes(colletion []string, values []string) []int {
	result := make([]int, 0)
	for _, vv := range values {
		for i, v := range colletion {
			if v == vv {
				result = append(result, i)
			}
		}
	}
	return result
}

//findIndexes finds all indexes that correspond all columns names in collection
func findIndex(colletion []string, value string) int {
	for i, v := range colletion {
		if v == value {
			return i
		}
	}
	return -1
}
