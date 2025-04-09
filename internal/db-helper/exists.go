package dbhelper

import (
	"fmt"
	"os"
	"strings"
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

	case "sqlite3":
		// In SQLite, the database is a file
		if name == "" {
			return false, fmt.Errorf("sqlite: file name is empty")
		}
		if !strings.HasSuffix(name, ".db") {
			name += ".db"
		}
		_, err := os.Stat(name)
		if os.IsNotExist(err) {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("sqlite: failed to stat file: %w", err)
		}
		return true, nil

	default:
		return false, fmt.Errorf("unsupported driver: %s", d.driver)
	}
}
