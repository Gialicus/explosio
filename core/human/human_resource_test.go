package human

import (
	"explosio/core/unit"
	"testing"
)

func TestHumanResource_SetTotalPrice(t *testing.T) {
	h := NewHumanResource("Plumber", "", *unit.NewDuration(1, unit.DurationUnitDay), unit.Price{})
	h.SetTotalPrice(*unit.NewPrice(100, "EUR"))
	if h.Price.Value != 100 || h.Price.Currency != "EUR" {
		t.Errorf("SetTotalPrice: Price = %+v, want 100 EUR", h.Price)
	}
}

func TestHumanResource_SetHourlyRate(t *testing.T) {
	t.Run("duration > 0 derives total price", func(t *testing.T) {
		h := NewHumanResource("Plumber", "", *unit.NewDuration(1, unit.DurationUnitDay), unit.Price{})
		h.SetHourlyRate(*unit.NewPrice(10, "EUR"))
		// 1 day = 24 hours, so Price = 10 * 24 = 240
		if h.Price.Value != 240 {
			t.Errorf("SetHourlyRate: Price = %v, want 240", h.Price.Value)
		}
		if h.CalculateHourlyRate() != 10 {
			t.Errorf("CalculateHourlyRate() = %v, want 10", h.CalculateHourlyRate())
		}
	})
	t.Run("duration 0 sets total price to 0", func(t *testing.T) {
		h := NewHumanResource("Plumber", "", *unit.NewDuration(0, unit.DurationUnitHour), unit.Price{})
		h.SetHourlyRate(*unit.NewPrice(100, "EUR"))
		if h.Price.Value != 0 {
			t.Errorf("SetHourlyRate with duration 0: Price = %v, want 0", h.Price.Value)
		}
	})
}

func TestHumanResource_SetDailyRate(t *testing.T) {
	t.Run("duration > 0 derives total price", func(t *testing.T) {
		h := NewHumanResource("Plumber", "", *unit.NewDuration(8, unit.DurationUnitHour), unit.Price{})
		h.SetDailyRate(*unit.NewPrice(80, "EUR"))
		// 8 hours = 1 working day, so Price = 80 * 8/8 = 80
		if h.Price.Value != 80 {
			t.Errorf("SetDailyRate: Price = %v, want 80", h.Price.Value)
		}
		if h.CalculateDailyRate() != 80 {
			t.Errorf("CalculateDailyRate() = %v, want 80", h.CalculateDailyRate())
		}
	})
	t.Run("duration 0 sets total price to 0", func(t *testing.T) {
		h := NewHumanResource("Plumber", "", *unit.NewDuration(0, unit.DurationUnitHour), unit.Price{})
		h.SetDailyRate(*unit.NewPrice(100, "EUR"))
		if h.Price.Value != 0 {
			t.Errorf("SetDailyRate with duration 0: Price = %v, want 0", h.Price.Value)
		}
	})
}

func TestHumanResourceBuilder_WithTotalPrice(t *testing.T) {
	h := NewHumanResourceBuilder().
		WithName("Plumber").
		WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
		WithTotalPrice(*unit.NewPrice(100, "EUR")).
		Build()
	if h.Price.Value != 100 {
		t.Errorf("WithTotalPrice: Price = %v, want 100", h.Price.Value)
	}
}

func TestHumanResourceBuilder_WithHourlyRate(t *testing.T) {
	h := NewHumanResourceBuilder().
		WithName("Plumber").
		WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
		WithHourlyRate(*unit.NewPrice(10, "EUR")).
		Build()
	if h.Price.Value != 240 {
		t.Errorf("WithHourlyRate: Price = %v, want 240 (10 * 24)", h.Price.Value)
	}
}

func TestHumanResourceBuilder_WithDailyRate(t *testing.T) {
	h := NewHumanResourceBuilder().
		WithName("Plumber").
		WithDuration(*unit.NewDuration(8, unit.DurationUnitHour)).
		WithDailyRate(*unit.NewPrice(80, "EUR")).
		Build()
	if h.Price.Value != 80 {
		t.Errorf("WithDailyRate: Price = %v, want 80", h.Price.Value)
	}
}
