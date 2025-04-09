package dbhelper

import "database/sql"

type DB interface {
	Close() error
	CreateTable(model any) error
	DatabaseExists(name string) (bool, error)
	Delete(name string) error
	Insert(model any) (int, error)
	Raw(ping bool) (*sql.DB, error)

	Version() (string, error)
	TestConnection() error
}

type dbWrapper struct {
	sqlDB  *sql.DB
	driver string
}
