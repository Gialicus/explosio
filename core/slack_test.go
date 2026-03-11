package core

import (
	"explosio/core/unit"
	"testing"
)

func TestCalculateSlack_NonCriticalHasSlack(t *testing.T) {
	root := NewActivity("Root", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	long := NewActivity("Long", "", *unit.NewDuration(3, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	short := NewActivity("Short", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	root.AddActivity(long)
	root.AddActivity(short)

	slackMap := root.CalculateSlack()
	path := root.CalculateCriticalPath()
	criticalSet := make(map[*Activity]bool)
	for _, a := range path {
		criticalSet[a] = true
	}

	if criticalSet[short] {
		t.Fatal("Short should not be on critical path")
	}
	info, ok := slackMap[short]
	if !ok {
		t.Fatal("Short should be in slack map")
	}
	if info.Slack <= 0 {
		t.Errorf("Short (non-critical) should have positive slack, got %.2f", info.Slack)
	}
	// Short has 1 day, Long has 3 days. Slack = 2 days = 48 hours
	if info.Slack < 40 {
		t.Errorf("Short slack should be ~48h, got %.2f", info.Slack)
	}
}
