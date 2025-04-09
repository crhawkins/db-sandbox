package dbhelper

import (
	"database/sql"
	"fmt"
)

func ConnectPostgreSQL(host string, port int, dbName string, user string, password string) (DB, error) {
	if port == 0 {
		port = 5432
	}

	const driver = "postgres"

	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbName, user, password,
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
