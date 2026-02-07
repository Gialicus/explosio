package unit

import "testing"

func TestDuration_ToHours(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		u     DurationUnit
		want  float64
	}{
		{"zero minutes", 0, DurationUnitMinute, 0},
		{"60 minutes", 60, DurationUnitMinute, 1},
		{"30 minutes", 30, DurationUnitMinute, 0.5},
		{"one hour", 1, DurationUnitHour, 1},
		{"zero hours", 0, DurationUnitHour, 0},
		{"one day", 1, DurationUnitDay, 24},
		{"two days", 2, DurationUnitDay, 48},
		{"one week", 1, DurationUnitWeek, 168},
		{"one month", 1, DurationUnitMonth, 720},
		{"one year", 1, DurationUnitYear, 8760},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDuration(tt.value, tt.u)
			got := d.ToHours()
			if got != tt.want {
				t.Errorf("ToHours() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_ToHours_NilReceiver(t *testing.T) {
	var d *Duration
	if got := d.ToHours(); got != 0 {
		t.Errorf("nil Duration.ToHours() = %v, want 0", got)
	}
}

func TestNewDuration(t *testing.T) {
	d := NewDuration(2, DurationUnitDay)
	if d == nil {
		t.Fatal("NewDuration returned nil")
	}
	if d.Value != 2 || d.Unit != DurationUnitDay {
		t.Errorf("NewDuration(2, day) = %+v", d)
	}
	if got := d.ToHours(); got != 48 {
		t.Errorf("NewDuration(2, day).ToHours() = %v, want 48", got)
	}
}
