package asset

import "explosio/core/unit"

// Asset represents an asset with name, description, price, quantity.
type Asset struct {
	Name        string
	Description string
	Price       unit.Price
	Duration    unit.Duration
}

// NewAsset creates an asset with name and description, zero price and quantity.
func NewAsset(name string, description string, price unit.Price, duration unit.Duration) *Asset {
	return &Asset{Name: name, Description: description, Price: price, Duration: duration}
}
