// Package human defines human resource types for project activities.
package human

import (
	"explosio/core/resource"
	"explosio/core/unit"
)

// HumanResource represents a person or role with a duration and price.
// It embeds resource.PricedResource for shared logic.
type HumanResource struct {
	resource.PricedResource
}

// NewHumanResource creates a human resource with name, description, duration, and price.
func NewHumanResource(name string, description string, duration unit.Duration, price unit.Price) *HumanResource {
	return &HumanResource{
		PricedResource: resource.PricedResource{
			Name:        name,
			Description: description,
			Duration:    duration,
			Price:       price,
		},
	}
}
