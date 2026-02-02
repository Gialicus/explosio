package lib

import (
	"regexp"
	"testing"
)

func TestNewProject(t *testing.T) {
	p := NewProject()
	if p == nil {
		t.Fatal("NewProject() must not return nil")
	}
	if p.Root != nil {
		t.Error("Root must be nil before Start()")
	}
}

func TestStart(t *testing.T) {
	p := NewProject()
	root := p.Start("A", "desc", 5)
	if p.Root == nil {
		t.Fatal("Root must be set after Start()")
	}
	if root != p.Root {
		t.Error("Start() must return the same activity as Root")
	}
	if root.Name != "A" || root.Description != "desc" || root.Duration != 5 {
		t.Errorf("Name=%q Description=%q Duration=%d, want A desc 5", root.Name, root.Description, root.Duration)
	}
	if root.MinDuration != 5 {
		t.Errorf("MinDuration want 5, got %d", root.MinDuration)
	}
	actID := regexp.MustCompile(`^ACT-\d{3}$`)
	if !actID.MatchString(root.ID) {
		t.Errorf("ID must match ACT-00x, got %q", root.ID)
	}
}

func TestNode(t *testing.T) {
	p := NewProject()
	a := p.Node("B", "b", 3)
	if a == nil {
		t.Fatal("Node() must not return nil")
	}
	if a.Name != "B" || a.Description != "b" || a.Duration != 3 {
		t.Errorf("Name=%q Description=%q Duration=%d, want B b 3", a.Name, a.Description, a.Duration)
	}
	actID := regexp.MustCompile(`^ACT-\d{3}$`)
	if !actID.MatchString(a.ID) {
		t.Errorf("ID must match ACT-00x, got %q", a.ID)
	}
}

func TestUniqueIDs(t *testing.T) {
	p := NewProject()
	root := p.Start("R", "r", 1)
	n1 := p.Node("N1", "n1", 1)
	n2 := p.Node("N2", "n2", 1)
	ids := map[string]bool{root.ID: true, n1.ID: true, n2.ID: true}
	if len(ids) != 3 {
		t.Errorf("expected 3 unique IDs, got %d: %v", len(ids), ids)
	}
	for id := range ids {
		ok, _ := regexp.MatchString(`^ACT-\d{3}$`, id)
		if !ok {
			t.Errorf("ID %q does not match ACT-00x", id)
		}
	}
}

func TestWithHumanWithMaterialWithAsset(t *testing.T) {
	p := NewProject()
	a := p.Node("X", "x", 1).
		WithHuman("Role", "Human desc", 20, 1).
		WithMaterial("Mat", "Mat desc", 0.5, 10).
		WithAsset("Asset", "Asset desc", 1.0, 2)
	if len(a.Humans) != 1 {
		t.Fatalf("len(Humans) want 1, got %d", len(a.Humans))
	}
	if a.Humans[0].Role != "Role" || a.Humans[0].Description != "Human desc" || a.Humans[0].CostPerH != 20 || a.Humans[0].Quantity != 1 {
		t.Errorf("Humans[0]: Role=%q CostPerH=%v Quantity=%v", a.Humans[0].Role, a.Humans[0].CostPerH, a.Humans[0].Quantity)
	}
	if len(a.Materials) != 1 {
		t.Fatalf("len(Materials) want 1, got %d", len(a.Materials))
	}
	if a.Materials[0].Name != "Mat" || a.Materials[0].UnitCost != 0.5 || a.Materials[0].Quantity != 10 {
		t.Errorf("Materials[0]: Name=%q UnitCost=%v Quantity=%v", a.Materials[0].Name, a.Materials[0].UnitCost, a.Materials[0].Quantity)
	}
	if len(a.Assets) != 1 {
		t.Fatalf("len(Assets) want 1, got %d", len(a.Assets))
	}
	if a.Assets[0].Name != "Asset" || a.Assets[0].CostPerUse != 1.0 || a.Assets[0].Quantity != 2 {
		t.Errorf("Assets[0]: Name=%q CostPerUse=%v Quantity=%v", a.Assets[0].Name, a.Assets[0].CostPerUse, a.Assets[0].Quantity)
	}
}

func TestCanCrash(t *testing.T) {
	p := NewProject()
	a := p.Node("C", "c", 5).CanCrash(1, 10.0)
	if a.MinDuration != 1 {
		t.Errorf("MinDuration want 1, got %d", a.MinDuration)
	}
	if a.CrashCostStep != 10.0 {
		t.Errorf("CrashCostStep want 10.0, got %v", a.CrashCostStep)
	}
}

func TestDependsOn(t *testing.T) {
	p := NewProject()
	root := p.Start("R", "r", 2)
	child1 := p.Node("C1", "c1", 1)
	child2 := p.Node("C2", "c2", 1)
	root.DependsOn(child1, child2)
	if len(root.SubActivities) != 2 {
		t.Fatalf("len(SubActivities) want 2, got %d", len(root.SubActivities))
	}
	seen := make(map[string]bool)
	for _, sub := range root.SubActivities {
		seen[sub.ID] = true
		if len(sub.Next) != 1 || sub.Next[0] != root.ID {
			t.Errorf("child %q Next want [%q], got %v", sub.ID, root.ID, sub.Next)
		}
	}
	if !seen[child1.ID] || !seen[child2.ID] {
		t.Error("SubActivities must contain both child1 and child2")
	}
}

// TestDSLAndEngine_Integration costruisce un albero a tre livelli e verifica CPM.
func TestDSLAndEngine_Integration(t *testing.T) {
	p := NewProject()
	engine := &AnalysisEngine{}
	childL := p.Node("Left", "left", 2)
	childR := p.Node("Right", "right", 4)
	mid := p.Node("Mid", "mid", 1).DependsOn(childL, childR)
	root := p.Start("Root", "root", 2).DependsOn(mid)
	engine.ComputeCPM(root)
	// Left: 0-2, Right: 0-4, Mid: 4-5, Root: 5-7
	if childL.ES != 0 || childL.EF != 2 {
		t.Errorf("Left: ES=0 EF=2, got ES=%d EF=%d", childL.ES, childL.EF)
	}
	if childR.ES != 0 || childR.EF != 4 {
		t.Errorf("Right: ES=0 EF=4, got ES=%d EF=%d", childR.ES, childR.EF)
	}
	if mid.ES != 4 || mid.EF != 5 {
		t.Errorf("Mid: ES=4 EF=5, got ES=%d EF=%d", mid.ES, mid.EF)
	}
	if root.ES != 5 || root.EF != 7 {
		t.Errorf("Root: ES=5 EF=7, got ES=%d EF=%d", root.ES, root.EF)
	}
	// Critical path: Right -> Mid -> Root; Left has Slack=2
	if childR.Slack != 0 {
		t.Errorf("Right (critical) Slack=0, got %d", childR.Slack)
	}
	if childL.Slack != 2 {
		t.Errorf("Left Slack=2, got %d", childL.Slack)
	}
}
