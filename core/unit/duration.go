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
