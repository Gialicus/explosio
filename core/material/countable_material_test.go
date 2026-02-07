package material

import (
	"explosio/core/unit"
	"testing"
)

func TestCountableMaterial_CalculatePrice(t *testing.T) {
	price := *unit.NewPrice(0.5, "EUR")
	c := NewCountableMaterial("Screws", "Mounting screws", price, 100)
	got := c.CalculatePrice()
	// Unit price only, not multiplied by Quantity
	if got != 0.5 {
		t.Errorf("CalculatePrice() = %v, want 0.5", got)
	}
}

func TestCountableMaterialBuilder_Build(t *testing.T) {
	price := *unit.NewPrice(25, "EUR")
	c := NewCountableMaterialBuilder().
		WithName("Switches").
		WithDescription("Light switches").
		WithPrice(price).
		WithQuantity(4).
		Build()
	if c == nil {
		t.Fatal("Build() returned nil")
	}
	if c.Name != "Switches" || c.Quantity != 4 {
		t.Errorf("Build() name/quantity = %q, %d", c.Name, c.Quantity)
	}
	if c.CalculatePrice() != 25 {
		t.Errorf("CalculatePrice() = %v, want 25", c.CalculatePrice())
	}
}
