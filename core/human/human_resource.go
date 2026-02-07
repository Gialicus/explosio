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

// SetName sets the name and returns the receiver for chaining.
func (h *HumanResource) SetName(name string) *HumanResource {
	h.Name = name
	return h
}

// SetDescription sets the description and returns the receiver for chaining.
func (h *HumanResource) SetDescription(description string) *HumanResource {
	h.Description = description
	return h
}

// SetDuration sets the duration and returns the receiver for chaining.
func (h *HumanResource) SetDuration(duration unit.Duration) *HumanResource {
	h.Duration = duration
	return h
}

// SetPrice sets the price and returns the receiver for chaining.
func (h *HumanResource) SetPrice(price unit.Price) *HumanResource {
	h.Price = price
	return h
}
