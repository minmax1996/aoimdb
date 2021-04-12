package datatypes

import (
	"reflect"
	"testing"
)

func TestTable_Select(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name string
		t    *Table
		args args
		want *Table
	}{
		{"error select", NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col3"}}, nil},
		{"1 col select",
			NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col1"}},
			NewTableSchema("", []string{"col1"}, []ColumnType{IntType})},
		{"2 col select",
			NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col1", "col2"}},
			NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType})},
		{"3 col select",
			NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col1", "col2", "col3"}},
			NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType})},
		{"1 col select withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1"}},
			NewTableWithRows("", []string{"col1"}, []ColumnType{IntType}, []Row{{11}, {22}})},
		{"2 col select  withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1", "col2"}},
			NewTableWithRows("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}, []Row{{11, "str"}, {22, "str2"}})},
		{"3 col select withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1", "col2", "col3"}},
			NewTableWithRows("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}, []Row{{11, "str"}, {22, "str2"}})},
		{"3 col select withRows in reverse order",
			NewTableWithRows("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col2", "col1"}},
			NewTableWithRows("", []string{"col2", "col1"}, []ColumnType{StringType, IntType}, []Row{{"str", 11}, {"str2", 22}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Select(tt.args.names); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_Insert(t *testing.T) {
	type args struct {
		names  []string
		values []interface{}
	}
	tests := []struct {
		name    string
		t       *Table
		args    args
		want    *Table
		wantErr bool
	}{
		{"example1", NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col1", "col2"}, []interface{}{12, "abcd"}},
			NewTableWithRows("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}, []Row{{12, "abcd"}}), false},
		{"with reverse order", NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}, []Row{{12, "abcd"}}), false},
		{"with move columns", NewTableSchema("", []string{"col1", "col2", "col3", "col4"}, []ColumnType{IntType, StringType, IntType, IntType}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2", "col3", "col4"}, []ColumnType{IntType, StringType, IntType, IntType}, []Row{{12, "abcd", nil, nil}}), false},
		{"with move columns", NewTableSchema("", []string{"col1", "col2", "col3", "col4"}, []ColumnType{IntType, StringType, IntType, IntType}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2", "col3", "col4"}, []ColumnType{IntType, StringType, IntType, IntType}, []Row{{12, "abcd", nil, nil}}), false},
		{"ErrInLengths", NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col1", "col2"}, []interface{}{12, "abcd", "324"}},
			NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), true},
		{"ErrInNotFoundIndex", NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col1", "col3"}, []interface{}{12, "abcd"}},
			NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), true},
		{"ErrInTypeCheck", NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), args{[]string{"col1", "col2"}, []interface{}{12, 321312}},
			NewTableSchema("", []string{"col1", "col2"}, []ColumnType{IntType, StringType}), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Insert(tt.args.names, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("Table.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("Table.NotEqual() = %v, want %v", tt.t, tt.want)
			}
		})
	}
}
