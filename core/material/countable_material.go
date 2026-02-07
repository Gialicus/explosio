package material

import "explosio/core/unit"

// CountableMaterial is a material that is countable, for example: a set of 5 screws
type CountableMaterial struct {
	Name        string
	Description string
	Price       unit.Price
	Quantity    int
}

func NewCountableMaterial(name string, description string, price unit.Price, quantity int) *CountableMaterial {
	return &CountableMaterial{Name: name, Description: description, Price: price, Quantity: quantity}
}

func (c *CountableMaterial) CalculatePrice() float64 {
	return c.Price.Value
}
