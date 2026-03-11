package unit

import "testing"

func TestNewMeasurableQuantity(t *testing.T) {
	q := NewMeasurableQuantity(5, UnitKilogram)
	if q == nil {
		t.Fatal("NewMeasurableQuantity returned nil")
	}
	if q.Value != 5 || q.Unit != UnitKilogram {
		t.Errorf("NewMeasurableQuantity(5, kg) = %+v", q)
	}
}

func TestMeasurableQuantity_String(t *testing.T) {
	tests := []struct {
		value float64
		unit  MeasurableUnit
		want  string
	}{
		{5, UnitMeter, "5 m"},
		{10, UnitKilogram, "10 kg"},
		{15, UnitSquareMeter, "15 m²"},
	}
	for _, tt := range tests {
		q := NewMeasurableQuantity(tt.value, tt.unit)
		got := q.String()
		if got != tt.want {
			t.Errorf("MeasurableQuantity(%v, %s).String() = %q, want %q", tt.value, tt.unit, got, tt.want)
		}
	}
}

func TestMeasurableQuantity_SetValue(t *testing.T) {
	q := NewMeasurableQuantity(5, UnitMeter)
	got := q.SetValue(10)
	if got != q {
		t.Error("SetValue should return self for chaining")
	}
	if q.Value != 10 {
		t.Errorf("SetValue(10): Value = %v, want 10", q.Value)
	}
}

func TestMeasurableQuantity_SetUnit(t *testing.T) {
	q := NewMeasurableQuantity(5, UnitMeter)
	got := q.SetUnit(UnitKilogram)
	if got != q {
		t.Error("SetUnit should return self for chaining")
	}
	if q.Unit != UnitKilogram {
		t.Errorf("SetUnit(kg): Unit = %v, want kg", q.Unit)
	}
}
