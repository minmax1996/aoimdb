package table_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	dt "github.com/minmax1996/aoimdb/internal/aoimdb/datatypes"
	"github.com/minmax1996/aoimdb/internal/aoimdb/table"
)

func TestTable_Select(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name string
		t    *table.Table
		args args
		want *table.Table
	}{
		{"error select", table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}), args{[]string{"col3"}}, nil},
		{"1 col select",
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}), args{[]string{"col1"}},
			table.NewTableSchema("", []string{"col1"}, []dt.Datatype{dt.Int})},
		{"2 col select",
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}), args{[]string{"col1", "col2"}},
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String})},
		{"3 col select",
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}), args{[]string{"col1", "col2", "col3"}},
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String})},
		{"1 col select withRows",
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1"}},
			table.NewTableWithRows("", []string{"col1"}, []dt.Datatype{dt.Int}, []table.Row{{11}, {22}})},
		{"2 col select  withRows",
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1", "col2"}},
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{11, "str"}, {22, "str2"}})},
		{"3 col select withRows",
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{11, "str"}, {22, "str2"}}), args{[]string{"col1", "col2", "col3"}},
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{11, "str"}, {22, "str2"}})},
		{"3 col select withRows in reverse order",
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{11, "str"}, {22, "str2"}}), args{[]string{"col2", "col1"}},
			table.NewTableWithRows("", []string{"col2", "col1"}, []dt.Datatype{dt.String, dt.Int}, []table.Row{{"str", 11}, {"str2", 22}})},
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
	tests := map[string]struct {
		t       *table.Table
		args    args
		want    *table.Table
		wantErr bool
	}{

		"example1": {
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}),
			args{[]string{"col1", "col2"}, []interface{}{12, "abcd"}},
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String},
				[]table.Row{{12, "abcd"}}), false,
		},
		"with reverse order": {
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}),
			args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String},
				[]table.Row{{12, "abcd"}}), false,
		},
		"with move columns": {
			table.NewTableSchema("", []string{"col1", "col2", "col3", "col4"}, []dt.Datatype{dt.Int, dt.String, dt.Int, dt.Int}),
			args{[]string{"col2", "col1"}, []interface{}{"abcd", 12}},
			table.NewTableWithRows("", []string{"col1", "col2", "col3", "col4"}, []dt.Datatype{dt.Int, dt.String, dt.Int, dt.Int},
				[]table.Row{{12, "abcd", nil, nil}}), false,
		},
		"ErrInLengths": {
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}),
			args{[]string{"col1", "col2"}, []interface{}{12, "abcd", "324"}},
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}),
			true,
		},
		"ErrInNotFoundIndex": {
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}),
			args{[]string{"col1", "col3"}, []interface{}{12, "abcd"}},
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}),
			true,
		},
		"exampleFromAllStrings": {
			table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}),
			args{[]string{"col1", "col2"}, []interface{}{"12", "abcd"}},
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String},
				[]table.Row{{12, "abcd"}}), false,
		},
	}
	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			if err := tt.t.Insert(tt.args.names, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("Table.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want.DataRows, tt.t.DataRows); diff != "" {
				t.Errorf("Table.Insert() mismatch: {+want;-got}\n\t%s", diff)
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	table := table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String})
	for n := 0; n < b.N; n++ {
		_ = table.Insert([]string{"col1"}, []interface{}{n})
	}
}

func benchmarkFilterDiv(count int, b *testing.B) {
	ttable := table.NewTableSchema("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String})
	for n := 0; n < count; n++ {
		_ = ttable.Insert([]string{"col1"}, []interface{}{n})
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = ttable.Filter(func(m table.Row) bool {
			return m[0].(int)%2 == 0
		})
	}
}

func BenchmarkFilterDiv10(b *testing.B)  { benchmarkFilterDiv(10, b) }
func BenchmarkFilterDiv100(b *testing.B) { benchmarkFilterDiv(100, b) }

func TestTable_Filter(t *testing.T) {
	type args struct {
		filterfunc func(table.Row) bool
	}
	tests := []struct {
		name string
		t    *table.Table
		args args
		want *table.Table
	}{
		{"1 filtered",
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{9, "str"}, {22, "str2"}}), args{func(m table.Row) bool {
				return m[0].(int) < 10 && len(m[1].(string)) > 0
			}},
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{9, "str"}})},
		{"none filtered, panic recover",
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{{9, "str"}, {22, "str2"}}), args{func(m table.Row) bool {
				return m[1].(int) < 10 && len(m[0].(string)) > 0
			}},
			table.NewTableWithRows("", []string{"col1", "col2"}, []dt.Datatype{dt.Int, dt.String}, []table.Row{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Filter(tt.args.filterfunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
