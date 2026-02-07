package core

import "explosio/core/unit"

// ActivityBuilder builds an activity.
type ActivityBuilder struct {
	activity *Activity
}

// NewActivityBuilder creates a new activity builder with empty name and description, zero duration and zero EUR price.
func NewActivityBuilder() *ActivityBuilder {
	return &ActivityBuilder{
		activity: NewActivity("", ""),
	}
}

// WithName sets the name and returns the builder for chaining.
func (b *ActivityBuilder) WithName(name string) *ActivityBuilder {
	b.activity.Name = name
	return b
}

// WithDescription sets the description and returns the builder for chaining.
func (b *ActivityBuilder) WithDescription(description string) *ActivityBuilder {
	b.activity.Description = description
	return b
}

// WithDuration sets the duration and returns the builder for chaining.
func (b *ActivityBuilder) WithDuration(duration unit.Duration) *ActivityBuilder {
	b.activity.Duration = duration
	return b
}

// WithPrice sets the price and returns the builder for chaining.
func (b *ActivityBuilder) WithPrice(price unit.Price) *ActivityBuilder {
	b.activity.Price = price
	return b
}

// Build returns the built activity.
func (b *ActivityBuilder) Build() *Activity {
	return b.activity
}
