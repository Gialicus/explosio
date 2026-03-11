package material

import (
	"explosio/core/unit"
	"testing"
)

func TestMeasurableMaterial_CalculatePrice(t *testing.T) {
	price := *unit.NewPrice(10.5, "EUR")
	qty := *unit.NewMeasurableQuantity(5, unit.UnitKilogram)
	m := NewMeasurableMaterial("Cement", "Bags", price, qty)
	got := m.CalculatePrice()
	// Total price = unit price * quantity
	if got != 52.5 {
		t.Errorf("CalculatePrice() = %v, want 52.5", got)
	}
}

func TestMeasurableMaterial_SetTotalPrice(t *testing.T) {
	t.Run("quantity > 0 derives unit price", func(t *testing.T) {
		m := NewMeasurableMaterial("Cable", "", unit.Price{}, *unit.NewMeasurableQuantity(10, unit.UnitMeter))
		m.SetTotalPrice(*unit.NewPrice(200, "EUR"))
		if m.Price.Value != 20 {
			t.Errorf("SetTotalPrice: unit price = %v, want 20", m.Price.Value)
		}
		if m.Price.Currency != "EUR" {
			t.Errorf("SetTotalPrice: currency = %q, want EUR", m.Price.Currency)
		}
		if m.CalculatePrice() != 200 {
			t.Errorf("CalculatePrice() = %v, want 200", m.CalculatePrice())
		}
	})
	t.Run("quantity 0 sets unit price to 0", func(t *testing.T) {
		m := NewMeasurableMaterial("Empty", "", unit.Price{}, *unit.NewMeasurableQuantity(0, unit.UnitMeter))
		m.SetTotalPrice(*unit.NewPrice(100, "EUR"))
		if m.Price.Value != 0 {
			t.Errorf("SetTotalPrice with quantity 0: unit price = %v, want 0", m.Price.Value)
		}
		if m.Price.Currency != "EUR" {
			t.Errorf("SetTotalPrice: currency = %q, want EUR", m.Price.Currency)
		}
	})
}

func TestMeasurableMaterialBuilder_WithTotalPrice(t *testing.T) {
	m, err := NewMeasurableMaterialBuilder().
		WithName("Cable").
		WithQuantity(*unit.NewMeasurableQuantity(10, unit.UnitMeter)).
		WithTotalPrice(*unit.NewPrice(200, "EUR")).
		Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if m.Price.Value != 20 {
		t.Errorf("WithTotalPrice: unit price = %v, want 20", m.Price.Value)
	}
	if m.CalculatePrice() != 200 {
		t.Errorf("CalculatePrice() = %v, want 200", m.CalculatePrice())
	}
}

func TestMeasurableMaterialBuilder_Build(t *testing.T) {
	price := *unit.NewPrice(20, "EUR")
	qty := *unit.NewMeasurableQuantity(10, unit.UnitMeter)
	m, err := NewMeasurableMaterialBuilder().
		WithName("Cable").
		WithDescription("Copper cable").
		WithPrice(price).
		WithQuantity(qty).
		Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if m == nil {
		t.Fatal("Build() returned nil")
	}
	if m.Name != "Cable" || m.Description != "Copper cable" {
		t.Errorf("Build() name/description = %q, %q", m.Name, m.Description)
	}
	if m.Price.Value != 20 || m.Quantity.Value != 10 {
		t.Errorf("Build() price/quantity = %v, %v", m.Price, m.Quantity)
	}
	if m.CalculatePrice() != 200 {
		t.Errorf("CalculatePrice() = %v, want 200 (20 * 10)", m.CalculatePrice())
	}
}
