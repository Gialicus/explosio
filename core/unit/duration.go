// Package unit defines types for durations, prices, and measurable quantities.
package unit

import "fmt"

// DurationUnit is the unit of duration (minute, hour, day, ...).
type DurationUnit string

const (
	DurationUnitMinute DurationUnit = "minute"
	DurationUnitHour   DurationUnit = "hour"
	DurationUnitDay    DurationUnit = "day"
	DurationUnitWeek   DurationUnit = "week"
	DurationUnitMonth  DurationUnit = "month"
	DurationUnitYear   DurationUnit = "year"
)

const (
	WORKING_HOURS_PER_DAY float64 = 8.0
)

// Duration represents a time interval (value plus unit).
type Duration struct {
	Value float64
	Unit  DurationUnit
}

// NewDuration creates a duration with value and unit.
func NewDuration(value float64, unit DurationUnit) *Duration {
	return &Duration{Value: value, Unit: unit}
}

// String formats the duration for output (e.g. "2 hour").
func (d *Duration) String() string {
	return fmt.Sprintf("%.0f %s", d.Value, d.Unit)
}

// SetValue sets the value and returns the pointer for chaining.
func (d *Duration) SetValue(value float64) *Duration {
	d.Value = value
	return d
}

// SetUnit sets the unit and returns the pointer for chaining.
func (d *Duration) SetUnit(unit DurationUnit) *Duration {
	d.Unit = unit
	return d
}

// ToHours returns the duration in hours for comparison and aggregation.
func (d *Duration) ToHours() float64 {
	if d == nil {
		return 0
	}
	switch d.Unit {
	case DurationUnitMinute:
		return d.Value / 60
	case DurationUnitHour:
		return d.Value
	case DurationUnitDay:
		return d.Value * 24
	case DurationUnitWeek:
		return d.Value * 24 * 7
	case DurationUnitMonth:
		return d.Value * 24 * 30
	case DurationUnitYear:
		return d.Value * 24 * 365
	default:
		return d.Value
	}
}
