package dbhelper

import (
	"fmt"
	"reflect"
	"strings"
)

func (d *dbWrapper) CreateTable(model any) error {
	tType := reflect.TypeOf(model)

	if tType.Kind() != reflect.Struct {
		return fmt.Errorf("CreateTable: provided value is not a struct")
	}

	tableName := strings.ToLower(tType.Name())
	columns := []string{}
	fkConstraints := []string{}

	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		fieldType := field.Type
		fieldName := strings.ToLower(field.Name)

		switch fieldType.Kind() {
		case reflect.Int:
			if field.Name == "ID" {
				switch d.driver {
				case "postgres":
					columns = append(columns, "id SERIAL PRIMARY KEY")
				case "sqlite3":
					columns = append(columns, "id INTEGER PRIMARY KEY AUTOINCREMENT")
				default:
					return fmt.Errorf("unsupported driver: %s", d.driver)
				}
			} else {
				columns = append(columns, fmt.Sprintf("%s INTEGER", fieldName))
			}

		case reflect.String:
			columns = append(columns, fmt.Sprintf("%s TEXT", fieldName))

		case reflect.Float64:
			columns = append(columns, fmt.Sprintf("%s DOUBLE PRECISION", fieldName))

		case reflect.Struct:
			// Check if it has an ID field → assume it's a foreign key
			idField, ok := fieldType.FieldByName("ID")
			if ok && idField.Type.Kind() == reflect.Int {
				refTable := strings.ToLower(fieldType.Name())
				colName := fmt.Sprintf("%s_id", refTable)
				columns = append(columns, fmt.Sprintf("%s INTEGER", colName))
				fk := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(id)", colName, refTable)
				fkConstraints = append(fkConstraints, fk)
			} else {
				return fmt.Errorf("unsupported embedded struct field: %s", field.Name)
			}

		default:
			return fmt.Errorf("unsupported field type for %s: %s", field.Name, fieldType.Name())
		}
	}

	allColumns := append(columns, fkConstraints...)
	createStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (%s);`, tableName, strings.Join(allColumns, ", "))

	fmt.Println(createStmt)

	_, err := d.sqlDB.Exec(createStmt)
	if err != nil {
		return fmt.Errorf("failed to create table %q: %w", tableName, err)
	}

	return nil
}
