package db

import (
	"database/sql"
	"fmt"

	"gopkg.in/guregu/null.v3"
)

// Thread is associated with Comments
type Thread struct {
	ID    int64
	URI   string
	Title null.String
}

func (db *database) NewThread(uri string, title null.String) (Thread, error) {
	stmt, err := db.Prepare(`INSERT INTO threads (uri, title) VALUES (?, ?)`)
	if err != nil {
		return Thread{}, fmt.Errorf("NewThread failed %s", err.Error())
	}
	res, err := stmt.Exec(uri, title)
	if err != nil {
		return Thread{}, fmt.Errorf("NewThread failed %s", err.Error())
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return Thread{}, fmt.Errorf("NewThread failed %s", err.Error())
	}
	return Thread{ID: lastID, URI: uri, Title: title}, nil
}

func (db *database) GetThreadWithID(ID int64) (Thread, error) {
	var th Thread
	err := db.QueryRow(`SELECT * FROM threads WHERE id=?`, ID).Scan(
		&th.ID, &th.URI, &th.Title,
	)
	if err != nil {
		return Thread{}, err
	}
	return th, nil
}

func (db *database) GetThreadWithURI(uri string) (Thread, error) {
	var th Thread
	err := db.QueryRow(`SELECT * FROM threads WHERE uri=?`, uri).Scan(
		&th.ID, &th.URI, &th.Title,
	)
	if err != nil {
		return Thread{}, err
	}
	return th, nil
}

func (db *database) Contains(uri string) (bool, error) {
	var title string

	err := db.QueryRow(`SELECT title FROM threads WHERE uri=?`, uri).Scan(
		&title,
	)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
