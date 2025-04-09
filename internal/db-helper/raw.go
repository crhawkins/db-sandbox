package dbhelper

import (
	"database/sql"
	"fmt"
)

func (d *dbWrapper) Raw(ping bool) (*sql.DB, error) {
	if d.sqlDB == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	if ping {
		if err := d.TestConnection(); err != nil {
			return nil, err
		}
	}

	return d.sqlDB, nil
}
