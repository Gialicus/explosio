package material

import "explosio/core/unit"

// CountableMaterialBuilder builds a countable material.
type CountableMaterialBuilder struct {
	countableMaterial *CountableMaterial
}

// NewCountableMaterialBuilder creates a new countable material builder with empty name and description, zero price and zero quantity.
func NewCountableMaterialBuilder() *CountableMaterialBuilder {
	return &CountableMaterialBuilder{
		countableMaterial: &CountableMaterial{
			Name:        "",
			Description: "",
			Price:       unit.Price{Value: 0, Currency: "EUR"},
			Quantity:    0,
		},
	}
}

// WithName sets the name and returns the builder for chaining.
func (b *CountableMaterialBuilder) WithName(name string) *CountableMaterialBuilder {
	b.countableMaterial.Name = name
	return b
}

// WithDescription sets the description and returns the builder for chaining.
func (b *CountableMaterialBuilder) WithDescription(description string) *CountableMaterialBuilder {
	b.countableMaterial.Description = description
	return b
}

// WithPrice sets the price and returns the builder for chaining.
func (b *CountableMaterialBuilder) WithPrice(price unit.Price) *CountableMaterialBuilder {
	b.countableMaterial.Price = price
	return b
}

// WithQuantity sets the quantity and returns the builder for chaining.
func (b *CountableMaterialBuilder) WithQuantity(quantity int) *CountableMaterialBuilder {
	b.countableMaterial.Quantity = quantity
	return b
}

// WithTotalPrice sets the total price and derives the unit price from quantity.
// Quantity must be set before calling this. If Quantity is 0, unit price becomes 0.
func (b *CountableMaterialBuilder) WithTotalPrice(totalPrice unit.Price) *CountableMaterialBuilder {
	b.countableMaterial.SetTotalPrice(totalPrice)
	return b
}

// Build returns the built countable material.
func (b *CountableMaterialBuilder) Build() *CountableMaterial {
	return b.countableMaterial
}
