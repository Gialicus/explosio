package asset

import "explosio/core/unit"

// Asset represents everything that can has a price and a duration.
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

// CalculatePrice returns the total price (asset plus all materials and sub-activities).
func (a *Asset) CalculatePrice() float64 {
	return a.Price.Value
}

// CalculateDuration returns the total duration (asset plus sub-activities) in the root's value/unit.
func (a *Asset) CalculateDuration() float64 {
	return a.Duration.Value
}
