package db

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gopkg.in/guregu/null.v3"
)

// preferenceAccessor defines all usual access ops avail for comment.
type preferenceAccessor interface {
	GetPreference(key string) (null.String, error)
	SetPreference(key, value string) error
}

func initPreference(db *database) error {
	if v, err := db.GetPreference("session-key"); err != nil || v.Valid {
		return err
	}

	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	sessionKey := hex.EncodeToString(b)

	err = db.SetPreference("session-key", sessionKey)
	if err != nil {
		return err
	}
	return nil
}

// GetPreference get the corresponding value of the key.
func (db *database) GetPreference(key string) (null.String, error) {
	var value null.String
	err := db.QueryRow("SELECT value FROM preferences WHERE key=?", key).Scan(&value)
	if err != nil {
		return null.NewString("", false), fmt.Errorf("db.GetPreference: %v", err)
	}
	return value, nil
}

// SetPreference set the corresponding value of the key.
func (db *database) SetPreference(key, value string) error {
	stmt, err := db.Prepare("INSERT INTO preferences (key, value) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("db.SetPreference: %v", err)
	}
	_, err = stmt.Exec(key, value)
	if err != nil {
		return fmt.Errorf("db.SetPreference: %v", err)
	}
	return nil
}
