package core

import (
	"explosio/core/material"
	"explosio/core/resource/asset"
	"explosio/core/resource/human"
	"explosio/core/unit"
	"testing"
)

func TestActivity_Clone(t *testing.T) {
	root := NewActivity("Root", "Test", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(100, "EUR"))
	child := NewActivity("Child", "Child", *unit.NewDuration(2, unit.DurationUnitDay), *unit.NewPrice(50, "EUR"))
	root.AddActivity(child)
	root.AddCountableMaterial(material.NewCountableMaterial("Screws", "", *unit.NewPrice(2, "EUR"), 10))
	root.AddHumanResource(human.NewHumanResource("Worker", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(80, "EUR")))
	root.AddAsset(asset.NewAsset("Tool", "", *unit.NewPrice(20, "EUR"), *unit.NewDuration(0, unit.DurationUnitDay)))

	clone := root.Clone()
	if clone == root {
		t.Fatal("Clone should return a different pointer")
	}
	if clone.Name != root.Name {
		t.Errorf("Clone.Name = %q, want %q", clone.Name, root.Name)
	}
	if clone.CalculatePrice() != root.CalculatePrice() {
		t.Errorf("Clone price = %.2f, want %.2f", clone.CalculatePrice(), root.CalculatePrice())
	}
	if len(clone.Activities) != len(root.Activities) {
		t.Errorf("Clone has %d children, want %d", len(clone.Activities), len(root.Activities))
	}
	if clone.Activities[0] == root.Activities[0] {
		t.Error("Clone's child should be a different pointer")
	}
}
