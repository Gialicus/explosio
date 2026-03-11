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
