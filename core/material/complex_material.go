package material

import "explosio/core/unit"

// ComplexMaterial is a material that is composed of multiple measurable materials
// for example: a set of 5 pipes 1 meter long
type ComplexMaterial struct {
	Name                string
	Description         string
	Price               unit.Price
	UnitQuantity        int
	MeasurableMaterials []MeasurableMaterial
}

func NewComplexMaterial(name string, description string, price unit.Price, unitQuantity int, measurableMaterials []MeasurableMaterial) *ComplexMaterial {
	return &ComplexMaterial{Name: name, Description: description, Price: price, UnitQuantity: unitQuantity, MeasurableMaterials: measurableMaterials}
}

func (c *ComplexMaterial) CalculatePrice() float64 {
	price := c.Price.Value
	for _, measurableMaterial := range c.MeasurableMaterials {
		price += measurableMaterial.CalculatePrice()
	}
	return price
}
