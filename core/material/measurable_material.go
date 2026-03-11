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

// CalculatePrice returns the total price of the material (unit price multiplied by quantity).
func (m *MeasurableMaterial) CalculatePrice() float64 {
	return m.Price.Value * m.Quantity.Value
}

// SetTotalPrice sets the total price and derives the unit price from quantity.
// If Quantity.Value is 0, unit price is set to 0 (avoids division by zero).
func (m *MeasurableMaterial) SetTotalPrice(totalPrice unit.Price) {
	if m.Quantity.Value == 0 {
		m.Price = unit.Price{Value: 0, Currency: totalPrice.Currency}
		return
	}
	m.Price = unit.Price{
		Value:    totalPrice.Value / m.Quantity.Value,
		Currency: totalPrice.Currency,
	}
}

// Clone returns a deep copy of the measurable material.
func (m *MeasurableMaterial) Clone() *MeasurableMaterial {
	return NewMeasurableMaterial(m.Name, m.Description, m.Price, m.Quantity)
}
