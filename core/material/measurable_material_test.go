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
	if got != 10.5 {
		t.Errorf("CalculatePrice() = %v, want 10.5", got)
	}
}

func TestMeasurableMaterialBuilder_Build(t *testing.T) {
	price := *unit.NewPrice(20, "EUR")
	qty := *unit.NewMeasurableQuantity(10, unit.UnitMeter)
	m := NewMeasurableMaterialBuilder().
		WithName("Cable").
		WithDescription("Copper cable").
		WithPrice(price).
		WithQuantity(qty).
		Build()
	if m == nil {
		t.Fatal("Build() returned nil")
	}
	if m.Name != "Cable" || m.Description != "Copper cable" {
		t.Errorf("Build() name/description = %q, %q", m.Name, m.Description)
	}
	if m.Price.Value != 20 || m.Quantity.Value != 10 {
		t.Errorf("Build() price/quantity = %v, %v", m.Price, m.Quantity)
	}
	if m.CalculatePrice() != 20 {
		t.Errorf("CalculatePrice() = %v, want 20", m.CalculatePrice())
	}
}
