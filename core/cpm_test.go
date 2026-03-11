package core

import (
	"explosio/core/unit"
	"testing"
)

func TestCalculateCriticalPath_WithDependsOn(t *testing.T) {
	// A (1d) -> B (2d) -> C (1d). D (1d) depends on A. So D runs parallel to B,C.
	// Critical path: A -> B -> C = 4 days. D has slack.
	a := NewActivity("A", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	b := NewActivity("B", "", *unit.NewDuration(2, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	c := NewActivity("C", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	d := NewActivity("D", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))

	a.AddActivity(b)
	b.AddActivity(c)
	a.AddActivity(d)
	d.AddDependsOn(a)

	path := a.CalculateCriticalPath()
	slackMap := a.CalculateSlack()

	// Critical path should include A, B, C (the longest chain)
	pathNames := make(map[string]bool)
	for _, p := range path {
		pathNames[p.Name] = true
	}
	if !pathNames["A"] || !pathNames["B"] || !pathNames["C"] {
		t.Errorf("Critical path should include A, B, C, got %v", pathNames)
	}

	// D should have positive slack (parallel to B-C path)
	var info SlackInfo
	var found bool
	for act, inf := range slackMap {
		if act.Name == "D" {
			info = inf
			found = true
			break
		}
	}
	if !found {
		t.Fatal("D not found in slack map")
	}
	if info.Slack <= 0 {
		t.Errorf("D should have positive slack, got %.2f", info.Slack)
	}
}
