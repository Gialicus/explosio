package material

import "explosio/core/unit"

// MeasurableMaterial is a material with measurable quantity (e.g. 5 kg cement, 10 m cable).
type MeasurableMaterial struct {
	Name        string
	Description string
	Price       unit.Price
	Quantity    unit.MeasurableQuantity
}

// NewMeasurableMaterial creates a material with measurable quantity (e.g. kg, m).
func NewMeasurableMaterial(name string, description string, price unit.Price, quantity unit.MeasurableQuantity) *MeasurableMaterial {
	return &MeasurableMaterial{Name: name, Description: description, Price: price, Quantity: quantity}
}

// CalculatePrice returns the material price.
func (m *MeasurableMaterial) CalculatePrice() float64 {
	return m.Price.Value
}
