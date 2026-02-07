package material

import "explosio/core/unit"

// CountableMaterial is a countable material (e.g. screws, pieces).
type CountableMaterial struct {
	Name        string
	Description string
	Price       unit.Price
	Quantity    int
}

// NewCountableMaterial creates a countable material (e.g. screws, pieces).
func NewCountableMaterial(name string, description string, price unit.Price, quantity int) *CountableMaterial {
	return &CountableMaterial{Name: name, Description: description, Price: price, Quantity: quantity}
}

// CalculatePrice returns the unit price of the material (not multiplied by Quantity).
func (c *CountableMaterial) CalculatePrice() float64 {
	return c.Price.Value
}
