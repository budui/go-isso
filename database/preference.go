package database

import "fmt"

// GetPreference get preference use key
func (d *Database) GetPreference(key string) (string, error) {
	var value string
	err := d.DB.QueryRow(d.statement["preference_get"], key).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("GetPreference failed. %w", err)
	}
	return value, nil
}

// SetPreference set preference with key value pairs.
func (d *Database) SetPreference(key string, value string) error {
	result, err := d.DB.Exec(d.statement["preference_set"], key, value)
	if err != nil {
		return fmt.Errorf("SetPreference failed. %w", err)
	}
	row, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SetPreference failed. %w", err)
	}
	if row != 1 {
		return ErrNotExpectRow
	}
	return nil
}



