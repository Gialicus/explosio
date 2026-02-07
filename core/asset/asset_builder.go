package asset

import "explosio/core/unit"

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

// Build builds the asset.
func (b *AssetBuilder) Build() *Asset {
	return b.asset
}
