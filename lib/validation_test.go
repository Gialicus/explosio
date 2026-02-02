package lib

import (
	"strings"
	"testing"
)

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

func TestActivity_ValidateBasic(t *testing.T) {
	root := buildActivity("A", "A", 10, nil)
	err := root.ValidateBasic()
	if err != nil {
		t.Errorf("ValidateBasic valid: want nil, got %v", err)
	}
}

func TestActivity_ValidateBasic_Nil(t *testing.T) {
	var root *Activity
	err := root.ValidateBasic()
	if err == nil {
		t.Fatal("ValidateBasic nil: want error")
	}
}

func TestActivity_ValidateBasic_NegativeDuration(t *testing.T) {
	root := buildActivity("A", "A", -5, nil)
	err := root.ValidateBasic()
	if err == nil {
		t.Fatal("ValidateBasic negative duration: want error")
	}
}
