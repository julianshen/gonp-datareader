package yahoo

import (
	"encoding/csv"
	"errors"
	"io"
)

var (
	// ErrEmptyCSV is returned when CSV data is empty
	ErrEmptyCSV = errors.New("CSV data is empty")
)

// ParsedData represents parsed CSV data from Yahoo Finance.
type ParsedData struct {
	// Columns contains the column names from the CSV header
	Columns []string
	// Rows contains the data rows as maps from column name to value
	Rows []map[string]string
}

// GetColumn returns all values for a given column name.
func (p *ParsedData) GetColumn(name string) []string {
	if p == nil || len(p.Rows) == 0 {
		return nil
	}

	values := make([]string, 0, len(p.Rows))
	for _, row := range p.Rows {
		if val, ok := row[name]; ok {
			values = append(values, val)
		}
	}

	if len(values) == 0 {
		return nil
	}

	return values
}

// ParseCSV parses CSV data from Yahoo Finance.
func ParseCSV(reader io.Reader) (*ParsedData, error) {
	csvReader := csv.NewReader(reader)

	// Read all records
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Check for empty data
	if len(records) == 0 {
		return nil, ErrEmptyCSV
	}

	// First row is the header
	header := records[0]
	if len(header) == 0 {
		return nil, ErrEmptyCSV
	}

	// Parse data rows
	rows := make([]map[string]string, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		record := records[i]
		row := make(map[string]string)

		for j, value := range record {
			if j < len(header) {
				row[header[j]] = value
			}
		}

		rows = append(rows, row)
	}

	return &ParsedData{
		Columns: header,
		Rows:    rows,
	}, nil
}
