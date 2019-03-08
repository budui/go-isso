package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/RayHY/go-isso/internal/pkg/conf"
	"github.com/RayHY/go-isso/internal/pkg/db/sqlite3"
)

type database struct {
	*sql.DB
	guard conf.Guard
}

// AcceptableError used to type
type AcceptableError struct {
	originErr error
}

func (e *AcceptableError) Error() string {
	return e.originErr.Error()
}

func newError(errString string) error {
	return &AcceptableError{originErr: errors.New(errString)}
}

// NewAccessor generate an new DB worker
func NewAccessor(dbConfig conf.Database, guard conf.Guard) (Accessor, error) {
	var db *sql.DB
	var err error

	switch strings.ToLower(dbConfig.Dialect) {
	case "sqlite3":
		db, err = sqlite3.CreateDatabase(dbConfig.Sqlite3)
	default:
		err = fmt.Errorf("unsupported dialect %v", dbConfig.Dialect)
		db = nil
	}
	if err != nil {
		return nil, fmt.Errorf("NewAccessor failed. [%s]", err.Error())
	}

	dbw := &database{db, guard}
	_ = initPreference(dbw)
	return dbw, nil
}

func (db *database) Close() error {
	return db.Close()
}
