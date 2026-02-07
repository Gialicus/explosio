package material

import "explosio/core"

type CountableMaterial struct {
	Name        string
	Description string
	Price       core.Price
	Quantity    int
}