package lib

import "testing"

func TestGetFinancials_PositiveMargin(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{{UnitCost: 10, Quantity: 1}} // cost 10
	fin := engine.GetFinancials(root, 15)
	assertFloatEqual(t, fin.TotalCost, 10)
	assertFloatEqual(t, fin.Margin, 5)
	assertFloatEqual(t, fin.Markup, 50) // (5/10)*100
	if !fin.IsViable {
		t.Error("IsViable want true")
	}
}

func TestGetFinancials_NegativeMargin(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{{UnitCost: 20, Quantity: 1}}
	fin := engine.GetFinancials(root, 15)
	assertFloatEqual(t, fin.TotalCost, 20)
	assertFloatEqual(t, fin.Margin, -5)
	if fin.IsViable {
		t.Error("IsViable want false when margin < 0")
	}
}

func TestGetFinancials_ZeroCost(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 1, nil) // no resources, cost 0
	fin := engine.GetFinancials(root, 10)
	assertFloatEqual(t, fin.TotalCost, 0)
	assertFloatEqual(t, fin.Margin, 10)
	assertFloatEqual(t, fin.Markup, 0)
	if !fin.IsViable {
		t.Error("IsViable want true when margin > 0 and cost 0")
	}
}

func TestGetFinancialsForPrices(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{{UnitCost: 10, Quantity: 1}}
	prices := []float64{10, 15, 20}
	fin := engine.GetFinancialsForPrices(root, prices)
	if len(fin) != 3 {
		t.Fatalf("GetFinancialsForPrices want 3, got %d", len(fin))
	}
	assertFloatEqual(t, fin[0].Margin, 0)
	assertFloatEqual(t, fin[1].Margin, 5)
	assertFloatEqual(t, fin[2].Margin, 10)
}
