package material

import "explosio/core"

type ComplexMaterial struct {
	Name        string
	Description string
	Price       core.Price
	Quantity    core.MeasurableQuantity
}