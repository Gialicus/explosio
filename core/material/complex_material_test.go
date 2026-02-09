package material

import (
	"explosio/core/unit"
	"testing"
)

func TestComplexMaterial_CalculatePrice_NilMeasurableMaterial(t *testing.T) {
	c := NewComplexMaterial("Only", "No measurable", *unit.NewPrice(50, "EUR"), 1, nil)
	got := c.CalculatePrice()
	if got != 50 {
		t.Errorf("CalculatePrice() with nil MeasurableMaterial = %v, want 50", got)
	}
}

func TestComplexMaterial_CalculatePrice(t *testing.T) {
	measPrice := *unit.NewPrice(10, "EUR")
	measQty := *unit.NewMeasurableQuantity(2, unit.UnitMeter)
	meas := NewMeasurableMaterial("Pipe 2m", "Copper pipe", measPrice, measQty)
	complexPrice := *unit.NewPrice(50, "EUR")
	c := NewComplexMaterial("Pipes", "Plumbing", complexPrice, 5, meas)
	got := c.CalculatePrice()
	// Complex price (50) + measurable material price (10)
	want := 60.0
	if got != want {
		t.Errorf("CalculatePrice() = %v, want %v", got, want)
	}
}

func TestComplexMaterialBuilder_Build(t *testing.T) {
	meas := NewMeasurableMaterialBuilder().
		WithName("Unit").
		WithPrice(*unit.NewPrice(5, "EUR")).
		WithQuantity(*unit.NewMeasurableQuantity(1, unit.UnitMeter)).
		Build()
	c := NewComplexMaterialBuilder().
		WithName("Bundle").
		WithDescription("Bundle of units").
		WithPrice(*unit.NewPrice(20, "EUR")).
		WithUnitQuantity(3).
		WithMeasurableMaterial(meas).
		Build()
	if c == nil {
		t.Fatal("Build() returned nil")
	}
	if c.Name != "Bundle" || c.UnitQuantity != 3 || c.MeasurableMaterial != meas {
		t.Errorf("Build() fields: name=%q unitQty=%d meas=%p", c.Name, c.UnitQuantity, c.MeasurableMaterial)
	}
	// 20 + 5 = 25
	if c.CalculatePrice() != 25 {
		t.Errorf("CalculatePrice() = %v, want 25", c.CalculatePrice())
	}
}
