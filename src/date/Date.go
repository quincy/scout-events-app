package date

import (
	"encoding/json"
	"fmt"
	"time"
)

// Date provides helpful functions for representing a date without respect to a specific time.
type Date struct {
	time.Time
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func ParseDate(value string) (Date, error) {
	t, err := time.Parse("01/02/2006", value)
	if err != nil {
		return Date{}, err
	}
	return Date{t}, nil
}

func (d Date) String() string {
	return d.Format("01/02/2006")
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

func UnmarshalJSON(data []byte) (Date, error) {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return Date{}, err
	}

	parsed, err := ParseDate(dateStr)
	if err != nil {
		return Date{}, err
	}
	return parsed, nil
}
