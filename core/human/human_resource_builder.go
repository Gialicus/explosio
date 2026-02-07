package human

import "explosio/core/unit"

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

// Build builds the human resource.
func (b *HumanResourceBuilder) Build() *HumanResource {
	return b.humanResource
}
