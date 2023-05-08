package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

func Load(db *sql.DB) (Schema, error) {
	var tables []Table
	tableNames, err := tableList(db)
	if err != nil {
		return Schema{}, fmt.Errorf("getting tables: %w", err)
	}
	for _, table := range tableNames {
		tableDetail, err := describeTable(table, db)
		if err != nil {
			return Schema{}, fmt.Errorf("describing table %v: %w", table, err)
		}
		tables = append(tables, tableDetail)
	}
	return Schema{
		Tables: tables,
	}, nil
}

// Schema represents a simplified database schema, containing a list of tables
type Schema struct {
	Tables []Table
}

// String returns the SQL query to create all tables in the schema
func (s Schema) String() string {
	var tables []string
	for _, table := range s.Tables {
		tables = append(tables, table.String())
	}
	return strings.Join(tables, "\n\n")
}

// Table represents a table in a database
type Table struct {
	Name    string
	Columns []Column
}

// String returns the SQL query to create a table with its columns.
func (t Table) String() string {
	var columns []string
	for _, column := range t.Columns {
		columns = append(columns, fmt.Sprintf("%s %s", column.Name, column.Type))
	}

	return fmt.Sprintf("CREATE TABLE %s (%s)", t.Name, strings.Join(columns, ", "))
}

// Column represents a column in a table
type Column struct {
	Name string
	Type string
}

func tableList(db *sql.DB) ([]string, error) {
	var tables []string
	rows, err := db.Query("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';")
	if err != nil {
		return nil, fmt.Errorf("listing tables: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables, nil
}

// describeTable returns a list of Column descriptions for a given table
func describeTable(table string, db *sql.DB) (Table, error) {
	var columns []Column
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1", table)
	if err != nil {
		return Table{}, fmt.Errorf("describing tables: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var column Column
		rows.Scan(&column.Name, &column.Type)
		columns = append(columns, column)
	}
	return Table{
		Name:    table,
		Columns: columns,
	}, nil
}
