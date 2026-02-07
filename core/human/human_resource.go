package human

import "explosio/core/unit"

// HumanResource represents a human resource with name, description, duration, price.
type HumanResource struct {
	Name        string
	Description string
	Duration    unit.Duration
	Price       unit.Price
}

// NewHumanResource creates a human resource with name and description, zero duration and zero EUR price.
func NewHumanResource(name string, description string, duration unit.Duration, price unit.Price) *HumanResource {
	return &HumanResource{Name: name, Description: description, Duration: duration, Price: price}
}
