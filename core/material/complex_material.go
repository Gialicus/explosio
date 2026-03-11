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

// CalculatePrice returns the complex price plus the measurable material price multiplied by unit quantity.
// If MeasurableMaterial is nil, returns only the complex Price.Value.
func (c *ComplexMaterial) CalculatePrice() float64 {
	price := c.Price.Value
	if c.MeasurableMaterial != nil {
		price += float64(c.UnitQuantity) * c.MeasurableMaterial.CalculatePrice()
	}
	return price
}

// Clone returns a deep copy of the complex material.
func (c *ComplexMaterial) Clone() *ComplexMaterial {
	var meas *MeasurableMaterial
	if c.MeasurableMaterial != nil {
		meas = c.MeasurableMaterial.Clone()
	}
	return NewComplexMaterial(c.Name, c.Description, c.Price, c.UnitQuantity, meas)
}
