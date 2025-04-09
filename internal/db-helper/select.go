package dbhelper

import (
	"fmt"
	"strings"
)

func CreateSelectSQL(base string, fields []string, where []string, joins []string) string {
	baseTable := strings.ToLower(base)

	// SELECT clause
	var selectParts []string
	for _, field := range fields {
		selectParts = append(selectParts, parseFieldToColumn(baseTable, field))
	}
	selectClause := "SELECT " + strings.Join(selectParts, ", ")

	// FROM clause
	fromClause := fmt.Sprintf("FROM %s", baseTable)

	// JOIN clause (supporting deep paths like Company.Country)
	joinClause := ""
	seen := map[string]bool{}
	for _, joinPath := range joins {
		parts := strings.Split(joinPath, ".")
		prev := baseTable
		for _, part := range parts {
			curr := strings.ToLower(part)
			key := prev + "_" + curr
			if seen[key] {
				prev = curr
				continue
			}
			joinClause += fmt.Sprintf(" LEFT JOIN %s ON %s.%s_id = %s.id", curr, prev, curr, curr)
			seen[key] = true
			prev = curr
		}
	}

	// WHERE clause
	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE " + strings.Join(where, " AND ")
	}

	return fmt.Sprintf("%s %s%s%s;", selectClause, fromClause, joinClause, whereClause)
}

func parseFieldToColumn(base string, field string) string {
	parts := strings.Split(field, ".")
	if len(parts) == 1 {
		return fmt.Sprintf("%s.%s", strings.ToLower(base), strings.ToLower(parts[0]))
	}
	// Handle nested field like Company.Country.Name
	table := strings.ToLower(parts[len(parts)-2])
	column := strings.ToLower(parts[len(parts)-1])
	return fmt.Sprintf("%s.%s", table, column)
}
