package core

import (
	"explosio/core/unit"
	"testing"
)

func TestActivityBuilder_Build(t *testing.T) {
	dur := *unit.NewDuration(2, unit.DurationUnitDay)
	price := *unit.NewPrice(1000, "EUR")
	a, err := NewActivityBuilder().
		WithName("Test").
		WithDescription("Test activity").
		WithDuration(dur).
		WithPrice(price).
		Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if a == nil {
		t.Fatal("Build() returned nil")
	}
	if a.Name != "Test" || a.Description != "Test activity" {
		t.Errorf("Build() name/description = %q, %q", a.Name, a.Description)
	}
	if a.Duration.Value != 2 || a.Duration.Unit != unit.DurationUnitDay {
		t.Errorf("Build() duration = %+v", a.Duration)
	}
	if a.Price.Value != 1000 || a.Price.Currency != "EUR" {
		t.Errorf("Build() price = %+v", a.Price)
	}
	if a.Activities != nil || len(a.Activities) != 0 {
		t.Errorf("Build() Activities should be empty, got len=%d", len(a.Activities))
	}
	if a.ComplexMaterials != nil || len(a.ComplexMaterials) != 0 {
		t.Errorf("Build() ComplexMaterials should be empty, got len=%d", len(a.ComplexMaterials))
	}
	if a.CountableMaterials != nil || len(a.CountableMaterials) != 0 {
		t.Errorf("Build() CountableMaterials should be empty, got len=%d", len(a.CountableMaterials))
	}
	if a.MeasurableMaterials != nil || len(a.MeasurableMaterials) != 0 {
		t.Errorf("Build() MeasurableMaterials should be empty, got len=%d", len(a.MeasurableMaterials))
	}
}

func TestActivityBuilder_Build_validation(t *testing.T) {
	t.Run("empty name returns error", func(t *testing.T) {
		_, err := NewActivityBuilder().
			WithName("").
			WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
			WithPrice(*unit.NewPrice(100, "EUR")).
			Build()
		if err == nil {
			t.Error("Build() with empty name should return error")
		}
	})
	t.Run("negative duration returns error", func(t *testing.T) {
		_, err := NewActivityBuilder().
			WithName("Test").
			WithDuration(*unit.NewDuration(-1, unit.DurationUnitDay)).
			WithPrice(*unit.NewPrice(100, "EUR")).
			Build()
		if err == nil {
			t.Error("Build() with negative duration should return error")
		}
	})
	t.Run("negative price returns error", func(t *testing.T) {
		_, err := NewActivityBuilder().
			WithName("Test").
			WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
			WithPrice(*unit.NewPrice(-10, "EUR")).
			Build()
		if err == nil {
			t.Error("Build() with negative price should return error")
		}
	})
	t.Run("empty currency returns error", func(t *testing.T) {
		_, err := NewActivityBuilder().
			WithName("Test").
			WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
			WithPrice(*unit.NewPrice(100, "")).
			Build()
		if err == nil {
			t.Error("Build() with empty currency should return error")
		}
	})
}
