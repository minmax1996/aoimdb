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
		{"error select", NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col3"}}, nil},
		{"1 col select",
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col1"}},
			NewTableSchema("", []string{"col1"}, []reflect.Type{reflect.TypeOf(1)})},
		{"2 col select",
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col1", "col2"}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")})},
		{"3 col select",
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col1", "col2", "col3"}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")})},
		{"1 col select withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1"}},
			NewTableWithRows("", []string{"col1"}, []reflect.Type{reflect.TypeOf(1)}, []Row{{11}, {22}})},
		{"2 col select  withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1", "col2"}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{11, "str"}, {22, "str2"}})},
		{"3 col select withRows",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1", "col2", "col3"}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{11, "str"}, {22, "str2"}})},
		{"3 col select withRows in reverse order",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{11, "str"}, {22, "str2"}}), args{[]string{"col2", "col1"}},
			NewTableWithRows("", []string{"col2", "col1"}, []reflect.Type{reflect.TypeOf(""), reflect.TypeOf(1)}, []Row{{"str", 11}, {"str2", 22}})},
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
		{"example1", NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col1", "col2"}, []interface{}{12, "abcd"}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{12, "abcd"}}), false},
		{"with reverse order", NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{12, "abcd"}}), false},
		{"with move columns", NewTableSchema("", []string{"col1", "col2", "col3", "col4"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf(""), reflect.TypeOf(1), reflect.TypeOf(1)}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2", "col3", "col4"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf(""), reflect.TypeOf(1), reflect.TypeOf(1)}, []Row{{12, "abcd", nil, nil}}), false},
		{"with move columns", NewTableSchema("", []string{"col1", "col2", "col3", "col4"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf(""), reflect.TypeOf(1), reflect.TypeOf(1)}), args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			NewTableWithRows("", []string{"col1", "col2", "col3", "col4"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf(""), reflect.TypeOf(1), reflect.TypeOf(1)}, []Row{{12, "abcd", nil, nil}}), false},
		{"ErrInLengths", NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col1", "col2"}, []interface{}{12, "abcd", "324"}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), true},
		{"ErrInNotFoundIndex", NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col1", "col3"}, []interface{}{12, "abcd"}},
			NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), true},
		{"exampleFromAllStrings", NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}), args{[]string{"col1", "col2"}, []interface{}{"12", "abcd"}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{12, "abcd"}}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Insert(tt.args.names, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("Table.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("Table.NotEqual() = %v, want %v", tt.t, tt.want)
				t.Errorf("%T %T", tt.t.DataRows[0][0], tt.want.DataRows[0][0])
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	table := NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")})
	for n := 0; n < b.N; n++ {
		table.Insert([]string{"col1"}, []interface{}{n})
	}
}

func benchmarkFilterDiv(count int, b *testing.B) {
	table := NewTableSchema("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")})
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
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{9, "str"}, {22, "str2"}}), args{func(m Row) bool {
				return m[0].(int) < 10 && len(m[1].(string)) > 0
			}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{9, "str"}})},
		{"none filtered, panic recover",
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{{9, "str"}, {22, "str2"}}), args{func(m Row) bool {
				return m[1].(int) < 10 && len(m[0].(string)) > 0
			}},
			NewTableWithRows("", []string{"col1", "col2"}, []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}, []Row{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Filter(tt.args.filterfunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
