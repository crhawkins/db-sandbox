package dbhelper

import (
	"database/sql"
	"fmt"
	"strings"
)

func NewPostgreSQL(host string, port int, user string, password string) (DB, error) {
	if port == 0 {
		port = 5432
	}

	const driver = "postgres"

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password,
	)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &dbWrapper{
		sqlDB:  db,
		driver: driver,
	}, nil
}

func NewSQLite(filepath string) (DB, error) {
	const driver = "sqlite3"
	if filepath == "" {
		return nil, fmt.Errorf("sqlite: file path is required")
	}
	if !strings.HasSuffix(filepath, ".db") {
		filepath += ".db"
	}

	db, err := sql.Open(driver, filepath)
	if err != nil {
		return nil, err
	}

	return &dbWrapper{
		sqlDB:  db,
		driver: driver,
	}, nil
}

func (d *dbWrapper) Close() error {
	if d.sqlDB != nil {
		return d.sqlDB.Close()
	}
	return nil
}

func (d *dbWrapper) TestConnection() error {
	if d.sqlDB == nil {
		return fmt.Errorf("connection is nil")
	}
	return d.sqlDB.Ping()
}
