package lib

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrintMermaidTo_NilRoot(t *testing.T) {
	var buf bytes.Buffer
	PrintMermaidTo(&buf, nil)
	out := buf.String()
	if out != "" {
		t.Errorf("PrintMermaidTo(nil): want empty output, got %q", out)
	}
}

func TestPrintMermaidTo_SingleNode(t *testing.T) {
	root := buildActivity("A", "Root", 5, nil)
	var buf bytes.Buffer
	PrintMermaidTo(&buf, root)
	out := buf.String()
	if !strings.Contains(out, "flowchart TB") {
		t.Errorf("output must contain flowchart TB: %q", out)
	}
	if !strings.Contains(out, "A") {
		t.Errorf("output must contain node ID A: %q", out)
	}
	if !strings.Contains(out, "Root") {
		t.Errorf("output must contain node name: %q", out)
	}
	if !strings.Contains(out, "5 min") {
		t.Errorf("output must contain duration: %q", out)
	}
}

func TestPrintMermaidTo_Tree(t *testing.T) {
	child1 := buildActivity("C1", "Child1", 2, nil)
	child2 := buildActivity("C2", "Child2", 3, nil)
	root := buildActivity("R", "Root", 1, []*Activity{child1, child2})
	var buf bytes.Buffer
	PrintMermaidTo(&buf, root)
	out := buf.String()
	if !strings.Contains(out, "flowchart TB") {
		t.Errorf("output must contain flowchart TB: %q", out)
	}
	if !strings.Contains(out, "R") || !strings.Contains(out, "C1") || !strings.Contains(out, "C2") {
		t.Errorf("output must contain all node IDs R, C1, C2: %q", out)
	}
	if !strings.Contains(out, "-->") {
		t.Errorf("output must contain edges: %q", out)
	}
	if !strings.Contains(out, "R --> C1") {
		t.Errorf("output must contain edge R --> C1: %q", out)
	}
	if !strings.Contains(out, "R --> C2") {
		t.Errorf("output must contain edge R --> C2: %q", out)
	}
}

func TestPrintMermaidTo_CriticalPath(t *testing.T) {
	child := buildActivity("C", "Child", 3, nil)
	root := buildActivity("R", "Root", 2, []*Activity{child})
	engine := &AnalysisEngine{}
	engine.ComputeCPM(root)
	var buf bytes.Buffer
	PrintMermaidTo(&buf, root)
	out := buf.String()
	if !strings.Contains(out, "style") {
		t.Errorf("output must contain style for critical path: %q", out)
	}
	if !strings.Contains(out, "fill:#f96") {
		t.Errorf("output must contain fill:#f96 for critical nodes: %q", out)
	}
	// Both root and child are on critical path (Slack == 0)
	if !strings.Contains(out, "style R ") && !strings.Contains(out, "style C ") {
		t.Errorf("output must style at least one critical node (R or C): %q", out)
	}
}

func TestPrintMermaidTo_WithHuman(t *testing.T) {
	root := buildActivity("A", "Root", 5, nil)
	root.Humans = []HumanResource{{Role: "Cameriere", CostPerH: 15, Quantity: 1}}
	var buf bytes.Buffer
	PrintMermaidTo(&buf, root)
	out := buf.String()
	if !strings.Contains(out, "_H_0") {
		t.Errorf("output must contain human node ID suffix _H_0: %q", out)
	}
	if !strings.Contains(out, "([\"") && !strings.Contains(out, "([") {
		t.Errorf("output must contain stadium shape for human: %q", out)
	}
	if !strings.Contains(out, "--> ") && !strings.Contains(out, "A -->") {
		t.Errorf("output must contain edge from activity to resource: %q", out)
	}
	if !strings.Contains(out, "Cameriere") {
		t.Errorf("output must contain role name: %q", out)
	}
	if !strings.Contains(out, "15") {
		t.Errorf("output must contain cost/h: %q", out)
	}
}

func TestPrintMermaidTo_WithMaterial(t *testing.T) {
	root := buildActivity("B", "Root", 1, nil)
	root.Materials = []MaterialResource{{Name: "Miscela", UnitCost: 0.02, Quantity: 18}}
	var buf bytes.Buffer
	PrintMermaidTo(&buf, root)
	out := buf.String()
	if !strings.Contains(out, "_M_0") {
		t.Errorf("output must contain material node ID suffix _M_0: %q", out)
	}
	if !strings.Contains(out, "[/\"") || !strings.Contains(out, "\"/]") {
		t.Errorf("output must contain parallelogram shape for material: %q", out)
	}
	if !strings.Contains(out, "B -->") {
		t.Errorf("output must contain edge from activity B to resource: %q", out)
	}
	if !strings.Contains(out, "Miscela") {
		t.Errorf("output must contain material name: %q", out)
	}
}

func TestPrintMermaidTo_WithAsset(t *testing.T) {
	root := buildActivity("C", "Root", 1, nil)
	root.Assets = []Asset{{Name: "Macchina", CostPerUse: 0.5, Quantity: 2}}
	var buf bytes.Buffer
	PrintMermaidTo(&buf, root)
	out := buf.String()
	if !strings.Contains(out, "_A_0") {
		t.Errorf("output must contain asset node ID suffix _A_0: %q", out)
	}
	if !strings.Contains(out, "[( \"") && !strings.Contains(out, "(\"") {
		t.Errorf("output must contain cylinder shape for asset: %q", out)
	}
	if !strings.Contains(out, "C -->") {
		t.Errorf("output must contain edge from activity C to resource: %q", out)
	}
	if !strings.Contains(out, "Macchina") {
		t.Errorf("output must contain asset name: %q", out)
	}
}

func TestWriteMermaidToFile_NilRoot(t *testing.T) {
	path := filepath.Join(os.TempDir(), "explosio_nil_test.mmd")
	err := WriteMermaidToFile(nil, path)
	if err != nil {
		t.Errorf("WriteMermaidToFile(nil): want nil error, got %v", err)
	}
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
		t.Error("WriteMermaidToFile(nil) must not create file")
	}
}

func TestWriteMermaidToFile_Success(t *testing.T) {
	root := buildActivity("A", "Root", 5, nil)
	path := filepath.Join(os.TempDir(), "explosio_test.mmd")
	defer os.Remove(path)
	err := WriteMermaidToFile(root, path)
	if err != nil {
		t.Fatalf("WriteMermaidToFile: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	out := string(data)
	if !strings.Contains(out, "flowchart TB") {
		t.Errorf("file must contain flowchart TB: %q", out)
	}
	if !strings.Contains(out, "A") {
		t.Errorf("file must contain node ID A: %q", out)
	}
}
