package dbhelper

import "database/sql"

type DB interface {
	Close() error
	Create(name string) error
	CreateIfNotExists(name string) error
	DatabaseExists(name string) (bool, error)
	Raw(ping bool) (*sql.DB, error)
	Version() (string, error)
	TestConnection() error
}

type dbWrapper struct {
	sqlDB  *sql.DB
	driver string
}
