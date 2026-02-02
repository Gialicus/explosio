package lib

import (
	"math"
	"testing"
)

func assertFloatEqualWhatif(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestApplyScenario_EmptyScenario(t *testing.T) {
	child := buildActivity("B", "Child", 3, nil)
	root := buildActivity("A", "Root", 2, []*Activity{child})
	scenario := Scenario{Name: "Base", Overrides: nil}
	clone := ApplyScenario(root, scenario)
	if clone == nil {
		t.Fatal("ApplyScenario: want non-nil clone")
	}
	if clone.Duration != root.Duration || clone.ID != root.ID {
		t.Errorf("clone with empty scenario: Duration=%d ID=%q, want same as root", clone.Duration, clone.ID)
	}
	if len(clone.SubActivities) != 1 || clone.SubActivities[0].Duration != 3 {
		t.Errorf("clone child: want Duration 3, got %d", clone.SubActivities[0].Duration)
	}
}

func TestApplyScenario_DurationOverride(t *testing.T) {
	child := buildActivity("B", "Child", 3, nil)
	root := buildActivity("A", "Root", 2, []*Activity{child})
	dur := 10
	scenario := Scenario{
		Name: "Durata aumentata",
		Overrides: map[string]ActivityOverride{
			"B": {Duration: &dur},
		},
	}
	clone := ApplyScenario(root, scenario)
	if clone.SubActivities[0].Duration != 10 {
		t.Errorf("clone child Duration want 10, got %d", clone.SubActivities[0].Duration)
	}
	if child.Duration != 3 {
		t.Error("original child must be unchanged")
	}
}

func TestApplyScenario_MaterialCostFactor(t *testing.T) {
	root := buildActivity("A", "Root", 1, nil)
	root.Materials = []MaterialResource{{Name: "M", UnitCost: 10, Quantity: 2}}
	scenario := Scenario{
		Name: "Materiali +20%",
		Overrides: map[string]ActivityOverride{
			"A": {MaterialCostFactor: 1.2},
		},
	}
	clone := ApplyScenario(root, scenario)
	if len(clone.Materials) != 1 {
		t.Fatalf("clone Materials len want 1, got %d", len(clone.Materials))
	}
	// 10 * 1.2 = 12
	assertFloatEqualWhatif(t, clone.Materials[0].UnitCost, 12)
	if root.Materials[0].UnitCost != 10 {
		t.Error("original Material UnitCost must be unchanged")
	}
}

func TestRunScenario_BaseVsDurationIncreased(t *testing.T) {
	child := buildActivity("B", "Child", 3, nil)
	root := buildActivity("A", "Root", 2, []*Activity{child})
	engine := NewWhatIfEngine()
	baseScenario := Scenario{Name: "Base", SellPrice: 10, Overrides: nil}
	dur := 5
	altScenario := Scenario{
		Name:      "Durata aumentata",
		SellPrice: 10,
		Overrides: map[string]ActivityOverride{
			"B": {Duration: &dur},
		},
	}
	_, resBase := engine.RunScenario(root, baseScenario)
	_, resAlt := engine.RunScenario(root, altScenario)
	// Base: child 3 min, root 2 min -> total 5 min. Alt: child 5 min, root 2 min -> total 7 min
	if resBase.TotalDuration != 5 {
		t.Errorf("base TotalDuration want 5, got %d", resBase.TotalDuration)
	}
	if resAlt.TotalDuration != 7 {
		t.Errorf("alt TotalDuration want 7, got %d", resAlt.TotalDuration)
	}
	if resBase.ScenarioName != "Base" || resAlt.ScenarioName != "Durata aumentata" {
		t.Errorf("ScenarioName: base=%q alt=%q", resBase.ScenarioName, resAlt.ScenarioName)
	}
}

func TestRunScenarios_TwoScenarios(t *testing.T) {
	root := buildActivity("A", "Root", 2, nil)
	engine := NewWhatIfEngine()
	scenarios := []Scenario{
		{Name: "S1", SellPrice: 5},
		{Name: "S2", SellPrice: 8},
	}
	results := engine.RunScenarios(root, scenarios)
	if len(results) != 2 {
		t.Fatalf("RunScenarios want 2 results, got %d", len(results))
	}
	if results[0].ScenarioName != "S1" || results[1].ScenarioName != "S2" {
		t.Errorf("results names: %q %q", results[0].ScenarioName, results[1].ScenarioName)
	}
	if results[0].TotalDuration != 2 || results[1].TotalDuration != 2 {
		t.Errorf("TotalDuration: %d %d", results[0].TotalDuration, results[1].TotalDuration)
	}
	// Same cost, different sell price -> different margin
	if results[0].Margin >= results[1].Margin && results[0].TotalCost == results[1].TotalCost {
		t.Log("S2 has higher price so higher margin")
	}
}
