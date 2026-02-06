package presenter

import (
	"bytes"
	"explosio/lib/domain"
	"strings"
	"testing"
)

func TestPrintReportTo_SingleActivity(t *testing.T) {
	root := domain.BuildActivityForTest("A", "Desc A", 5, nil)
	root.Slack = 0
	var buf bytes.Buffer
	PrintReportTo(&buf, root, 0, true, "")
	out := buf.String()
	if !strings.Contains(out, "A") {
		t.Errorf("output must contain activity name A: %q", out)
	}
	if !strings.Contains(out, "Desc A") {
		t.Errorf("output must contain description: %q", out)
	}
	if !strings.Contains(out, "5 min") {
		t.Errorf("output must contain duration: %q", out)
	}
	if !strings.Contains(out, "[CRITICAL]") {
		t.Errorf("output must contain [CRITICAL] when Slack==0: %q", out)
	}
}

func TestPrintReportTo_NonCritical(t *testing.T) {
	root := domain.BuildActivityForTest("B", "Desc B", 2, nil)
	root.Slack = 3
	var buf bytes.Buffer
	PrintReportTo(&buf, root, 0, true, "")
	out := buf.String()
	if !strings.Contains(out, "B") {
		t.Errorf("output must contain activity name B: %q", out)
	}
	if strings.Contains(out, "[CRITICAL]") {
		t.Errorf("output must not contain [CRITICAL] when Slack!=0: %q", out)
	}
	if !strings.Contains(out, "[ ]") {
		t.Errorf("output must contain [ ] for non-critical: %q", out)
	}
}

func TestPrintReportTo_WithHumanAndMaterial(t *testing.T) {
	root := domain.BuildActivityForTest("R", "Root", 1, nil)
	root.Humans = []domain.HumanResource{{Role: "Cameriere", CostPerH: 15, Quantity: 1}}
	root.Materials = []domain.MaterialResource{{Name: "Caffè", UnitCost: 0.02, Quantity: 18}}
	var buf bytes.Buffer
	PrintReportTo(&buf, root, 0, true, "")
	out := buf.String()
	if !strings.Contains(out, "Human") {
		t.Errorf("output must contain Human when Humans present: %q", out)
	}
	if !strings.Contains(out, "Cameriere") {
		t.Errorf("output must contain role name: %q", out)
	}
	if !strings.Contains(out, "Caffè") {
		t.Errorf("output must contain material name: %q", out)
	}
}

func TestPrintReportTo_WithAsset(t *testing.T) {
	root := domain.BuildActivityForTest("R", "Root", 1, nil)
	root.Assets = []domain.Asset{{Name: "Macchina", Description: "Espresso", CostPerUse: 0.5, Quantity: 2}}
	var buf bytes.Buffer
	PrintReportTo(&buf, root, 0, true, "")
	out := buf.String()
	if !strings.Contains(out, "Asset") {
		t.Errorf("output must contain Asset when Assets present: %q", out)
	}
	if !strings.Contains(out, "Macchina") {
		t.Errorf("output must contain asset name: %q", out)
	}
}

func TestPrintReportTo_Tree(t *testing.T) {
	child := domain.BuildActivityForTest("C", "Child", 2, nil)
	child.Slack = 0
	root := domain.BuildActivityForTest("R", "Root", 1, []*domain.Activity{child})
	root.Slack = 0
	var buf bytes.Buffer
	PrintReportTo(&buf, root, 0, true, "")
	out := buf.String()
	if !strings.Contains(out, "R") || !strings.Contains(out, "C") {
		t.Errorf("output must contain both root and child names: %q", out)
	}
	if !strings.Contains(out, "└──") && !strings.Contains(out, "├──") {
		t.Errorf("output must contain tree markers for children: %q", out)
	}
}

func TestPrintReport_NoPanic(t *testing.T) {
	root := domain.BuildActivityForTest("X", "X", 1, nil)
	PrintReport(root, 0, true, "")
}
