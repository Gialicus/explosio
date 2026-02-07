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

func DefaultDuration() *Duration {
	return &Duration{Value: 0, Unit: DurationUnitHour}
}

func NewDuration(value float64, unit DurationUnit) *Duration {
	return &Duration{Value: value, Unit: unit}
}

func (d *Duration) String() string {
	return fmt.Sprintf("%.0f %s", d.Value, d.Unit)
}
