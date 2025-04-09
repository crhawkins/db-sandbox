package dbhelper

import (
	"fmt"
	"regexp"

	_ "github.com/lib/pq" // import your driver
)

func CreatePostgreSQL(host string, port int, dbName string, user string, password string) error {
	if dbName == "" {
		return fmt.Errorf("file name is empty")
	}

	if err := validateDatabaseName(dbName); err != nil {
		return err
	}

	db, err := ConnectPostgreSQL(host, port, "postgres", user, password)
	if err != nil {
		return err
	}

	raw, err := db.Raw(false)
	if err != nil {
		return fmt.Errorf("failed to access raw DB: %w", err)
	}
	_, err = raw.Exec(fmt.Sprintf(`CREATE DATABASE "%s";`, dbName))
	if err != nil {
		return fmt.Errorf("postgres: failed to create database: %w", err)
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
