package asset

import "explosio/core/unit"

// Asset represents everything that can have a price and a duration.
type Asset struct {
	Name        string
	Description string
	Price       unit.Price
	Duration    unit.Duration
}

// NewAsset creates an asset with name and description, zero price and duration.
func NewAsset(name string, description string, price unit.Price, duration unit.Duration) *Asset {
	return &Asset{Name: name, Description: description, Price: price, Duration: duration}
}

// CalculatePrice returns the asset price.
func (a *Asset) CalculatePrice() float64 {
	return a.Price.Value
}

// CalculateDuration returns the asset duration value.
func (a *Asset) CalculateDuration() float64 {
	return a.Duration.Value
}

// CalculateHourlyRate returns the hourly rate of the asset.
// Returns 0 if duration is zero to avoid division by zero.
func (a *Asset) CalculateHourlyRate() float64 {
	hours := a.Duration.ToHours()
	if hours == 0 {
		return 0
	}
	return a.Price.Value / hours
}

// CalculateDailyRate returns the daily rate of the asset.
// Returns 0 if duration is zero to avoid division by zero.
func (a *Asset) CalculateDailyRate() float64 {
	hours := a.Duration.ToHours()
	if hours == 0 {
		return 0
	}
	return a.Price.Value / hours * unit.WorkingHoursPerDay
}

// SetTotalPrice sets the total price. Hourly and daily rates are derived from Price/Duration.
func (a *Asset) SetTotalPrice(totalPrice unit.Price) {
	a.Price = totalPrice
}

// SetHourlyRate sets the hourly rate and derives the total price from duration.
// If Duration.ToHours() is 0, total price is set to 0.
func (a *Asset) SetHourlyRate(rate unit.Price) {
	hours := a.Duration.ToHours()
	if hours == 0 {
		a.Price = unit.Price{Value: 0, Currency: rate.Currency}
		return
	}
	a.Price = unit.Price{Value: rate.Value * hours, Currency: rate.Currency}
}

// SetDailyRate sets the daily rate (per working day) and derives the total price from duration.
// If Duration.ToHours() is 0, total price is set to 0.
func (a *Asset) SetDailyRate(rate unit.Price) {
	hours := a.Duration.ToHours()
	if hours == 0 {
		a.Price = unit.Price{Value: 0, Currency: rate.Currency}
		return
	}
	a.Price = unit.Price{
		Value:    rate.Value * hours / unit.WorkingHoursPerDay,
		Currency: rate.Currency,
	}
}
