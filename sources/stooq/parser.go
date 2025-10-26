package stooq

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
)

// ParsedData represents parsed Stooq CSV data.
type ParsedData struct {
	Columns []string
	Rows    []map[string]string
}

// ParseCSV parses Stooq CSV response data.
func ParseCSV(data []byte) (*ParsedData, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}

	// Validate header
	if len(header) == 0 {
		return nil, fmt.Errorf("empty CSV header")
	}

	// Read all rows
	rows := make([]map[string]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV row: %w", err)
		}

		// Skip rows with wrong number of columns
		if len(record) != len(header) {
			continue
		}

		// Create row map
		row := make(map[string]string)
		for i, value := range record {
			row[header[i]] = value
		}
		rows = append(rows, row)
	}

	// Sort rows by date (ascending)
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i]["Date"] < rows[j]["Date"]
	})

	return &ParsedData{
		Columns: header,
		Rows:    rows,
	}, nil
}
