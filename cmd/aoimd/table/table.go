package table

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/minmax1996/aoimdb/cmd/aoimd/datatypes"
)

type Rows []Row

func (rr *Rows) Export() [][]string {
	if rr == nil {
		return [][]string{}
	}
	result := make([][]string, 0, len(*rr))
	for _, row := range *rr {
		result = append(result, row.Export())
	}
	return result
}

type Row []interface{}

func (r *Row) Export() []string {
	if r == nil {
		return []string{}
	}
	result := make([]string, 0, len(*r))
	for _, el := range *r {
		result = append(result, fmt.Sprintf("%v", el))
	}
	return result
}

type Table struct {
	sync.RWMutex
	Name        string
	ColumnNames []string
	ColumnTypes []datatypes.Datatype
	DataRows    Rows
}

// NewTableSchema initialize new table scheam with no data
func NewTableSchema(tableName string, names []string, types []datatypes.Datatype) *Table {
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

// NewTableWithRows initailexe new table with rows
func NewTableWithRows(tableName string, names []string, types []datatypes.Datatype, rows []Row) *Table {
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

// Insert inserts into rows new row with values if values pass typecheck
func (t *Table) Insert(names []string, values []interface{}) error {
	var err error
	if len(names) != len(values) {
		return errors.New("lens not equal")
	}
	t.Lock()
	defer t.Unlock()

	row := make(Row, len(t.ColumnNames))
	for i, val := range values {
		ind := findIndex(t.ColumnNames, names[i])
		if ind == -1 {
			return errors.New("cant find index")
		}

		if reflect.TypeOf(val) != t.ColumnTypes[ind].Type {
			row[ind], err = datatypes.ConvertToColumnType(val, t.ColumnTypes[ind].Kind())
			if err != nil {
				return errors.Unwrap(fmt.Errorf("failed to convert to table type %w", err))
			}
		} else {
			row[ind] = val
		}
	}
	t.DataRows = append(t.DataRows, row)
	//TODO insert indexes here
	return nil
}

// Filter filters dataRows by given function
func (t *Table) Filter(filterfunc func(Row) bool) (resultTable *Table) {
	//we using named return value to properly return resultTable after recover from panic
	t.RLock()
	defer t.RUnlock()
	defer func() {
		_ = recover()
	}()

	resultTable = &Table{
		ColumnNames: t.ColumnNames,
		ColumnTypes: t.ColumnTypes,
		DataRows:    make([]Row, 0, len(t.DataRows)),
	}
	//TODO index
	//filterfunc will panic if in user input has wrong type convertions
	for _, v := range t.DataRows {
		//filerMap := v.GetMap(t.ColumnNames)
		if filterfunc(v) {
			resultTable.DataRows = append(resultTable.DataRows, v)
		}
	}

	return resultTable
}

// Delete not sure leave it this way or change later
func (t *Table) Delete(filterfunc func(Row) bool) (err error) {
	//we using named return value to properly return resultTable after recover from panic
	defer func() {
		_ = recover()
		err = errors.New("panic recovered")
	}()
	t.Lock()
	defer t.Unlock()
	newDataRows := make([]Row, 0)
	//filterfunc will panic if in user input has wrong type convertion
	for _, v := range t.DataRows {
		//TODO index
		if !filterfunc(v) {
			newDataRows = append(newDataRows, v)
		}
	}
	//TODO rewrite index
	t.DataRows = newDataRows
	return nil
}

// Select returns filtered table with only given names in given order
func (t *Table) Select(names []string) *Table {
	indexes := findIndexes(t.ColumnNames, names)
	if len(indexes) == 0 {
		return nil
	}
	t.RLock()
	defer t.RUnlock()
	//make empty table with fixed lens
	resultTable := &Table{
		ColumnNames: make([]string, len(indexes)),
		ColumnTypes: make([]datatypes.Datatype, len(indexes)),
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

// findIndexes finds all indexes that correspond all columns names in collection
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

// findIndexes finds all indexes that correspond all columns names in collection
func findIndex(colletion []string, value string) int {
	for i, v := range colletion {
		if v == value {
			return i
		}
	}
	return -1
}
