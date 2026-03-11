package core

import (
	"bytes"
	"explosio/core/material"
	"explosio/core/resource/asset"
	"explosio/core/resource/human"
	"explosio/core/unit"
	"testing"
)

func buildSerializationTestTree(t *testing.T) *Activity {
	t.Helper()
	root := NewActivity("Root", "Root activity", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(100, "EUR"))
	child := NewActivity("Child", "Child activity", *unit.NewDuration(2, unit.DurationUnitDay), *unit.NewPrice(50, "EUR"))
	root.AddActivity(child)

	meas, err := material.NewMeasurableMaterialBuilder().
		WithName("Cable").
		WithPrice(*unit.NewPrice(2, "EUR")).
		WithQuantity(*unit.NewMeasurableQuantity(10, unit.UnitMeter)).
		Build()
	if err != nil {
		t.Fatalf("MeasurableMaterial Build: %v", err)
	}
	child.AddMeasurableMaterial(meas)
	child.AddHumanResource(human.NewHumanResource("Worker", "Test worker", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(80, "EUR")))
	child.AddAsset(asset.NewAsset("Tool", "Test tool", *unit.NewPrice(20, "EUR"), *unit.NewDuration(0, unit.DurationUnitDay)))

	return root
}

func TestProject_WriteReadJSON(t *testing.T) {
	root := buildSerializationTestTree(t)
	proj := NewProject(root)

	var buf bytes.Buffer
	if err := proj.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}

	read, err := ReadJSON(&buf)
	if err != nil {
		t.Fatalf("ReadJSON: %v", err)
	}

	if read.Version != ProjectVersion {
		t.Errorf("Version = %q, want %q", read.Version, ProjectVersion)
	}
	if read.Root == nil {
		t.Fatal("Root is nil after read")
	}
	if read.Root.Name != "Root" {
		t.Errorf("Root.Name = %q, want Root", read.Root.Name)
	}
	if len(read.Root.Activities) != 1 {
		t.Fatalf("Root has %d activities, want 1", len(read.Root.Activities))
	}
	child := read.Root.Activities[0]
	if child.Name != "Child" {
		t.Errorf("Child.Name = %q, want Child", child.Name)
	}
	if len(child.MeasurableMaterials) != 1 {
		t.Fatalf("Child has %d measurable materials, want 1", len(child.MeasurableMaterials))
	}
	if child.MeasurableMaterials[0].Name != "Cable" {
		t.Errorf("Material name = %q, want Cable", child.MeasurableMaterials[0].Name)
	}

	// Verify price calculation is preserved
	origPrice := root.CalculatePrice()
	readPrice := read.Root.CalculatePrice()
	if readPrice != origPrice {
		t.Errorf("Price after round-trip = %.2f, want %.2f", readPrice, origPrice)
	}
}

func TestProject_WriteReadYAML(t *testing.T) {
	root := buildSerializationTestTree(t)
	proj := NewProject(root)

	var buf bytes.Buffer
	if err := proj.WriteYAML(&buf); err != nil {
		t.Fatalf("WriteYAML: %v", err)
	}

	read, err := ReadYAML(&buf)
	if err != nil {
		t.Fatalf("ReadYAML: %v", err)
	}

	if read.Root == nil {
		t.Fatal("Root is nil after read")
	}
	if read.Root.Name != "Root" {
		t.Errorf("Root.Name = %q, want Root", read.Root.Name)
	}

	origPrice := root.CalculatePrice()
	readPrice := read.Root.CalculatePrice()
	if readPrice != origPrice {
		t.Errorf("Price after YAML round-trip = %.2f, want %.2f", readPrice, origPrice)
	}
}
