package unit

import "fmt"

type MeasurableQuantity struct {
	Value float64
	Unit  MeasurableUnit
}

func NewMeasurableQuantity(value float64, unit MeasurableUnit) *MeasurableQuantity {
	return &MeasurableQuantity{Value: value, Unit: unit}
}

func (m *MeasurableQuantity) String() string {
	return fmt.Sprintf("%.0f %s", m.Value, m.Unit)
}

func (m *MeasurableQuantity) SetValue(value float64) *MeasurableQuantity {
	m.Value = value
	return m
}

func (m *MeasurableQuantity) SetUnit(unit MeasurableUnit) *MeasurableQuantity {
	m.Unit = unit
	return m
}
