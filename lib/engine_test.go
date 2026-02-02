package lib

import (
	"math"
	"strings"
	"testing"
)

const floatEpsilon = 1e-9

func assertFloatEqual(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > floatEpsilon {
		t.Errorf("got %v, want %v", got, want)
	}
}

// buildActivity crea un'attività minimale per test (senza DSL).
func buildActivity(id, name string, duration int, subs []*Activity) *Activity {
	a := &Activity{
		ID:          id,
		Name:        name,
		Description: name,
		Duration:    duration,
		MinDuration: duration,
		SubActivities: subs,
	}
	for _, s := range subs {
		s.Next = append(s.Next, a.ID)
	}
	return a
}

func TestComputeCPM_NilRoot(t *testing.T) {
	engine := &AnalysisEngine{}
	engine.ComputeCPM(nil)
	// non deve andare in panic
}

func TestComputeCPM_SingleNode(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 5, nil)
	engine.ComputeCPM(root)
	if root.ES != 0 || root.EF != 5 || root.LS != 0 || root.LF != 5 || root.Slack != 0 {
		t.Errorf("single node: ES=0 EF=5 LS=0 LF=5 Slack=0, got ES=%d EF=%d LS=%d LF=%d Slack=%d",
			root.ES, root.EF, root.LS, root.LF, root.Slack)
	}
}

func TestComputeCPM_TwoLevels(t *testing.T) {
	engine := &AnalysisEngine{}
	child := buildActivity("B", "B", 3, nil)
	root := buildActivity("A", "A", 2, []*Activity{child})
	engine.ComputeCPM(root)
	// figlio: ES=0, EF=3, LS=0, LF=3, Slack=0
	if child.ES != 0 || child.EF != 3 || child.LS != 0 || child.LF != 3 || child.Slack != 0 {
		t.Errorf("child: ES=0 EF=3 LS=0 LF=3 Slack=0, got ES=%d EF=%d LS=%d LF=%d Slack=%d",
			child.ES, child.EF, child.LS, child.LF, child.Slack)
	}
	// radice: ES=3, EF=5, LS=3, LF=5, Slack=0
	if root.ES != 3 || root.EF != 5 || root.LS != 3 || root.LF != 5 || root.Slack != 0 {
		t.Errorf("root: ES=3 EF=5 LS=3 LF=5 Slack=0, got ES=%d EF=%d LS=%d LF=%d Slack=%d",
			root.ES, root.EF, root.LS, root.LF, root.Slack)
	}
}

func TestComputeCPM_TwoChildrenParallel(t *testing.T) {
	engine := &AnalysisEngine{}
	short := buildActivity("S", "Short", 2, nil)
	long := buildActivity("L", "Long", 4, nil)
	root := buildActivity("R", "Root", 1, []*Activity{short, long})
	engine.ComputeCPM(root)
	// short: ES=0 EF=2, long: ES=0 EF=4; root ES=4 EF=5
	if short.ES != 0 || short.EF != 2 {
		t.Errorf("short: ES=0 EF=2, got ES=%d EF=%d", short.ES, short.EF)
	}
	if long.ES != 0 || long.EF != 4 {
		t.Errorf("long: ES=0 EF=4, got ES=%d EF=%d", long.ES, long.EF)
	}
	if root.ES != 4 || root.EF != 5 {
		t.Errorf("root: ES=4 EF=5, got ES=%d EF=%d", root.ES, root.EF)
	}
	// backward: root LF=5 LS=4; long LF=4 LS=0 Slack=0; short LF=4 LS=2 Slack=2
	if long.Slack != 0 {
		t.Errorf("long (critical): Slack=0, got %d", long.Slack)
	}
	if short.Slack != 2 {
		t.Errorf("short: Slack=2, got %d", short.Slack)
	}
}

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

func TestValidate_Valid(t *testing.T) {
	engine := &AnalysisEngine{}
	child := buildActivity("B", "B", 2, nil)
	root := buildActivity("A", "A", 3, []*Activity{child})
	err := engine.Validate(root)
	if err != nil {
		t.Errorf("Validate valid tree: want nil, got %v", err)
	}
}

func TestValidate_NilRoot(t *testing.T) {
	engine := &AnalysisEngine{}
	err := engine.Validate(nil)
	if err == nil {
		t.Fatal("Validate nil root: want error")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error should mention nil: %v", err)
	}
}

func TestValidate_NegativeDuration(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", -1, nil)
	err := engine.Validate(root)
	if err == nil {
		t.Fatal("Validate negative duration: want error")
	}
	if !strings.Contains(err.Error(), "duration negative") {
		t.Errorf("error should mention duration: %v", err)
	}
}

func TestValidate_MinDurationGreaterThanDuration(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 5, nil)
	root.MinDuration = 10
	err := engine.Validate(root)
	if err == nil {
		t.Fatal("Validate MinDuration > Duration: want error")
	}
	if !strings.Contains(err.Error(), "min duration greater") {
		t.Errorf("error should mention min duration: %v", err)
	}
}

func TestGetTotalDuration(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 5, nil)
	engine.ComputeCPM(root)
	got := engine.GetTotalDuration(root)
	if got != 5 {
		t.Errorf("GetTotalDuration want 5, got %d", got)
	}
	gotNil := engine.GetTotalDuration(nil)
	if gotNil != 0 {
		t.Errorf("GetTotalDuration(nil) want 0, got %d", gotNil)
	}
}

func TestGetCriticalPath(t *testing.T) {
	engine := &AnalysisEngine{}
	short := buildActivity("S", "Short", 2, nil)
	long := buildActivity("L", "Long", 4, nil)
	root := buildActivity("R", "Root", 1, []*Activity{short, long})
	engine.ComputeCPM(root)
	path := engine.GetCriticalPath(root)
	// Only long and root have Slack==0
	if len(path) != 2 {
		t.Errorf("GetCriticalPath want 2 activities, got %d", len(path))
	}
	ids := make(map[string]bool)
	for _, a := range path {
		ids[a.ID] = true
	}
	if !ids["R"] || !ids["L"] {
		t.Errorf("CriticalPath should contain R and L, got %v", ids)
	}
	if ids["S"] {
		t.Error("Short (Slack=2) should not be in CriticalPath")
	}
}

func TestGetCPMSummary(t *testing.T) {
	engine := &AnalysisEngine{}
	child := buildActivity("B", "B", 3, nil)
	root := buildActivity("A", "A", 2, []*Activity{child})
	engine.ComputeCPM(root)
	sum := engine.GetCPMSummary(root)
	if sum.TotalDuration != 5 {
		t.Errorf("TotalDuration want 5, got %d", sum.TotalDuration)
	}
	if sum.ActivityCount != 2 {
		t.Errorf("ActivityCount want 2, got %d", sum.ActivityCount)
	}
	if len(sum.CriticalPath) != 2 {
		t.Errorf("CriticalPath len want 2, got %d", len(sum.CriticalPath))
	}
}

func TestWalk(t *testing.T) {
	engine := &AnalysisEngine{}
	child := buildActivity("B", "B", 1, nil)
	root := buildActivity("A", "A", 1, []*Activity{child})
	var ids []string
	engine.Walk(root, func(a *Activity) { ids = append(ids, a.ID) })
	if len(ids) != 2 {
		t.Errorf("Walk: want 2 nodes, got %d", len(ids))
	}
	if ids[0] != "A" || ids[1] != "B" {
		t.Errorf("Walk pre-order: want [A B], got %v", ids)
	}
	engine.Walk(nil, func(*Activity) { t.Error("Walk(nil) should not call f") })
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

func TestActivitiesByES(t *testing.T) {
	engine := &AnalysisEngine{}
	c1 := buildActivity("C1", "C1", 2, nil)
	c2 := buildActivity("C2", "C2", 4, nil)
	root := buildActivity("R", "R", 1, []*Activity{c1, c2})
	engine.ComputeCPM(root)
	list := engine.ActivitiesByES(root)
	if len(list) != 3 {
		t.Fatalf("ActivitiesByES want 3, got %d", len(list))
	}
	for i := 1; i < len(list); i++ {
		if list[i].ES < list[i-1].ES {
			t.Errorf("ActivitiesByES not sorted: %d < %d at %d", list[i].ES, list[i-1].ES, i)
		}
	}
}
