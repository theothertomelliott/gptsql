package schema

import (
	"database/sql"
	"fmt"
	"time"
)

type postgresLoader struct{}

func (p *postgresLoader) tableList(db *sql.DB) ([]string, error) {
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
func (p *postgresLoader) describeTable(table string, db *sql.DB) (Table, error) {
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

	exampleRow, err := db.Query(fmt.Sprintf("SELECT * FROM %v LIMIT 1;", table))
	if err != nil {
		return Table{}, fmt.Errorf("loading example row: %w", err)
	}
	defer exampleRow.Close()
	var values []interface{}
	// Read exactly one row, assuming there are any
	for exampleRow.Next() && len(values) == 0 {
		for i := 0; i < len(columns); i++ {
			var value interface{}
			values = append(values, &value)
		}
		err = exampleRow.Scan(values...)
		if err != nil {
			return Table{}, fmt.Errorf("scanning example row: %w", err)
		}
	}

	var stringValues []string
	for _, rvp := range values {
		// Based on code from sqltocsv
		rawValue := *rvp.(*interface{})
		var value interface{}
		byteArray, ok := rawValue.([]byte)
		if ok {
			value = string(byteArray)
		} else {
			value = rawValue
		}

		float64Value, ok := value.(float64)
		if ok {
			value = fmt.Sprintf("%v", float64Value)
		} else {
			float32Value, ok := value.(float32)
			if ok {
				value = fmt.Sprintf("%v", float32Value)
			}
		}

		timeValue, ok := value.(time.Time)
		if ok {
			value = timeValue.Format(time.RFC1123)
		}

		if value == nil {
			stringValues = append(stringValues, "")
		} else {
			stringValues = append(stringValues, fmt.Sprintf("%v", value))
		}
	}

	return Table{
		Name:      table,
		Columns:   columns,
		SampleRow: stringValues,
	}, nil
}
