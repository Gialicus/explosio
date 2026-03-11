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

// SetTotalPrice sets the total price. Hourly and daily rates are derived from Price/Duration.
func (h *HumanResource) SetTotalPrice(totalPrice unit.Price) {
	h.Price = totalPrice
}

// SetHourlyRate sets the hourly rate and derives the total price from duration.
// If Duration.ToHours() is 0, total price is set to 0.
func (h *HumanResource) SetHourlyRate(rate unit.Price) {
	hours := h.Duration.ToHours()
	if hours == 0 {
		h.Price = unit.Price{Value: 0, Currency: rate.Currency}
		return
	}
	h.Price = unit.Price{Value: rate.Value * hours, Currency: rate.Currency}
}

// SetDailyRate sets the daily rate (per working day) and derives the total price from duration.
// If Duration.ToHours() is 0, total price is set to 0.
func (h *HumanResource) SetDailyRate(rate unit.Price) {
	hours := h.Duration.ToHours()
	if hours == 0 {
		h.Price = unit.Price{Value: 0, Currency: rate.Currency}
		return
	}
	h.Price = unit.Price{
		Value:    rate.Value * hours / unit.WorkingHoursPerDay,
		Currency: rate.Currency,
	}
}
