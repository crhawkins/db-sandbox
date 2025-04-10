package dbhelper

import "fmt"

func (d *dbWrapper) Version() (string, error) {
	var versionQuery string

	switch d.driver {
	case "postgres":
		versionQuery = "SELECT version();"
	default:
		return "", fmt.Errorf("unsupported driver: %s", d.driver)
	}

	var version string
	err := d.sqlDB.QueryRow(versionQuery).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version: %w", err)
	}

	return version, nil
}
