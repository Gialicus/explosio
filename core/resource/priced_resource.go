// Package resource defines shared types for priced resources (assets, human resources).
package resource

import "explosio/core/unit"

// PricedResource holds common fields and logic for types with name, description, price, and duration.
// Asset and HumanResource embed this to avoid code duplication.
type PricedResource struct {
	Name        string
	Description string
	Price       unit.Price
	Duration    unit.Duration
}

// CalculatePrice returns the price value.
func (p *PricedResource) CalculatePrice() float64 {
	return p.Price.Value
}

// CalculateDuration returns the duration value.
func (p *PricedResource) CalculateDuration() float64 {
	return p.Duration.Value
}

// CalculateHourlyRate returns the hourly rate.
// Returns 0 if duration is zero to avoid division by zero.
func (p *PricedResource) CalculateHourlyRate() float64 {
	hours := p.Duration.ToHours()
	if hours == 0 {
		return 0
	}
	return p.Price.Value / hours
}

// CalculateDailyRate returns the daily rate (per working day).
// Returns 0 if duration is zero to avoid division by zero.
func (p *PricedResource) CalculateDailyRate() float64 {
	hours := p.Duration.ToHours()
	if hours == 0 {
		return 0
	}
	return p.Price.Value / hours * unit.WorkingHoursPerDay
}

// SetTotalPrice sets the total price. Hourly and daily rates are derived from Price/Duration.
func (p *PricedResource) SetTotalPrice(totalPrice unit.Price) {
	p.Price = totalPrice
}

// SetHourlyRate sets the hourly rate and derives the total price from duration.
// If Duration.ToHours() is 0, total price is set to 0.
func (p *PricedResource) SetHourlyRate(rate unit.Price) {
	hours := p.Duration.ToHours()
	if hours == 0 {
		p.Price = unit.Price{Value: 0, Currency: rate.Currency}
		return
	}
	p.Price = unit.Price{Value: rate.Value * hours, Currency: rate.Currency}
}

// SetDailyRate sets the daily rate (per working day of WorkingHoursPerDay hours) and derives the total price from duration.
// Duration is converted to calendar hours via ToHours() (e.g. 1 day = 24h). Total = rate * (hours / WorkingHoursPerDay).
// Example: 8 hours duration, 80 EUR/day → total = 80 * 8/8 = 80 EUR. 24 hours, 80 EUR/day → total = 80 * 24/8 = 240 EUR.
// If Duration.ToHours() is 0, total price is set to 0.
func (p *PricedResource) SetDailyRate(rate unit.Price) {
	hours := p.Duration.ToHours()
	if hours == 0 {
		p.Price = unit.Price{Value: 0, Currency: rate.Currency}
		return
	}
	p.Price = unit.Price{
		Value:    rate.Value * hours / unit.WorkingHoursPerDay,
		Currency: rate.Currency,
	}
}
