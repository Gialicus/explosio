package material

import "explosio/core/unit"

// CountableMaterial is a material that is countable, for example: a set of 5 screws
type CountableMaterial struct {
	Name        string
	Description string
	Price       unit.Price
	Quantity    int
}
