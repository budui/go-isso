package sqlite3

import (
	"testing"

	"github.com/RayHY/go-isso/internal/pkg/conf"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreateDatabase(t *testing.T) {
	type args struct {
		sqliteConf conf.Sqlite3
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"mem", args{conf.Sqlite3{":memory:"}}, false},
		{"comments", args{conf.Sqlite3{"comments_test.db"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDatabase(tt.args.sqliteConf)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Errorf("no error happen but got nil database")
			}
		})
	}
}
