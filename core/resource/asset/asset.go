// Package asset defines asset types for project activities (equipment, tools, etc.).
package asset

import (
	"explosio/core/resource"
	"explosio/core/unit"
)

// Asset represents everything that can have a price and a duration.
// It embeds resource.PricedResource for shared logic.
type Asset struct {
	resource.PricedResource
}

// NewAsset creates an asset with name and description, zero price and duration.
func NewAsset(name string, description string, price unit.Price, duration unit.Duration) *Asset {
	return &Asset{
		PricedResource: resource.PricedResource{
			Name:        name,
			Description: description,
			Price:       price,
			Duration:    duration,
		},
	}
}

// Clone returns a deep copy of the asset.
func (a *Asset) Clone() *Asset {
	return NewAsset(a.Name, a.Description, a.Price, a.Duration)
}
