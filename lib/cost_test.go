package lib

import "testing"

func TestGetTotalCost_Empty(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 5, nil)
	got := engine.GetTotalCost(root)
	assertFloatEqual(t, got, 0)
}

func TestGetTotalCost_WithHuman(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 60, nil) // 1 ora
	root.Humans = []HumanResource{{Role: "R", CostPerH: 30, Quantity: 1}}
	// (30/60) * 60 * 1 = 30
	got := engine.GetTotalCost(root)
	assertFloatEqual(t, got, 30)
}

func TestGetTotalCost_WithMaterial(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 5, nil)
	root.Materials = []MaterialResource{{Name: "M", UnitCost: 2.5, Quantity: 4}}
	// 2.5 * 4 = 10
	got := engine.GetTotalCost(root)
	assertFloatEqual(t, got, 10)
}

func TestGetTotalCost_Tree(t *testing.T) {
	engine := &AnalysisEngine{}
	child := buildActivity("B", "B", 3, nil)
	child.Materials = []MaterialResource{{UnitCost: 1, Quantity: 10}} // 10
	root := buildActivity("A", "A", 2, []*Activity{child})
	root.Humans = []HumanResource{{CostPerH: 60, Quantity: 1}} // (60/60)*2*1 = 2
	got := engine.GetTotalCost(root)
	assertFloatEqual(t, got, 12)
}

func TestGetCostBreakdown(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 60, nil)
	root.Humans = []HumanResource{{CostPerH: 60, Quantity: 1}}   // 60
	root.Materials = []MaterialResource{{UnitCost: 5, Quantity: 2}} // 10
	root.Assets = []Asset{{CostPerUse: 1, Quantity: 3}}           // 3
	b := engine.GetCostBreakdown(root)
	assertFloatEqual(t, b.Human, 60)
	assertFloatEqual(t, b.Material, 10)
	assertFloatEqual(t, b.Asset, 3)
}

func TestGetBreakEvenPrice(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{{UnitCost: 12, Quantity: 1}}
	got := engine.GetBreakEvenPrice(root)
	assertFloatEqual(t, got, 12)
	if engine.GetBreakEvenPrice(nil) != 0 {
		t.Error("GetBreakEvenPrice(nil) want 0")
	}
}
