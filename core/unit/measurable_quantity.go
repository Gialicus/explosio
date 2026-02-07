package unit

import "fmt"

type MeasurableQuantity struct {
	Value float64
	Unit  MeasurableUnit
}

func DefaultMeasurableQuantity() *MeasurableQuantity {
	return &MeasurableQuantity{Value: 0, Unit: UnitMeter}
}

func NewMeasurableQuantity(value float64, unit MeasurableUnit) *MeasurableQuantity {
	return &MeasurableQuantity{Value: value, Unit: unit}
}

func (m *MeasurableQuantity) String() string {
	return fmt.Sprintf("%.0f %s", m.Value, m.Unit)
}
