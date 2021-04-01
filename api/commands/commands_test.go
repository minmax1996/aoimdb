package commands

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {

	RegisterCommand(NewAuthCommand(nil))
	RegisterCommand(NewSelectCommand(nil))
	RegisterCommand(NewGetCommand(nil))
	RegisterCommand(NewSetCommand(nil))
	RegisterCommand(NewKeysCommand(nil))
	RegisterCommand(NewExitCommand(nil))

	type args struct {
		input string
		sep   string
	}
	tests := []struct {
		name    string
		args    args
		want    Commander
		want1   []string
		wantErr bool
	}{
		{"auth", args{"auth admin pass", " "}, GetCommand("auth"), []string{"admin", "pass"}, false},
		{"select", args{"select 3", " "}, GetCommand("select"), []string{"3"}, false},
		{"get", args{"get key", " "}, GetCommand("get"), []string{"key"}, false},
		{"get2", args{"get 2.key", " "}, GetCommand("get"), []string{"2.key"}, false},
		{"set", args{"set 2.key 231", " "}, GetCommand("set"), []string{"2.key", "231"}, false},
		{"keys", args{"keys", " "}, GetCommand("keys"), []string{}, false},
		{"keys", args{"keys 2.ddd", " "}, GetCommand("keys"), []string{"2.ddd"}, false},
		{"authErr", args{"auth admin", " "}, nil, nil, true},
		{"getErr", args{"get", " "}, nil, nil, true},
		{"selectErr", args{"select", " "}, nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseCommand(tt.args.input, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCommand() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseCommand() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
