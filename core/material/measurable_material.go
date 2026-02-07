package material

import "explosio/core/unit"

// MeasurableMaterial is a material that is measurable, for example: 5kg of cement
type MeasurableMaterial struct {
	Name        string
	Description string
	Price       unit.Price
	Quantity    unit.MeasurableQuantity
}

func NewMeasurableMaterial(name string, description string, price unit.Price, quantity unit.MeasurableQuantity) *MeasurableMaterial {
	return &MeasurableMaterial{Name: name, Description: description, Price: price, Quantity: quantity}
}

func (m *MeasurableMaterial) CalculatePrice() float64 {
	return m.Price.Value
}
