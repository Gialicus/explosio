package unit

import "fmt"

// MeasurableQuantity is a quantity with unit (e.g. 5 kg, 10 m).
type MeasurableQuantity struct {
	Value float64
	Unit  MeasurableUnit
}

// NewMeasurableQuantity creates a quantity with value and unit.
func NewMeasurableQuantity(value float64, unit MeasurableUnit) *MeasurableQuantity {
	return &MeasurableQuantity{Value: value, Unit: unit}
}

// String formats the quantity for output (e.g. "5 kg").
func (m *MeasurableQuantity) String() string {
	return fmt.Sprintf("%.0f %s", m.Value, m.Unit)
}

// SetValue sets the value and returns the pointer for chaining.
func (m *MeasurableQuantity) SetValue(value float64) *MeasurableQuantity {
	m.Value = value
	return m
}

// SetUnit sets the unit and returns the pointer for chaining.
func (m *MeasurableQuantity) SetUnit(unit MeasurableUnit) *MeasurableQuantity {
	m.Unit = unit
	return m
}
