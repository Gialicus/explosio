package material

import (
	"explosio/core/unit"
	"testing"
)

func TestCountableMaterial_CalculatePrice(t *testing.T) {
	price := *unit.NewPrice(0.5, "EUR")
	c := NewCountableMaterial("Screws", "Mounting screws", price, 100)
	got := c.CalculatePrice()
	// Total price = unit price * quantity
	if got != 50 {
		t.Errorf("CalculatePrice() = %v, want 50", got)
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
	if c.CalculatePrice() != 100 {
		t.Errorf("CalculatePrice() = %v, want 100 (25 * 4)", c.CalculatePrice())
	}
}
