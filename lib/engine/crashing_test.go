package engine

import (
	"explosio/lib/domain"
	"testing"
)

func TestGetMaxCrashPotential_Single(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", 10, nil)
	root.MinDuration = 6
	root.CrashCostStep = 2.5
	timeSaved, extraCost := eng.GetMaxCrashPotential(root)
	if timeSaved != 4 {
		t.Errorf("timeSaved want 4, got %d", timeSaved)
	}
	assertFloatEqual(t, extraCost, 10)
}

func TestGetMaxCrashPotential_Tree(t *testing.T) {
	eng := &AnalysisEngine{}
	c1 := domain.BuildActivityForTest("C1", "C1", 5, nil)
	c1.MinDuration = 3
	c1.CrashCostStep = 1
	c2 := domain.BuildActivityForTest("C2", "C2", 4, nil)
	c2.MinDuration = 4
	c2.CrashCostStep = 2
	root := domain.BuildActivityForTest("R", "R", 2, []*domain.Activity{c1, c2})
	root.MinDuration = 1
	root.CrashCostStep = 0.5
	timeSaved, extraCost := eng.GetMaxCrashPotential(root)
	if timeSaved != 3 {
		t.Errorf("timeSaved want 3, got %d", timeSaved)
	}
	assertFloatEqual(t, extraCost, 2.5)
}

func TestCrashWithBudget(t *testing.T) {
	eng := &AnalysisEngine{}
	child := domain.BuildActivityForTest("C", "C", 4, nil)
	child.MinDuration = 1
	child.CrashCostStep = 1
	root := domain.BuildActivityForTest("R", "R", 1, []*domain.Activity{child})
	eng.ComputeCPM(root)
	timeSaved, actualCost := eng.CrashWithBudget(root, 5)
	if actualCost > 5 {
		t.Errorf("actualCost must be <= 5, got %v", actualCost)
	}
	if timeSaved != 3 {
		t.Errorf("timeSaved want 3, got %d", timeSaved)
	}
	assertFloatEqual(t, actualCost, 3)
	timeSaved2, actualCost2 := eng.CrashWithBudget(root, 1)
	if timeSaved2 != 0 || actualCost2 != 0 {
		t.Errorf("with budget 1: want (0,0), got (%d,%v)", timeSaved2, actualCost2)
	}
}

func TestCrashToSaveTime(t *testing.T) {
	eng := &AnalysisEngine{}
	c1 := domain.BuildActivityForTest("C1", "C1", 5, nil)
	c1.MinDuration = 2
	c1.CrashCostStep = 1
	root := domain.BuildActivityForTest("R", "R", 1, []*domain.Activity{c1})
	eng.ComputeCPM(root)
	extraCost, achieved := eng.CrashToSaveTime(root, 2)
	if !achieved {
		t.Error("CrashToSaveTime(2) want achieved true")
	}
	assertFloatEqual(t, extraCost, 3)
}
