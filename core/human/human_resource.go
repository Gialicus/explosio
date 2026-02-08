package human

import "explosio/core/unit"

type HumanResource struct {
	Name        string
	Description string
	Duration    unit.Duration
	Price       unit.Price
}

func NewHumanResource(name string, description string, duration unit.Duration, price unit.Price) *HumanResource {
	return &HumanResource{Name: name, Description: description, Duration: duration, Price: price}
}

// CalculatePrice returns the total price (human resource plus all materials and sub-activities).
func (h *HumanResource) CalculatePrice() float64 {
	return h.Price.Value
}

// CalculateDuration returns the total duration (human resource plus sub-activities) in the root's value/unit.
func (h *HumanResource) CalculateDuration() float64 {
	return h.Duration.Value
}
