package resources

import (
	"explosio/lib/domain"
	"testing"
)

func TestForEachResource_Empty(t *testing.T) {
	a := domain.BuildActivityForTest("A", "A", 1, nil)
	n := 0
	ForEachResource(a, func(domain.Resource) { n++ })
	if n != 0 {
		t.Errorf("ForEachResource(no resources) called callback %d times, want 0", n)
	}
}

func TestForEachResource_WithResources(t *testing.T) {
	a := domain.BuildActivityForTest("A", "A", 1, nil)
	a.Humans = []domain.HumanResource{{Role: "R1", CostPerH: 10, Quantity: 1}}
	a.Materials = []domain.MaterialResource{{Name: "M1", UnitCost: 5, Quantity: 2}}
	var count int
	ForEachResource(a, func(r domain.Resource) { count++ })
	if count != 2 {
		t.Errorf("ForEachResource got %d resources, want 2", count)
	}
}

func TestWalkResources_Count(t *testing.T) {
	c := domain.BuildActivityForTest("C", "C", 1, nil)
	c.Humans = []domain.HumanResource{{Role: "H", CostPerH: 1, Quantity: 1}}
	root := domain.BuildActivityForTest("R", "R", 1, []*domain.Activity{c})
	var count int
	WalkResources(root, func(*domain.Activity, domain.Resource) { count++ })
	if count != 1 {
		t.Errorf("WalkResources got %d resources, want 1", count)
	}
}

func TestResourceDisplayName(t *testing.T) {
	if name := ResourceDisplayName(domain.HumanResource{Role: "Dev"}); name != "Dev" {
		t.Errorf("ResourceDisplayName(Human) = %q, want Dev", name)
	}
	if name := ResourceDisplayName(domain.MaterialResource{Name: "Caffè"}); name != "Caffè" {
		t.Errorf("ResourceDisplayName(Material) = %q, want Caffè", name)
	}
	if name := ResourceDisplayName(domain.Asset{Name: "Macchina"}); name != "Macchina" {
		t.Errorf("ResourceDisplayName(Asset) = %q, want Macchina", name)
	}
}
