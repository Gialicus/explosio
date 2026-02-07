package unit

import "fmt"

type DurationUnit string

const (
	DurationUnitMinute DurationUnit = "minute"
	DurationUnitHour   DurationUnit = "hour"
	DurationUnitDay    DurationUnit = "day"
	DurationUnitWeek   DurationUnit = "week"
	DurationUnitMonth  DurationUnit = "month"
	DurationUnitYear   DurationUnit = "year"
)

type Duration struct {
	Value float64
	Unit  DurationUnit
}

func NewDuration(value float64, unit DurationUnit) *Duration {
	return &Duration{Value: value, Unit: unit}
}

func (d *Duration) String() string {
	return fmt.Sprintf("%.0f %s", d.Value, d.Unit)
}

func (d *Duration) SetValue(value float64) *Duration {
	d.Value = value
	return d
}

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
