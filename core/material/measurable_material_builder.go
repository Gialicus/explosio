package material

import "explosio/core/unit"

// MeasurableMaterialBuilder builds a measurable material.
type MeasurableMaterialBuilder struct {
	measurableMaterial *MeasurableMaterial
}

// NewMeasurableMaterialBuilder creates a new measurable material builder with empty name and description, zero price and zero quantity (meter).
func NewMeasurableMaterialBuilder() *MeasurableMaterialBuilder {
	return &MeasurableMaterialBuilder{
		measurableMaterial: &MeasurableMaterial{
			Name:        "",
			Description: "",
			Price:       unit.Price{Value: 0, Currency: "EUR"},
			Quantity:    unit.MeasurableQuantity{Value: 0, Unit: unit.UnitMeter},
		},
	}
}

// WithName sets the name and returns the builder for chaining.
func (b *MeasurableMaterialBuilder) WithName(name string) *MeasurableMaterialBuilder {
	b.measurableMaterial.Name = name
	return b
}

// WithDescription sets the description and returns the builder for chaining.
func (b *MeasurableMaterialBuilder) WithDescription(description string) *MeasurableMaterialBuilder {
	b.measurableMaterial.Description = description
	return b
}

// WithPrice sets the price and returns the builder for chaining.
func (b *MeasurableMaterialBuilder) WithPrice(price unit.Price) *MeasurableMaterialBuilder {
	b.measurableMaterial.Price = price
	return b
}

// WithQuantity sets the quantity and returns the builder for chaining.
func (b *MeasurableMaterialBuilder) WithQuantity(quantity unit.MeasurableQuantity) *MeasurableMaterialBuilder {
	b.measurableMaterial.Quantity = quantity
	return b
}

// Build returns the built measurable material.
func (b *MeasurableMaterialBuilder) Build() *MeasurableMaterial {
	return b.measurableMaterial
}
