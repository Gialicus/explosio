package core

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
