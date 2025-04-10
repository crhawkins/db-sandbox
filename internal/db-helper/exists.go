package dbhelper

import (
	"fmt"
)

func (d *dbWrapper) DatabaseExists(name string) (bool, error) {
	switch d.driver {
	case "postgres":
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);`
		err := d.sqlDB.QueryRow(query, name).Scan(&exists)
		if err != nil {
			return false, fmt.Errorf("postgres: failed to check database existence: %w", err)
		}
		return exists, nil

	default:
		return false, fmt.Errorf("unsupported driver: %s", d.driver)
	}
}
