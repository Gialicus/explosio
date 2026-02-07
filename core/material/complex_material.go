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
