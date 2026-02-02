package lib

import "testing"

func TestGetMaxCrashPotential_Single(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 10, nil)
	root.MinDuration = 6
	root.CrashCostStep = 2.5
	timeSaved, extraCost := engine.GetMaxCrashPotential(root)
	if timeSaved != 4 {
		t.Errorf("timeSaved want 4, got %d", timeSaved)
	}
	assertFloatEqual(t, extraCost, 10) // 4 * 2.5
}

func TestGetMaxCrashPotential_Tree(t *testing.T) {
	engine := &AnalysisEngine{}
	c1 := buildActivity("C1", "C1", 5, nil)
	c1.MinDuration = 3
	c1.CrashCostStep = 1
	c2 := buildActivity("C2", "C2", 4, nil)
	c2.MinDuration = 4
	c2.CrashCostStep = 2
	root := buildActivity("R", "R", 2, []*Activity{c1, c2})
	root.MinDuration = 1
	root.CrashCostStep = 0.5
	timeSaved, extraCost := engine.GetMaxCrashPotential(root)
	// root: 1 min, 0.5; c1: 2 min, 2; c2: 0; total time 3, cost 2.5
	if timeSaved != 3 {
		t.Errorf("timeSaved want 3, got %d", timeSaved)
	}
	assertFloatEqual(t, extraCost, 2.5)
}

func TestCrashWithBudget(t *testing.T) {
	engine := &AnalysisEngine{}
	// Single critical child: 4 min, crashable to 1 min (3 min save, cost 3).
	child := buildActivity("C", "C", 4, nil)
	child.MinDuration = 1
	child.CrashCostStep = 1
	root := buildActivity("R", "R", 1, []*Activity{child})
	engine.ComputeCPM(root)
	timeSaved, actualCost := engine.CrashWithBudget(root, 5)
	if actualCost > 5 {
		t.Errorf("actualCost must be <= 5, got %v", actualCost)
	}
	if timeSaved != 3 {
		t.Errorf("timeSaved want 3, got %d", timeSaved)
	}
	assertFloatEqual(t, actualCost, 3)
	// Budget too low: can't afford any crash.
	timeSaved2, actualCost2 := engine.CrashWithBudget(root, 1)
	if timeSaved2 != 0 || actualCost2 != 0 {
		t.Errorf("with budget 1: want (0,0), got (%d,%v)", timeSaved2, actualCost2)
	}
}

func TestCrashToSaveTime(t *testing.T) {
	engine := &AnalysisEngine{}
	c1 := buildActivity("C1", "C1", 5, nil)
	c1.MinDuration = 2
	c1.CrashCostStep = 1 // 3 min, 3 cost
	root := buildActivity("R", "R", 1, []*Activity{c1})
	engine.ComputeCPM(root)
	extraCost, achieved := engine.CrashToSaveTime(root, 2)
	if !achieved {
		t.Error("CrashToSaveTime(2) want achieved true")
	}
	assertFloatEqual(t, extraCost, 3)
}
