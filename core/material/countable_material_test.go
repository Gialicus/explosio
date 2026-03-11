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

func TestCountableMaterial_SetTotalPrice(t *testing.T) {
	t.Run("quantity > 0 derives unit price", func(t *testing.T) {
		c := NewCountableMaterial("Screws", "", unit.Price{}, 100)
		c.SetTotalPrice(*unit.NewPrice(50, "EUR"))
		if c.Price.Value != 0.5 {
			t.Errorf("SetTotalPrice: unit price = %v, want 0.5", c.Price.Value)
		}
		if c.Price.Currency != "EUR" {
			t.Errorf("SetTotalPrice: currency = %q, want EUR", c.Price.Currency)
		}
		if c.CalculatePrice() != 50 {
			t.Errorf("CalculatePrice() = %v, want 50", c.CalculatePrice())
		}
	})
	t.Run("quantity 0 sets unit price to 0", func(t *testing.T) {
		c := NewCountableMaterial("Empty", "", unit.Price{}, 0)
		c.SetTotalPrice(*unit.NewPrice(100, "EUR"))
		if c.Price.Value != 0 {
			t.Errorf("SetTotalPrice with quantity 0: unit price = %v, want 0", c.Price.Value)
		}
		if c.Price.Currency != "EUR" {
			t.Errorf("SetTotalPrice: currency = %q, want EUR", c.Price.Currency)
		}
	})
}

func TestCountableMaterialBuilder_WithTotalPrice(t *testing.T) {
	c := NewCountableMaterialBuilder().
		WithName("Viti").
		WithQuantity(100).
		WithTotalPrice(*unit.NewPrice(50, "EUR")).
		Build()
	if c.Price.Value != 0.5 {
		t.Errorf("WithTotalPrice: unit price = %v, want 0.5", c.Price.Value)
	}
	if c.CalculatePrice() != 50 {
		t.Errorf("CalculatePrice() = %v, want 50", c.CalculatePrice())
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
