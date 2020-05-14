package database

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"wrong.wang/x/go-isso/isso"
)

func TestDatabase_Thread(t *testing.T) {
	t.Run("not exist", func(t *testing.T) {
		got, err := db.GetThreadByURI(context.Background(), "/not-exist")
		if !errors.Is(err, isso.ErrStorageNotFound) {
			t.Errorf("Database.GetThreadByURI() error = %v, wantErr %v", err, isso.ErrStorageNotFound)
			return
		}
		emptyt := isso.Thread{}
		if !reflect.DeepEqual(got, emptyt) {
			t.Errorf("Database.GetThreadByURI() = %v, want %v", got, emptyt)
		}

		got, err = db.GetThreadByID(context.Background(), 1024)
		if !errors.Is(err, isso.ErrStorageNotFound) {
			t.Errorf("Database.GetThreadByID() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, emptyt) {
			t.Errorf("Database.GetThreadByID() = %v, want %v", got, emptyt)
		}
	})
	t.Run("normal", func(t *testing.T) {
		newt1, err := db.NewThread(context.Background(), "/uri", "/hello", "wrong.wang")
		if (err != nil) != false {
			t.Errorf("Database.NewThread() error = %v, wantErr %v", err, false)
			return
		}

		newt2, err := db.NewThread(context.Background(), "/about", "", "https://wrong.wang")
		if (err != nil) != false {
			t.Errorf("Database.NewThread() %v", err)
			return
		}

		got, err := db.GetThreadByURI(context.Background(), "/uri")
		if err != nil {
			t.Errorf("Database.GetThreadByURI() error = %v, wantErr %v", err, false)
			return
		}
		if !reflect.DeepEqual(got, newt1) {
			t.Errorf("Database.GetThreadByID() = %v, want %v", got, newt1)
		}

		got, err = db.GetThreadByID(context.Background(), newt2.ID)
		if err != nil {
			t.Errorf("Database.GetThreadByID() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, newt2) {
			t.Errorf("Database.GetThreadByID() = %v, want %v", got, newt2)
		}
	})
}
