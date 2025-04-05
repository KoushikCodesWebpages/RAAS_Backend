package controllers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// parseCSV - Parses CSV data into a slice of the given model type while ignoring duplicates
func parseCSV[T any](file io.Reader) ([]T, error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // Allows handling variable-length rows

	// Read all lines
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(lines) < 2 {
		return nil, errors.New("CSV file must have at least one data row")
	}

	// Get headers
	headers := lines[0]
	var records []T
	seenJobIDs := make(map[int64]bool) // Track seen job IDs to prevent duplicates

	// Iterate through each row
	for rowIndex, row := range lines[1:] {
		if len(row) != len(headers) {
			return nil, fmt.Errorf("record on line %d: expected %d fields, got %d", rowIndex+2, len(headers), len(row))
		}

		var record T
		val := reflect.ValueOf(&record).Elem()
		var jobID int64

		for i, field := range headers {
			field = strings.TrimSpace(field)
			fieldVal := strings.TrimSpace(row[i])

			// Case-insensitive struct field lookup
			structField := val.FieldByNameFunc(func(name string) bool {
				return strings.EqualFold(name, field)
			})

			if !structField.IsValid() || !structField.CanSet() {
				fmt.Printf("Skipping field: %s (Not found in struct)\n", field)
				continue
			}

			switch structField.Kind() {
			case reflect.String:
				structField.SetString(fieldVal)

			case reflect.Int, reflect.Int64:
				intVal, err := strconv.ParseInt(fieldVal, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("record on line %d: invalid integer in field '%s' (value: %s)", rowIndex+2, field, fieldVal)
				}
				structField.SetInt(intVal)

				// Store jobId separately for duplicate check
				if strings.EqualFold(field, "jobId") {
					jobID = intVal
				}

			case reflect.Bool:
				boolVal, err := strconv.ParseBool(fieldVal)
				if err != nil {
					return nil, fmt.Errorf("record on line %d: invalid boolean in field '%s' (value: %s)", rowIndex+2, field, fieldVal)
				}
				structField.SetBool(boolVal)
			}
		}

		// Skip duplicate jobIds
		if seenJobIDs[jobID] {
			fmt.Printf("Skipping duplicate jobId: %d\n", jobID)
			continue
		}
		seenJobIDs[jobID] = true // Mark jobId as seen

		fmt.Printf("Parsed record: %+v\n", record)
		records = append(records, record)
	}

	return records, nil
}
