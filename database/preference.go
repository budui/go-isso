package database

import "wrong.wang/x/go-isso/logger"

// GetPreference get preference use key
func (d *Database) GetPreference(key string) (string, error) {
	logger.Debug("key: %s", key)
	var value string
	err := d.DB.QueryRow(d.statement["preference_get"], key).Scan(&value)
	if err != nil {
		return "", wraperror(err)
	}
	return value, nil
}

// SetPreference set preference with key value pairs.
func (d *Database) SetPreference(key string, value string) error {
	logger.Debug("key: %s, value %s", key, value)
	result, err := d.DB.Exec(d.statement["preference_set"], key, value)
	if err != nil {
		return wraperror(err)
	}
	row, err := result.RowsAffected()
	if err != nil {
		return wraperror(err)
	}
	if row != 1 {
		return wraperror(ErrNotExpectRow)
	}
	return nil
}
