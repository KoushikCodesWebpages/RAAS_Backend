package utils

import (
	"time"
	"strings"
	"encoding/json"

)

// DateOnly is a custom type to format date as YYYY-MM-DD
type DateOnly struct {
	time.Time
}

const dateFormat = "2006-01-02"

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		// Empty string means zero time
		d.Time = time.Time{}
		return nil
	}
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Time.Format(dateFormat))
}

// You can also add a String() method if needed
func (d DateOnly) String() string {
	return d.Time.Format(dateFormat)
}

func ToDateOnly(t time.Time) DateOnly {
	return DateOnly{Time: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())}
}

// ToTime converts a DateOnly value back to a time.Time (with a time of 00:00:00).
func ToTime(d DateOnly) time.Time {
	return d.Time // DateOnly is a wrapper around time.Time, so this is simple
}
