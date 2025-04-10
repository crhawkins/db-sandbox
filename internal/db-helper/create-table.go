package dbhelper

import (
	"fmt"
	"reflect"
	"strings"
)

func (d *dbWrapper) CreateTable(model any) error {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Struct {
		return fmt.Errorf("CreateTable: provided value is not a struct")
	}

	tableName := strings.ToLower(modelType.Name())
	columns := []string{}
	fkConstraints := []string{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldType := field.Type
		fieldName := strings.ToLower(field.Name)
		tag := parseTagOptions(field.Tag.Get("meta"))

		notNull := ""
		if tag["required"] == "true" || tag["nullable"] == "false" {
			notNull = " NOT NULL"
		}

		defaultValue := ""
		if def, ok := tag["default"]; ok {
			defaultValue = fmt.Sprintf(" DEFAULT %s", def)
		}

		checkConstraint := ""
		if check, ok := tag["check"]; ok {
			checkConstraint = fmt.Sprintf(" CHECK (%s)", check)
		}

		switch fieldType.Kind() {
		case reflect.Int:
			if field.Name == "ID" {
				switch d.driver {
				case "postgres":
					columns = append(columns, "id SERIAL PRIMARY KEY")
				default:
					return fmt.Errorf("unsupported driver: %s", d.driver)
				}
			} else {
				col := fmt.Sprintf("%s INTEGER%s%s%s", fieldName, notNull, defaultValue, checkConstraint)
				if tag["unique"] == "true" {
					col += " UNIQUE"
				}
				columns = append(columns, col)
			}

		case reflect.String:
			colType := "TEXT"
			if tag["max_len"] != "" {
				colType = fmt.Sprintf("VARCHAR(%s)", tag["max_len"])
			}
			col := fmt.Sprintf("%s %s%s%s%s", fieldName, colType, notNull, defaultValue, checkConstraint)
			if tag["unique"] == "true" {
				col += " UNIQUE"
			}
			columns = append(columns, col)

		case reflect.Float64:
			col := fmt.Sprintf("%s DOUBLE PRECISION%s%s%s", fieldName, notNull, defaultValue, checkConstraint)
			if tag["unique"] == "true" {
				col += " UNIQUE"
			}
			columns = append(columns, col)

		case reflect.Struct:
			idField, ok := fieldType.FieldByName("ID")
			if ok && idField.Type.Kind() == reflect.Int {
				refTable := strings.ToLower(fieldType.Name())
				colName := fmt.Sprintf("%s_id", refTable)
				columns = append(columns, fmt.Sprintf("%s INTEGER%s", colName, notNull))

				action := strings.ToUpper(tag["on_delete"])
				if action != "CASCADE" && action != "SET NULL" && action != "SET DEFAULT" {
					action = "RESTRICT"
				}

				fk := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(id) ON DELETE %s", colName, refTable, action)
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

func parseTagOptions(tag string) map[string]string {
	options := map[string]string{}
	pairs := strings.Split(tag, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 {
			options[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		} else {
			// Treat as boolean flag
			options[pair] = "true"
		}
	}
	return options
}
