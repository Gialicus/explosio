package core

import (
	"explosio/core/unit"
	"testing"
)

func TestActivity_IsMilestone(t *testing.T) {
	milestone := NewActivity("M1", "Checkpoint", *unit.NewDuration(0, unit.DurationUnitDay), *unit.NewPrice(0, "EUR"))
	if !milestone.IsMilestone() {
		t.Error("Zero-duration activity should be milestone")
	}

	normal := NewActivity("A", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(10, "EUR"))
	if normal.IsMilestone() {
		t.Error("Non-zero duration activity should not be milestone")
	}
}
