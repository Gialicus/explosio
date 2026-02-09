package core

import (
	"explosio/core/asset"
	"explosio/core/human"
	"explosio/core/material"
	"explosio/core/unit"
	"testing"
)

// buildQueryTestTree returns a small tree: root (R) with two children (A, B).
// R has no materials/resources; A has 1 countable, 1 measurable; B has 1 complex, 1 human, 1 asset.
// Used to test Get* methods.
func buildQueryTestTree(t *testing.T) *Activity {
	t.Helper()
	root := NewActivity("R", "Root", *unit.NewDuration(0, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	a := NewActivity("A", "Child A", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(10, "EUR"))
	b := NewActivity("B", "Child B", *unit.NewDuration(2, unit.DurationUnitDay), *unit.NewPrice(20, "EUR"))
	root.AddActivity(a)
	root.AddActivity(b)

	a.AddCountableMaterial(material.NewCountableMaterial("Screws", "", *unit.NewPrice(5, "EUR"), 10))
	a.AddMeasurableMaterial(material.NewMeasurableMaterial("Cable", "", *unit.NewPrice(3, "EUR"), *unit.NewMeasurableQuantity(10, unit.UnitMeter)))

	measUnit := material.NewMeasurableMaterialBuilder().
		WithName("Pipe").
		WithPrice(*unit.NewPrice(10, "EUR")).
		WithQuantity(*unit.NewMeasurableQuantity(2, unit.UnitMeter)).
		Build()
	complexMat := material.NewComplexMaterialBuilder().
		WithName("Pipes").
		WithPrice(*unit.NewPrice(50, "EUR")).
		WithUnitQuantity(5).
		WithMeasurableMaterial(measUnit).
		Build()
	b.AddComplexMaterial(complexMat)
	b.AddHumanResource(human.NewHumanResource("Plumber", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(100, "EUR")))
	b.AddAsset(asset.NewAsset("Tool", "", *unit.NewPrice(30, "EUR"), *unit.NewDuration(0, unit.DurationUnitDay)))

	return root
}

func TestActivity_GetActivities(t *testing.T) {
	root := buildQueryTestTree(t)
	got := root.GetActivities()
	wantNames := []string{"R", "A", "B"}
	if len(got) != len(wantNames) {
		t.Fatalf("GetActivities() length = %d, want %d", len(got), len(wantNames))
	}
	for i, name := range wantNames {
		if got[i].Name != name {
			t.Errorf("GetActivities()[%d].Name = %q, want %q", i, got[i].Name, name)
		}
	}
	if got[0] != root {
		t.Error("GetActivities()[0] should be the root activity")
	}
}

func TestActivity_GetComplexMaterials(t *testing.T) {
	root := buildQueryTestTree(t)
	got := root.GetComplexMaterials()
	if len(got) != 1 {
		t.Fatalf("GetComplexMaterials() length = %d, want 1", len(got))
	}
	if got[0].Name != "Pipes" {
		t.Errorf("GetComplexMaterials()[0].Name = %q, want Pipes", got[0].Name)
	}
}

func TestActivity_GetCountableMaterials(t *testing.T) {
	root := buildQueryTestTree(t)
	got := root.GetCountableMaterials()
	if len(got) != 1 {
		t.Fatalf("GetCountableMaterials() length = %d, want 1", len(got))
	}
	if got[0].Name != "Screws" {
		t.Errorf("GetCountableMaterials()[0].Name = %q, want Screws", got[0].Name)
	}
}

func TestActivity_GetMeasurableMaterials(t *testing.T) {
	root := buildQueryTestTree(t)
	got := root.GetMeasurableMaterials()
	if len(got) != 1 {
		t.Fatalf("GetMeasurableMaterials() length = %d, want 1", len(got))
	}
	if got[0].Name != "Cable" {
		t.Errorf("GetMeasurableMaterials()[0].Name = %q, want Cable", got[0].Name)
	}
}

func TestActivity_GetHumanResources(t *testing.T) {
	root := buildQueryTestTree(t)
	got := root.GetHumanResources()
	if len(got) != 1 {
		t.Fatalf("GetHumanResources() length = %d, want 1", len(got))
	}
	if got[0].Name != "Plumber" {
		t.Errorf("GetHumanResources()[0].Name = %q, want Plumber", got[0].Name)
	}
}

func TestActivity_GetAssets(t *testing.T) {
	root := buildQueryTestTree(t)
	got := root.GetAssets()
	if len(got) != 1 {
		t.Fatalf("GetAssets() length = %d, want 1", len(got))
	}
	if got[0].Name != "Tool" {
		t.Errorf("GetAssets()[0].Name = %q, want Tool", got[0].Name)
	}
}

func TestActivity_GetActivities_Leaf(t *testing.T) {
	leaf := NewActivity("Leaf", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	got := leaf.GetActivities()
	if len(got) != 1 || got[0].Name != "Leaf" {
		t.Errorf("GetActivities() on leaf = %d items, want 1 with Name Leaf", len(got))
	}
}

func TestActivity_GetComplexMaterials_Empty(t *testing.T) {
	// Single activity with no materials returns no complex materials.
	empty := NewActivity("E", "", *unit.NewDuration(0, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	got := empty.GetComplexMaterials()
	if len(got) != 0 {
		t.Errorf("GetComplexMaterials() on activity with none = %d, want 0", len(got))
	}
}
