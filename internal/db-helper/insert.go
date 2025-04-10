package dbhelper

import (
	"fmt"
	"reflect"
	"strings"
)

func (d *dbWrapper) Insert(model any) (int, error) {
	tValue := reflect.ValueOf(model)
	tType := reflect.TypeOf(model)

	if tType.Kind() != reflect.Struct {
		return 0, fmt.Errorf("Insert: expected a struct, got %s", tType.Kind())
	}

	tableName := strings.ToLower(tType.Name())
	columns := []string{}
	placeholders := []string{}
	values := []any{}

	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		value := tValue.Field(i)

		// Skip ID (auto-increment)
		if field.Name == "ID" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String, reflect.Int:
			columns = append(columns, strings.ToLower(field.Name))
			values = append(values, value.Interface())

		case reflect.Float64:
			columns = append(columns, strings.ToLower(field.Name))
			values = append(values, value.Interface())

		case reflect.Struct:
			idField := value.FieldByName("ID")
			if idField.IsValid() && idField.Kind() == reflect.Int {
				refColumn := strings.ToLower(field.Type.Name()) + "_id"
				columns = append(columns, refColumn)
				values = append(values, idField.Interface())
			} else {
				return 0, fmt.Errorf("Insert: embedded struct %s does not contain an ID", field.Name)
			}

		default:
			return 0, fmt.Errorf("Insert: unsupported field type for %s", field.Name)
		}

		placeholders = append(placeholders, "?")
	}

	columnStr := strings.Join(columns, ", ")
	placeholderStr := strings.Join(placeholders, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnStr, placeholderStr)

	if d.driver == "postgres" {
		for i := range placeholders {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
		placeholderStr = strings.Join(placeholders, ", ")
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", tableName, columnStr, placeholderStr)

		var newID int
		err := d.sqlDB.QueryRow(query, values...).Scan(&newID)
		if err != nil {
			return 0, fmt.Errorf("Insert: failed to insert into %s: %w", tableName, err)
		}
		return newID, nil
	}

	result, err := d.sqlDB.Exec(query, values...)
	if err != nil {
		return 0, fmt.Errorf("Insert: failed to insert into %s: %w", tableName, err)
	}

	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}
