// Package material defines material types: complex, countable, and measurable.
package material

import "explosio/core/unit"

// ComplexMaterial is a material made of multiple units of a measurable material (e.g. 5 pipes of 1 meter).
type ComplexMaterial struct {
	Name               string
	Description        string
	Price              unit.Price
	UnitQuantity       int
	MeasurableMaterial *MeasurableMaterial
}

// NewComplexMaterial creates a complex material (e.g. N units of a measurable material).
func NewComplexMaterial(name string, description string, price unit.Price, unitQuantity int, measurableMaterial *MeasurableMaterial) *ComplexMaterial {
	return &ComplexMaterial{Name: name, Description: description, Price: price, UnitQuantity: unitQuantity, MeasurableMaterial: measurableMaterial}
}

// SetName sets the name and returns the receiver for chaining.
func (c *ComplexMaterial) SetName(name string) *ComplexMaterial {
	c.Name = name
	return c
}

// SetDescription sets the description and returns the receiver for chaining.
func (c *ComplexMaterial) SetDescription(description string) *ComplexMaterial {
	c.Description = description
	return c
}

// SetPrice sets the price and returns the receiver for chaining.
func (c *ComplexMaterial) SetPrice(price unit.Price) *ComplexMaterial {
	c.Price = price
	return c
}

// SetUnitQuantity sets the unit quantity and returns the receiver for chaining.
func (c *ComplexMaterial) SetUnitQuantity(unitQuantity int) *ComplexMaterial {
	c.UnitQuantity = unitQuantity
	return c
}

// SetMeasurableMaterial sets the measurable material and returns the receiver for chaining.
func (c *ComplexMaterial) SetMeasurableMaterial(measurableMaterial *MeasurableMaterial) *ComplexMaterial {
	c.MeasurableMaterial = measurableMaterial
	return c
}

// CalculatePrice returns the complex price plus the measurable material price.
func (c *ComplexMaterial) CalculatePrice() float64 {
	price := c.Price.Value
	price += c.MeasurableMaterial.CalculatePrice()
	return price
}
