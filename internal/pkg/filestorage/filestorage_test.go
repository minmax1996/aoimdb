package filestorage

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// DatabaseController Database structure copy
type DatabaseController struct {
	Databases    map[string]Database
	Users        map[string]string
	AccessTokens []string
}

type Database struct {
	Name   string
	Sets   []string
	HSets  map[string]string
	Tables map[string][]string
}

func TestBackup(t *testing.T) {
	dbExample := DatabaseController{
		Databases: map[string]Database{
			"ss": {},
		},
	}
	type args struct {
		db DatabaseController
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple", args{dbExample}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eac, err := encodeAndCompress(tt.args.db)
			if err != nil {
				t.Errorf("Backup().Encode error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var dd DatabaseController
			if err := decompressAndDecore(&dd, eac); err != nil {
				t.Errorf("Backup().Decode error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.args.db, dd); diff != "" {
				t.Errorf("Backup() mismatch: {+want;-got}\n\t%s", diff)
				return
			}
		})
	}
}
