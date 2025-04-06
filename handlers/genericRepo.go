package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	
)

func ParseCSV[T any](file io.Reader) ([]T, error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(lines) < 2 {
		return nil, errors.New("CSV file must have at least one data row")
	}

	headers := lines[0]
	var records []T
	seenIDs := make(map[int64]bool)

	for rowIndex, row := range lines[1:] {
		if len(row) != len(headers) {
			return nil, fmt.Errorf("record on line %d: expected %d fields, got %d", rowIndex+2, len(headers), len(row))
		}

		var record T
		val := reflect.ValueOf(&record).Elem()
		var id int64

		for i, field := range headers {
			field = strings.TrimSpace(field)
			fieldVal := strings.TrimSpace(row[i])

			structField := val.FieldByNameFunc(func(name string) bool {
				return strings.EqualFold(name, field)
			})

			if !structField.IsValid() || !structField.CanSet() {
				continue
			}

			switch structField.Kind() {
			case reflect.String:
				structField.SetString(fieldVal)
			case reflect.Int, reflect.Int64:
				intVal, err := strconv.ParseInt(fieldVal, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("record on line %d: invalid integer in field '%s'", rowIndex+2, field)
				}
				structField.SetInt(intVal)
				if strings.EqualFold(field, "jobId") {
					id = intVal
				}
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(fieldVal)
				if err != nil {
					return nil, fmt.Errorf("record on line %d: invalid boolean in field '%s'", rowIndex+2, field)
				}
				structField.SetBool(boolVal)
			}
		}

		if seenIDs[id] {
			continue
		}
		seenIDs[id] = true
		records = append(records, record)
	}

	return records, nil
}
