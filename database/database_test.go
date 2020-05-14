package database

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"wrong.wang/x/go-isso/isso"
)

var db *Database

func TestMain(m *testing.M) {
	var err error
	db, err = New("", 1*time.Second)
	if err != nil {
		log.Fatal("init database failed.")
	}
	code := m.Run()

	db.DB.Close()
	os.Exit(code)
}

func Test_databaseError_Format(t *testing.T) {
	t.Run("format", func(t *testing.T) {
		err := wraperror(sql.ErrNoRows)
		if err.Unwrap() != isso.ErrStorageNotFound {
			t.Errorf("wrapper return not wrapped isso error")
		}

		t.Logf("%v", err)
		t.Logf("%+s", err)
	})
}
