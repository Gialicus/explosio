package asset

import (
	"errors"
	"explosio/core/unit"
)

// AssetBuilder builds an asset.
type AssetBuilder struct {
	asset *Asset
}

// NewAssetBuilder creates a new asset builder with empty name and description, zero price and duration.
func NewAssetBuilder() *AssetBuilder {
	return &AssetBuilder{asset: NewAsset("", "", *unit.NewPrice(0, "EUR"), *unit.NewDuration(0, unit.DurationUnitHour))}
}

// WithName sets the name and returns the builder for chaining.
func (b *AssetBuilder) WithName(name string) *AssetBuilder {
	b.asset.Name = name
	return b
}

// WithDescription sets the description and returns the builder for chaining.
func (b *AssetBuilder) WithDescription(description string) *AssetBuilder {
	b.asset.Description = description
	return b
}

// WithPrice sets the price and returns the builder for chaining.
func (b *AssetBuilder) WithPrice(price unit.Price) *AssetBuilder {
	b.asset.Price = price
	return b
}

// WithDuration sets the duration and returns the builder for chaining.
func (b *AssetBuilder) WithDuration(duration unit.Duration) *AssetBuilder {
	b.asset.Duration = duration
	return b
}

// WithTotalPrice sets the total price. Hourly and daily rates are derived from Price/Duration.
func (b *AssetBuilder) WithTotalPrice(totalPrice unit.Price) *AssetBuilder {
	b.asset.SetTotalPrice(totalPrice)
	return b
}

// WithHourlyRate sets the hourly rate and derives the total price from duration.
// Duration must be set before calling this. If Duration.ToHours() is 0, total price becomes 0.
func (b *AssetBuilder) WithHourlyRate(rate unit.Price) *AssetBuilder {
	b.asset.SetHourlyRate(rate)
	return b
}

// WithDailyRate sets the daily rate (per working day) and derives the total price from duration.
// Duration must be set before calling this. If Duration.ToHours() is 0, total price becomes 0.
func (b *AssetBuilder) WithDailyRate(rate unit.Price) *AssetBuilder {
	b.asset.SetDailyRate(rate)
	return b
}

// Build builds the asset. Returns an error if name is empty, duration is negative, or price is invalid (negative value or empty currency).
func (b *AssetBuilder) Build() (*Asset, error) {
	if b.asset.Name == "" {
		return nil, errors.New("asset name cannot be empty")
	}
	if b.asset.Duration.Value < 0 {
		return nil, errors.New("asset duration cannot be negative")
	}
	if b.asset.Price.Value < 0 {
		return nil, errors.New("asset price cannot be negative")
	}
	if b.asset.Price.Currency == "" {
		return nil, errors.New("asset price currency cannot be empty")
	}
	return b.asset, nil
}
