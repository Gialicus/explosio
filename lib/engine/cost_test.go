package engine

import (
	"explosio/lib/domain"
	"testing"
)

func TestGetTotalCost_Empty(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", 5, nil)
	got := eng.GetTotalCost(root)
	assertFloatEqual(t, got, 0)
}

func TestGetTotalCost_WithHuman(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", 60, nil)
	root.Humans = []domain.HumanResource{{Role: "R", CostPerH: 30, Quantity: 1}}
	got := eng.GetTotalCost(root)
	assertFloatEqual(t, got, 30)
}

func TestGetTotalCost_WithMaterial(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", 5, nil)
	root.Materials = []domain.MaterialResource{{Name: "M", UnitCost: 2.5, Quantity: 4}}
	got := eng.GetTotalCost(root)
	assertFloatEqual(t, got, 10)
}

func TestGetTotalCost_Tree(t *testing.T) {
	eng := &AnalysisEngine{}
	child := domain.BuildActivityForTest("B", "B", 3, nil)
	child.Materials = []domain.MaterialResource{{UnitCost: 1, Quantity: 10}}
	root := domain.BuildActivityForTest("A", "A", 2, []*domain.Activity{child})
	root.Humans = []domain.HumanResource{{CostPerH: 60, Quantity: 1}}
	got := eng.GetTotalCost(root)
	assertFloatEqual(t, got, 12)
}

func TestGetCostBreakdown(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", 60, nil)
	root.Humans = []domain.HumanResource{{CostPerH: 60, Quantity: 1}}
	root.Materials = []domain.MaterialResource{{UnitCost: 5, Quantity: 2}}
	root.Assets = []domain.Asset{{CostPerUse: 1, Quantity: 3}}
	b := eng.GetCostBreakdown(root)
	assertFloatEqual(t, b.Human, 60)
	assertFloatEqual(t, b.Material, 10)
	assertFloatEqual(t, b.Asset, 3)
}

func TestGetBreakEvenPrice(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	root.Materials = []domain.MaterialResource{{UnitCost: 12, Quantity: 1}}
	got := eng.GetBreakEvenPrice(root)
	assertFloatEqual(t, got, 12)
	if eng.GetBreakEvenPrice(nil) != 0 {
		t.Error("GetBreakEvenPrice(nil) want 0")
	}
}
