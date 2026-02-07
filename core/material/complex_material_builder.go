package material

import "explosio/core/unit"

// ComplexMaterialBuilder builds a complex material.
type ComplexMaterialBuilder struct {
	complexMaterial *ComplexMaterial
}

// NewComplexMaterialBuilder creates a new complex material builder with empty name and description, zero price, zero unit quantity and nil measurable material.
func NewComplexMaterialBuilder() *ComplexMaterialBuilder {
	return &ComplexMaterialBuilder{
		complexMaterial: &ComplexMaterial{
			Name:               "",
			Description:        "",
			Price:              unit.Price{Value: 0, Currency: "EUR"},
			UnitQuantity:       0,
			MeasurableMaterial: nil,
		},
	}
}

// WithName sets the name and returns the builder for chaining.
func (b *ComplexMaterialBuilder) WithName(name string) *ComplexMaterialBuilder {
	b.complexMaterial.Name = name
	return b
}

// WithDescription sets the description and returns the builder for chaining.
func (b *ComplexMaterialBuilder) WithDescription(description string) *ComplexMaterialBuilder {
	b.complexMaterial.Description = description
	return b
}

// WithPrice sets the price and returns the builder for chaining.
func (b *ComplexMaterialBuilder) WithPrice(price unit.Price) *ComplexMaterialBuilder {
	b.complexMaterial.Price = price
	return b
}

// WithUnitQuantity sets the unit quantity and returns the builder for chaining.
func (b *ComplexMaterialBuilder) WithUnitQuantity(unitQuantity int) *ComplexMaterialBuilder {
	b.complexMaterial.UnitQuantity = unitQuantity
	return b
}

// WithMeasurableMaterial sets the measurable material and returns the builder for chaining.
func (b *ComplexMaterialBuilder) WithMeasurableMaterial(m *MeasurableMaterial) *ComplexMaterialBuilder {
	b.complexMaterial.MeasurableMaterial = m
	return b
}

// Build returns the built complex material.
func (b *ComplexMaterialBuilder) Build() *ComplexMaterial {
	return b.complexMaterial
}
