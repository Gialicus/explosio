package human

import (
	"errors"
	"explosio/core/unit"
)

// HumanResourceBuilder builds a human resource.
type HumanResourceBuilder struct {
	humanResource *HumanResource
}

// NewHumanResourceBuilder creates a new human resource builder.
func NewHumanResourceBuilder() *HumanResourceBuilder {
	return &HumanResourceBuilder{humanResource: NewHumanResource("", "", *unit.NewDuration(0, unit.DurationUnitHour), *unit.NewPrice(0, "EUR"))}
}

// WithName sets the name and returns the builder for chaining.
func (b *HumanResourceBuilder) WithName(name string) *HumanResourceBuilder {
	b.humanResource.Name = name
	return b
}

// WithDescription sets the description and returns the builder for chaining.
func (b *HumanResourceBuilder) WithDescription(description string) *HumanResourceBuilder {
	b.humanResource.Description = description
	return b
}

// WithDuration sets the duration and returns the builder for chaining.
func (b *HumanResourceBuilder) WithDuration(duration unit.Duration) *HumanResourceBuilder {
	b.humanResource.Duration = duration
	return b
}

// WithPrice sets the price and returns the builder for chaining.
func (b *HumanResourceBuilder) WithPrice(price unit.Price) *HumanResourceBuilder {
	b.humanResource.Price = price
	return b
}

// WithTotalPrice sets the total price. Hourly and daily rates are derived from Price/Duration.
func (b *HumanResourceBuilder) WithTotalPrice(totalPrice unit.Price) *HumanResourceBuilder {
	b.humanResource.SetTotalPrice(totalPrice)
	return b
}

// WithHourlyRate sets the hourly rate and derives the total price from duration.
// Duration must be set before calling this. If Duration.ToHours() is 0, total price becomes 0.
func (b *HumanResourceBuilder) WithHourlyRate(rate unit.Price) *HumanResourceBuilder {
	b.humanResource.SetHourlyRate(rate)
	return b
}

// WithDailyRate sets the daily rate (per working day) and derives the total price from duration.
// Duration must be set before calling this. If Duration.ToHours() is 0, total price becomes 0.
func (b *HumanResourceBuilder) WithDailyRate(rate unit.Price) *HumanResourceBuilder {
	b.humanResource.SetDailyRate(rate)
	return b
}

// Build builds the human resource. Returns an error if name is empty, duration is negative, or price is invalid (negative value or empty currency).
func (b *HumanResourceBuilder) Build() (*HumanResource, error) {
	if b.humanResource.Name == "" {
		return nil, errors.New("human resource name cannot be empty")
	}
	if b.humanResource.Duration.Value < 0 {
		return nil, errors.New("human resource duration cannot be negative")
	}
	if b.humanResource.Price.Value < 0 {
		return nil, errors.New("human resource price cannot be negative")
	}
	if b.humanResource.Price.Currency == "" {
		return nil, errors.New("human resource price currency cannot be empty")
	}
	return b.humanResource, nil
}
