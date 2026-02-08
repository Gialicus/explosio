package core

import (
	"explosio/core/asset"
	"explosio/core/human"
	"explosio/core/material"
	"explosio/core/unit"
	"testing"
)

func activityWithDefaults(name, desc string) *Activity {
	return NewActivity(name, desc, *unit.NewDuration(0, unit.DurationUnitHour), *unit.NewPrice(0, "EUR"))
}

func TestActivity_CalculatePrice(t *testing.T) {
	t.Run("own price only", func(t *testing.T) {
		a := activityWithDefaults("A", "")
		a.Price = *unit.NewPrice(100, "EUR")
		if got := a.CalculatePrice(); got != 100 {
			t.Errorf("CalculatePrice() = %v, want 100", got)
		}
	})
	t.Run("with one child", func(t *testing.T) {
		root := activityWithDefaults("Root", "")
		root.Price = *unit.NewPrice(10, "EUR")
		child := activityWithDefaults("Child", "")
		child.Price = *unit.NewPrice(20, "EUR")
		root.AddActivity(child)
		if got := root.CalculatePrice(); got != 30 {
			t.Errorf("CalculatePrice() = %v, want 30", got)
		}
	})
	t.Run("with materials", func(t *testing.T) {
		a := activityWithDefaults("A", "")
		a.Price = *unit.NewPrice(50, "EUR")
		a.AddCountableMaterial(material.NewCountableMaterial("S", "", *unit.NewPrice(5, "EUR"), 10))
		a.AddMeasurableMaterial(material.NewMeasurableMaterial("M", "", *unit.NewPrice(3, "EUR"), *unit.NewMeasurableQuantity(1, unit.UnitMeter)))
		// 50 + 5 + 3 = 58
		if got := a.CalculatePrice(); got != 58 {
			t.Errorf("CalculatePrice() = %v, want 58", got)
		}
	})
	t.Run("with human resources", func(t *testing.T) {
		a := activityWithDefaults("A", "")
		a.Price = *unit.NewPrice(10, "EUR")
		a.AddHumanResource(human.NewHumanResource("Dev", "", *unit.NewDuration(8, unit.DurationUnitHour), *unit.NewPrice(25, "EUR")))
		// 10 + 25 = 35
		if got := a.CalculatePrice(); got != 35 {
			t.Errorf("CalculatePrice() = %v, want 35", got)
		}
	})
	t.Run("with assets", func(t *testing.T) {
		a := activityWithDefaults("A", "")
		a.Price = *unit.NewPrice(10, "EUR")
		a.AddAsset(asset.NewAsset("Mac", "", *unit.NewPrice(15, "EUR"), *unit.NewDuration(1, unit.DurationUnitDay)))
		// 10 + 15 = 25
		if got := a.CalculatePrice(); got != 25 {
			t.Errorf("CalculatePrice() = %v, want 25", got)
		}
	})
	t.Run("with human resources and assets", func(t *testing.T) {
		a := activityWithDefaults("A", "")
		a.Price = *unit.NewPrice(5, "EUR")
		a.AddHumanResource(human.NewHumanResource("Dev", "", *unit.NewDuration(8, unit.DurationUnitHour), *unit.NewPrice(20, "EUR")))
		a.AddAsset(asset.NewAsset("Tool", "", *unit.NewPrice(10, "EUR"), *unit.NewDuration(1, unit.DurationUnitDay)))
		// 5 + 20 + 10 = 35
		if got := a.CalculatePrice(); got != 35 {
			t.Errorf("CalculatePrice() = %v, want 35", got)
		}
	})
}

func TestActivity_CalculateDuration(t *testing.T) {
	t.Run("no children", func(t *testing.T) {
		a := activityWithDefaults("A", "")
		a.Duration = *unit.NewDuration(5, unit.DurationUnitDay)
		if got := a.CalculateDuration(); got != 5 {
			t.Errorf("CalculateDuration() = %v, want 5", got)
		}
	})
	t.Run("with children", func(t *testing.T) {
		root := activityWithDefaults("Root", "")
		root.Duration = *unit.NewDuration(1, unit.DurationUnitDay)
		child := activityWithDefaults("Child", "")
		child.Duration = *unit.NewDuration(2, unit.DurationUnitDay)
		root.AddActivity(child)
		if got := root.CalculateDuration(); got != 3 {
			t.Errorf("CalculateDuration() = %v, want 3", got)
		}
	})
}

func TestActivity_CalculateCriticalPath(t *testing.T) {
	tests := []struct {
		name      string
		build     func() *Activity
		wantNames []string
	}{
		{
			name: "leaf",
			build: func() *Activity {
				a := activityWithDefaults("Leaf", "")
				a.Duration = *unit.NewDuration(1, unit.DurationUnitDay)
				return a
			},
			wantNames: []string{"Leaf"},
		},
		{
			name: "root and one child",
			build: func() *Activity {
				root := activityWithDefaults("Root", "")
				root.Duration = *unit.NewDuration(1, unit.DurationUnitDay)
				child := activityWithDefaults("Child", "")
				child.Duration = *unit.NewDuration(2, unit.DurationUnitDay)
				root.AddActivity(child)
				return root
			},
			wantNames: []string{"Root", "Child"},
		},
		{
			name: "two children longer path wins",
			build: func() *Activity {
				root := activityWithDefaults("Root", "")
				root.Duration = *unit.NewDuration(0, unit.DurationUnitDay)
				short := activityWithDefaults("Short", "")
				short.Duration = *unit.NewDuration(1, unit.DurationUnitDay)
				long := activityWithDefaults("Long", "")
				long.Duration = *unit.NewDuration(3, unit.DurationUnitDay)
				root.AddActivity(short)
				root.AddActivity(long)
				return root
			},
			wantNames: []string{"Root", "Long"},
		},
		{
			name: "nested longer branch",
			build: func() *Activity {
				root := activityWithDefaults("Root", "")
				root.Duration = *unit.NewDuration(0, unit.DurationUnitDay)
				a := activityWithDefaults("A", "")
				a.Duration = *unit.NewDuration(1, unit.DurationUnitDay)
				b := activityWithDefaults("B", "")
				b.Duration = *unit.NewDuration(1, unit.DurationUnitDay)
				deep := activityWithDefaults("Deep", "")
				deep.Duration = *unit.NewDuration(2, unit.DurationUnitDay)
				b.AddActivity(deep)
				root.AddActivity(a)
				root.AddActivity(b)
				return root
			},
			wantNames: []string{"Root", "B", "Deep"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.build()
			path := root.CalculateCriticalPath()
			if len(path) != len(tt.wantNames) {
				t.Fatalf("CalculateCriticalPath() length = %d, want %d", len(path), len(tt.wantNames))
			}
			for i, wantName := range tt.wantNames {
				if path[i].Name != wantName {
					t.Errorf("path[%d].Name = %q, want %q", i, path[i].Name, wantName)
				}
			}
		})
	}
}
