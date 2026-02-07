package material

import "explosio/core"

type MeasurableMaterial struct {
	Name        string
	Description string
	Price       core.Price
	Quantity    core.MeasurableQuantity
}