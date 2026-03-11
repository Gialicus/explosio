// Package unit defines date type for activity scheduling.
package unit

import "time"

// Date wraps time.Time for activity start/end dates.
type Date struct {
	Time time.Time
}

// NewDate creates a date from year, month, day.
func NewDate(year int, month time.Month, day int) Date {
	return Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// AddHours adds hours to the date.
func (d Date) AddHours(hours float64) Date {
	return Date{Time: d.Time.Add(time.Duration(hours * float64(time.Hour)))}
}

// String returns the date in YYYY-MM-DD format.
func (d Date) String() string {
	return d.Time.Format("2006-01-02")
}
