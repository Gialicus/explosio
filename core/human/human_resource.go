// Package human defines human resource types for project activities.
package human

import "explosio/core/unit"

// HumanResource represents a person or role with a duration and price.
type HumanResource struct {
	Name        string
	Description string
	Duration    unit.Duration
	Price       unit.Price
}

// NewHumanResource creates a human resource with name, description, duration, and price.
func NewHumanResource(name string, description string, duration unit.Duration, price unit.Price) *HumanResource {
	return &HumanResource{Name: name, Description: description, Duration: duration, Price: price}
}

// CalculatePrice returns the human resource price.
func (h *HumanResource) CalculatePrice() float64 {
	return h.Price.Value
}

// CalculateDuration returns the human resource duration value.
func (h *HumanResource) CalculateDuration() float64 {
	return h.Duration.Value
}

// CalculateHourlyRate returns the hourly rate of the human resource.
// Returns 0 if duration is zero to avoid division by zero.
func (h *HumanResource) CalculateHourlyRate() float64 {
	hours := h.Duration.ToHours()
	if hours == 0 {
		return 0
	}
	return h.Price.Value / hours
}

// CalculateDailyRate returns the daily rate of the human resource.
// Returns 0 if duration is zero to avoid division by zero.
func (h *HumanResource) CalculateDailyRate() float64 {
	hours := h.Duration.ToHours()
	if hours == 0 {
		return 0
	}
	return h.Price.Value / hours * unit.WorkingHoursPerDay
}
