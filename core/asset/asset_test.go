package asset

import (
	"explosio/core/unit"
	"testing"
)

func TestAsset_SetTotalPrice(t *testing.T) {
	a := NewAsset("Tool", "", unit.Price{}, *unit.NewDuration(1, unit.DurationUnitDay))
	a.SetTotalPrice(*unit.NewPrice(100, "EUR"))
	if a.Price.Value != 100 || a.Price.Currency != "EUR" {
		t.Errorf("SetTotalPrice: Price = %+v, want 100 EUR", a.Price)
	}
}

func TestAsset_SetHourlyRate(t *testing.T) {
	t.Run("duration > 0 derives total price", func(t *testing.T) {
		a := NewAsset("Tool", "", unit.Price{}, *unit.NewDuration(1, unit.DurationUnitDay))
		a.SetHourlyRate(*unit.NewPrice(10, "EUR"))
		// 1 day = 24 hours, so Price = 10 * 24 = 240
		if a.Price.Value != 240 {
			t.Errorf("SetHourlyRate: Price = %v, want 240", a.Price.Value)
		}
		if a.CalculateHourlyRate() != 10 {
			t.Errorf("CalculateHourlyRate() = %v, want 10", a.CalculateHourlyRate())
		}
	})
	t.Run("duration 0 sets total price to 0", func(t *testing.T) {
		a := NewAsset("Tool", "", unit.Price{}, *unit.NewDuration(0, unit.DurationUnitHour))
		a.SetHourlyRate(*unit.NewPrice(100, "EUR"))
		if a.Price.Value != 0 {
			t.Errorf("SetHourlyRate with duration 0: Price = %v, want 0", a.Price.Value)
		}
	})
}

func TestAsset_SetDailyRate(t *testing.T) {
	t.Run("duration > 0 derives total price", func(t *testing.T) {
		a := NewAsset("Tool", "", unit.Price{}, *unit.NewDuration(8, unit.DurationUnitHour))
		a.SetDailyRate(*unit.NewPrice(80, "EUR"))
		// 8 hours = 1 working day, so Price = 80 * 8/8 = 80
		if a.Price.Value != 80 {
			t.Errorf("SetDailyRate: Price = %v, want 80", a.Price.Value)
		}
		if a.CalculateDailyRate() != 80 {
			t.Errorf("CalculateDailyRate() = %v, want 80", a.CalculateDailyRate())
		}
	})
	t.Run("duration 0 sets total price to 0", func(t *testing.T) {
		a := NewAsset("Tool", "", unit.Price{}, *unit.NewDuration(0, unit.DurationUnitHour))
		a.SetDailyRate(*unit.NewPrice(100, "EUR"))
		if a.Price.Value != 0 {
			t.Errorf("SetDailyRate with duration 0: Price = %v, want 0", a.Price.Value)
		}
	})
}

func TestAssetBuilder_WithTotalPrice(t *testing.T) {
	a := NewAssetBuilder().
		WithName("Tool").
		WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
		WithTotalPrice(*unit.NewPrice(100, "EUR")).
		Build()
	if a.Price.Value != 100 {
		t.Errorf("WithTotalPrice: Price = %v, want 100", a.Price.Value)
	}
}

func TestAssetBuilder_WithHourlyRate(t *testing.T) {
	a := NewAssetBuilder().
		WithName("Tool").
		WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
		WithHourlyRate(*unit.NewPrice(10, "EUR")).
		Build()
	if a.Price.Value != 240 {
		t.Errorf("WithHourlyRate: Price = %v, want 240 (10 * 24)", a.Price.Value)
	}
}

func TestAssetBuilder_WithDailyRate(t *testing.T) {
	a := NewAssetBuilder().
		WithName("Tool").
		WithDuration(*unit.NewDuration(8, unit.DurationUnitHour)).
		WithDailyRate(*unit.NewPrice(80, "EUR")).
		Build()
	if a.Price.Value != 80 {
		t.Errorf("WithDailyRate: Price = %v, want 80", a.Price.Value)
	}
}
