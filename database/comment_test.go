package database

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"wrong.wang/x/go-isso/isso"
)

func TestDatabase_NewComment(t *testing.T) {
	type args struct {
		ctx        context.Context
		c          isso.Comment
		threadID   int64
		remoteAddr string
	}
	tests := []struct {
		name    string
		d       *Database
		args    args
		want    isso.Comment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.NewComment(tt.args.ctx, tt.args.c, tt.args.threadID, tt.args.remoteAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.NewComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Database.NewComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GetComment(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		_, err := db.GetComment(context.Background(), -1)
		if !errors.Is(err, isso.ErrStorageNotFound) {
			t.Errorf("Database.getComment() error = %v, wantErr %v", err, isso.ErrStorageNotFound)
			return
		}
	})
}
