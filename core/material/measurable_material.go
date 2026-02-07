package material

import "explosio/core/unit"

// MeasurableMaterial is a material that is measurable, for example: a pipe 1 meter long
type MeasurableMaterial struct {
	Name        string
	Description string
	Price       unit.Price
	Quantity    unit.MeasurableQuantity
}
