package core

import (
	"explosio/core/unit"
	"testing"
)

func TestValidate_NoErrors(t *testing.T) {
	root := NewActivity("Root", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(100, "EUR"))
	child := NewActivity("Child", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(50, "EUR"))
	root.AddActivity(child)

	r := root.Validate()
	if !r.Valid() {
		t.Errorf("Expected valid, got errors: %v", r.Errors)
	}
}

func TestValidate_CircularDependency(t *testing.T) {
	a := NewActivity("A", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	b := NewActivity("B", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	a.AddActivity(b)
	b.AddDependsOn(a) // B depends on A; A is parent of B - no cycle
	// A -> B, B depends on A. So B starts after A. No cycle.

	c := NewActivity("C", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	d := NewActivity("D", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	c.AddActivity(d)
	c.AddDependsOn(d) // C depends on D (its child) - cycle!
	d.AddDependsOn(c) // D depends on C - cycle!

	r := c.Validate()
	if r.Valid() {
		t.Error("Expected validation error for circular dependency")
	}
}
