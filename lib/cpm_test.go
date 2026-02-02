package lib

import "testing"

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

func TestPeriodType_ToMinutes_Invalid(t *testing.T) {
	invalidPeriod := PeriodType("invalid")
	minutes := invalidPeriod.ToMinutes()
	if minutes != 0 {
		t.Errorf("ToMinutes() for invalid period should return 0, got %d", minutes)
	}
}

func TestPeriodType_IsValid(t *testing.T) {
	validPeriods := []PeriodType{PeriodMinute, PeriodHour, PeriodDay, PeriodWeek, PeriodMonth, PeriodYear}
	for _, p := range validPeriods {
		if !p.IsValid() {
			t.Errorf("IsValid() should return true for %s", p)
		}
	}
	invalidPeriod := PeriodType("invalid")
	if invalidPeriod.IsValid() {
		t.Error("IsValid() should return false for invalid period")
	}
}
