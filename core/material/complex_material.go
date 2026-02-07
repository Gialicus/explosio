package material

import "explosio/core/unit"

// ComplexMaterial is a material that is composed of multiple measurable materials
// for example: a set of 5 pipes 1 meter long
type ComplexMaterial struct {
	Name               string
	Description        string
	Price              unit.Price
	UnitQuantity       int
	MeasurableMaterial *MeasurableMaterial
}

func NewComplexMaterial(name string, description string, price unit.Price, unitQuantity int, measurableMaterial *MeasurableMaterial) *ComplexMaterial {
	return &ComplexMaterial{Name: name, Description: description, Price: price, UnitQuantity: unitQuantity, MeasurableMaterial: measurableMaterial}
}

func (c *ComplexMaterial) SetName(name string) *ComplexMaterial {
	c.Name = name
	return c
}

func (c *ComplexMaterial) SetDescription(description string) *ComplexMaterial {
	c.Description = description
	return c
}

func (c *ComplexMaterial) SetPrice(price unit.Price) *ComplexMaterial {
	c.Price = price
	return c
}

func (c *ComplexMaterial) SetUnitQuantity(unitQuantity int) *ComplexMaterial {
	c.UnitQuantity = unitQuantity
	return c
}

func (c *ComplexMaterial) SetMeasurableMaterial(measurableMaterial *MeasurableMaterial) *ComplexMaterial {
	c.MeasurableMaterial = measurableMaterial
	return c
}

func (c *ComplexMaterial) CalculatePrice() float64 {
	price := c.Price.Value
	price += c.MeasurableMaterial.CalculatePrice()
	return price
}
