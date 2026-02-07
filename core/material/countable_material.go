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

// SetName sets the name and returns the receiver for chaining.
func (c *CountableMaterial) SetName(name string) *CountableMaterial {
	c.Name = name
	return c
}

// SetDescription sets the description and returns the receiver for chaining.
func (c *CountableMaterial) SetDescription(description string) *CountableMaterial {
	c.Description = description
	return c
}

// SetPrice sets the price and returns the receiver for chaining.
func (c *CountableMaterial) SetPrice(price unit.Price) *CountableMaterial {
	c.Price = price
	return c
}

// SetQuantity sets the quantity and returns the receiver for chaining.
func (c *CountableMaterial) SetQuantity(quantity int) *CountableMaterial {
	c.Quantity = quantity
	return c
}

// CalculatePrice returns the unit price of the material (not multiplied by Quantity).
func (c *CountableMaterial) CalculatePrice() float64 {
	return c.Price.Value
}
