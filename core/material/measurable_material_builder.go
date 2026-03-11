package material

import (
	"errors"
	"explosio/core/unit"
)

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

// WithTotalPrice sets the total price and derives the unit price from quantity.
// Quantity must be set before calling this. If Quantity.Value is 0, unit price becomes 0.
func (b *MeasurableMaterialBuilder) WithTotalPrice(totalPrice unit.Price) *MeasurableMaterialBuilder {
	b.measurableMaterial.SetTotalPrice(totalPrice)
	return b
}

// Build returns the built measurable material. Returns an error if name is empty, price is invalid (negative value or empty currency), or quantity is negative.
func (b *MeasurableMaterialBuilder) Build() (*MeasurableMaterial, error) {
	if b.measurableMaterial.Name == "" {
		return nil, errors.New("measurable material name cannot be empty")
	}
	if b.measurableMaterial.Price.Value < 0 {
		return nil, errors.New("measurable material price cannot be negative")
	}
	if b.measurableMaterial.Price.Currency == "" {
		return nil, errors.New("measurable material price currency cannot be empty")
	}
	if b.measurableMaterial.Quantity.Value < 0 {
		return nil, errors.New("measurable material quantity cannot be negative")
	}
	return b.measurableMaterial, nil
}
