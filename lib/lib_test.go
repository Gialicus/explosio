package lib

import "testing"

// TestDSLAndEngine_Integration verifica che la facade integri correttamente domain e engine.
func TestDSLAndEngine_Integration(t *testing.T) {
	p := NewProject()
	eng := &AnalysisEngine{}
	childL := p.Node("Left", "left", 2)
	childR := p.Node("Right", "right", 4)
	mid := p.Node("Mid", "mid", 1).DependsOn(childL, childR)
	root := p.Start("Root", "root", 2).DependsOn(mid)
	eng.ComputeCPM(root)
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
	if childR.Slack != 0 {
		t.Errorf("Right (critical) Slack=0, got %d", childR.Slack)
	}
	if childL.Slack != 2 {
		t.Errorf("Left Slack=2, got %d", childL.Slack)
	}
}
