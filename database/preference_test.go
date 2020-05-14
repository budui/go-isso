package database

import (
	"errors"
	"testing"

	"wrong.wang/x/go-isso/isso"
)

func TestDatabase_GetPreference(t *testing.T) {
	t.Run("not exist key", func(t *testing.T) {
		value, err := db.GetPreference("not-exist")
		if !errors.Is(err, isso.ErrStorageNotFound) {
			t.Errorf("Database.GetPreference() want isso.ErrStorageNotFound's wrapper, but got %v", err)
		}
		if value != "" {
			t.Errorf("expect null string, but got %s", value)
		}
	})
	t.Run("exist key", func(t *testing.T) {
		err := db.SetPreference("key", "value")
		if err != nil {
			t.Errorf("Database.SetPreference() error = %v, wantErr %v", err, false)
		}
		value, err := db.GetPreference("key")
		if err != nil {
			t.Errorf("Database.GetPreference() error = %v, wantErr %v", err, false)
		}
		if value != "value" {
			t.Errorf("expect `value`, but got %s", value)
		}
	})
}

func TestDatabase_SetPreference(t *testing.T) {
	t.Run("same key", func(t *testing.T) {
		err := db.SetPreference("key1", "value")
		if err != nil {
			t.Errorf("Database.SetPreference() error = %v, wantErr %v", err, false)
		}
		err = db.SetPreference("key1", "value1")
		if err == nil {
			t.Errorf("Database.SetPreference() error = %v, wantErr %v", err, true)
		}
		value, err := db.GetPreference("key1")
		if err != nil {
			t.Errorf("Database.GetPreference() error = %v, wantErr %v", err, false)
		}
		if value != "value" {
			t.Errorf("expect `value`, but got %s", value)
		}
	})
}
