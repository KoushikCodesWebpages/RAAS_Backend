package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

// MarshalData takes any Go object and returns a JSON string
func MarshalData(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshaling data: %w", err)
	}
	return string(bytes), nil
}

// UnmarshalData takes a JSON string and populates the target interface
func UnmarshalData[T any](jsonStr string) (T, error) {
	var target T
	if jsonStr == "" {
		return target, errors.New("cannot unmarshal from empty string")
	}
	err := json.Unmarshal([]byte(jsonStr), &target)
	if err != nil {
		return target, fmt.Errorf("error unmarshaling data: %w", err)
	}
	return target, nil
}
