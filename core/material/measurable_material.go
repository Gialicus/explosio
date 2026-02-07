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

// SetName sets the name and returns the receiver for chaining.
func (m *MeasurableMaterial) SetName(name string) *MeasurableMaterial {
	m.Name = name
	return m
}

// SetDescription sets the description and returns the receiver for chaining.
func (m *MeasurableMaterial) SetDescription(description string) *MeasurableMaterial {
	m.Description = description
	return m
}

// SetPrice sets the price and returns the receiver for chaining.
func (m *MeasurableMaterial) SetPrice(price unit.Price) *MeasurableMaterial {
	m.Price = price
	return m
}

// SetQuantity sets the quantity and returns the receiver for chaining.
func (m *MeasurableMaterial) SetQuantity(quantity unit.MeasurableQuantity) *MeasurableMaterial {
	m.Quantity = quantity
	return m
}

// CalculatePrice returns the material price.
func (m *MeasurableMaterial) CalculatePrice() float64 {
	return m.Price.Value
}
