package schema

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func Load(dbType string, db *sql.DB) (Schema, error) {

	var loader loader
	switch dbType {
	case "postgres":
		loader = &postgresLoader{}
	case "snowflake":
		loader = &snowflakeLoader{}
	default:
		return Schema{}, fmt.Errorf("unsupported database type %v", dbType)
	}

	var tables []Table
	tableNames, err := loader.tableList(db)
	if err != nil {
		return Schema{}, fmt.Errorf("getting tables: %w", err)
	}

	log.Printf("Got %v tables", len(tableNames))

	for _, table := range tableNames {
		tableDetail, err := loader.describeTable(table, db)
		if err != nil {
			return Schema{}, fmt.Errorf("describing table %v: %w", table, err)
		}
		tables = append(tables, tableDetail)
	}
	return Schema{
		Tables: tables,
	}, nil
}

type loader interface {
	tableList(db *sql.DB) ([]string, error)
	describeTable(table string, db *sql.DB) (Table, error)
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
	Name      string
	Columns   []Column
	SampleRow []string
}

// String returns the SQL query to create a table with its columns.
func (t Table) String() string {
	var columns []string
	for _, column := range t.Columns {
		columns = append(columns, fmt.Sprintf("%s %s", column.Name, column.Type))
	}

	return fmt.Sprintf(`CREATE TABLE %s (%s)
INSERT INTO %s VALUES (%s);`, t.Name, strings.Join(columns, ", "), t.Name, strings.Join(t.SampleRow, ", "))
}

// Column represents a column in a table
type Column struct {
	Name string
	Type string
}
