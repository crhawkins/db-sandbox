package dbhelper

import (
	"fmt"
	"os"
	"strings"
)

func (d *dbWrapper) Delete(name string) error {
	if name == "" {
		return fmt.Errorf("database name is required")
	}

	if err := validateDatabaseName(name); err != nil {
		return err
	}

	switch d.driver {
	case "postgres":
		// Terminate all existing connections to the database (required before DROP)
		_, err := d.sqlDB.Exec(fmt.Sprintf(`
			REVOKE CONNECT ON DATABASE "%s" FROM public;
			SELECT pg_terminate_backend(pid) 
			FROM pg_stat_activity 
			WHERE datname = '%s' AND pid <> pg_backend_pid();`, name, name))
		if err != nil {
			return fmt.Errorf("postgres: failed to terminate connections: %w", err)
		}

		// Drop the database
		_, err = d.sqlDB.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s";`, name))
		if err != nil {
			return fmt.Errorf("postgres: failed to delete database: %w", err)
		}
		return nil

	case "sqlite3":
		if !strings.HasSuffix(name, ".db") {
			name += ".db"
		}
		if err := os.Remove(name); err != nil {
			if os.IsNotExist(err) {
				return nil // already gone
			}
			return fmt.Errorf("sqlite: failed to delete database file: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("unsupported driver: %s", d.driver)
	}
}
