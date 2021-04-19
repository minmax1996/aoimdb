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
		{"error select", NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col3"}}, nil},
		{"1 col select",
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col1"}},
			NewTableSchema("", []string{"col1"}, []reflect.Kind{reflect.Int})},
		{"2 col select",
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col1", "col2"}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String})},
		{"3 col select",
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col1", "col2", "col3"}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String})},
		{"1 col select withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1"}},
			NewTableWithRows("", []string{"col1"}, []reflect.Kind{reflect.Int}, []Row{{11}, {22}})},
		{"2 col select  withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1", "col2"}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{11, "str"}, {22, "str2"}})},
		{"3 col select withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1", "col2", "col3"}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{11, "str"}, {22, "str2"}})},
		{"3 col select withRows in reverse order",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col2", "col1"}},
			NewTableWithRows("", []string{"col2", "col1"}, []reflect.Kind{reflect.String, reflect.Int}, []Row{{"str", 11}, {"str2", 22}})},
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
		{"example1", NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col1", "col2"}, []interface{}{12, "abcd"}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{12, "abcd"}}), false},
		{"with reverse order", NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{12, "abcd"}}), false},
		{"with move columns", NewTableSchema("", []string{"col1", "col2", "col3", "col4"}, []reflect.Kind{reflect.Int, reflect.String, reflect.Int, reflect.Int}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2", "col3", "col4"}, []reflect.Kind{reflect.Int, reflect.String, reflect.Int, reflect.Int}, []Row{{12, "abcd", nil, nil}}), false},
		{"with move columns", NewTableSchema("", []string{"col1", "col2", "col3", "col4"}, []reflect.Kind{reflect.Int, reflect.String, reflect.Int, reflect.Int}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2", "col3", "col4"}, []reflect.Kind{reflect.Int, reflect.String, reflect.Int, reflect.Int}, []Row{{12, "abcd", nil, nil}}), false},
		{"ErrInLengths", NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col1", "col2"}, []interface{}{12, "abcd", "324"}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), true},
		{"ErrInNotFoundIndex", NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col1", "col3"}, []interface{}{12, "abcd"}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), true},
		{"ErrInTypeCheck", NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), args{[]string{"col1", "col2"}, []interface{}{12, 321312}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}), true},
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

func BenchmarkInsert(b *testing.B) {
	table := NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String})
	for n := 0; n < b.N; n++ {
		table.Insert([]string{"col1"}, []interface{}{n})
	}
}

func benchmarkFilterDiv(count int, b *testing.B) {
	table := NewTableSchema("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String})
	for n := 0; n < count; n++ {
		table.Insert([]string{"col1"}, []interface{}{n})
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		table.Filter(func(m Row) bool {
			return m[0].(int)%2 == 0
		})
	}
}

func BenchmarkFilterDiv10(b *testing.B)  { benchmarkFilterDiv(10, b) }
func BenchmarkFilterDiv100(b *testing.B) { benchmarkFilterDiv(100, b) }

func TestTable_Filter(t *testing.T) {
	type args struct {
		filterfunc func(Row) bool
	}
	tests := []struct {
		name string
		t    *Table
		args args
		want *Table
	}{
		{"1 filtered",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{9, "str"}, {22, "str2"}}), args{func(m Row) bool {
				return m[0].(int) < 10 && len(m[1].(string)) > 0
			}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{9, "str"}})},
		{"none filtered, panic recover",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{{9, "str"}, {22, "str2"}}), args{func(m Row) bool {
				return m[1].(int) < 10 && len(m[0].(string)) > 0
			}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Kind{reflect.Int, reflect.String}, []Row{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Filter(tt.args.filterfunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
