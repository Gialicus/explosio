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

// CalculatePrice returns the total price of the material (unit price multiplied by quantity).
func (c *CountableMaterial) CalculatePrice() float64 {
	return c.Price.Value * float64(c.Quantity)
}

// SetTotalPrice sets the total price and derives the unit price from quantity.
// If Quantity is 0, unit price is set to 0 (avoids division by zero).
func (c *CountableMaterial) SetTotalPrice(totalPrice unit.Price) {
	if c.Quantity == 0 {
		c.Price = unit.Price{Value: 0, Currency: totalPrice.Currency}
		return
	}
	c.Price = unit.Price{
		Value:    totalPrice.Value / float64(c.Quantity),
		Currency: totalPrice.Currency,
	}
}

// Clone returns a deep copy of the countable material.
func (c *CountableMaterial) Clone() *CountableMaterial {
	return NewCountableMaterial(c.Name, c.Description, c.Price, c.Quantity)
}
