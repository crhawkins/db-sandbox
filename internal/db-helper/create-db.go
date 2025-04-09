package dbhelper

import (
	"fmt"
	"regexp"

	_ "github.com/lib/pq" // import your driver
)

func (d *dbWrapper) Create(name string) error {
	if name == "" {
		return fmt.Errorf("file name is empty")
	}

	if err := validateDatabaseName(name); err != nil {
		return err
	}

	switch d.driver {
	case "postgres":
		_, err := d.sqlDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s";`, name))
		if err != nil {
			return fmt.Errorf("postgres: failed to create database: %w", err)
		}
		return nil

	case "sqlite3":
		// No-op — file-based database already created via sql.Open()
		return nil

	default:
		return fmt.Errorf("unsupported driver: %s", d.driver)
	}
}

func (d *dbWrapper) CreateIfNotExists(name string) error {
	exists, err := d.DatabaseExists(name)
	if err != nil {
		return err
	}
	if !exists {
		return d.Create(name)
	}
	return nil
}

func validateDatabaseName(name string) error {
	if name == "" {
		return fmt.Errorf("database name is empty")
	}
	validName := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid database name: %q", name)
	}
	return nil
}
