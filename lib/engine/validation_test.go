package engine

import (
	"explosio/lib/domain"
	"strings"
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	eng := &AnalysisEngine{}
	child := domain.BuildActivityForTest("B", "B", 2, nil)
	root := domain.BuildActivityForTest("A", "A", 3, []*domain.Activity{child})
	err := eng.Validate(root)
	if err != nil {
		t.Errorf("Validate valid tree: want nil, got %v", err)
	}
}

func TestValidate_NilRoot(t *testing.T) {
	eng := &AnalysisEngine{}
	err := eng.Validate(nil)
	if err == nil {
		t.Fatal("Validate nil root: want error")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error should mention nil: %v", err)
	}
}

func TestValidate_NegativeDuration(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", -1, nil)
	err := eng.Validate(root)
	if err == nil {
		t.Fatal("Validate negative duration: want error")
	}
	if !strings.Contains(err.Error(), "duration") || !strings.Contains(err.Error(), "negative") {
		t.Errorf("error should mention duration and negative: %v", err)
	}
}

func TestValidate_MinDurationGreaterThanDuration(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", 5, nil)
	root.MinDuration = 10
	err := eng.Validate(root)
	if err == nil {
		t.Fatal("Validate MinDuration > Duration: want error")
	}
	if !strings.Contains(err.Error(), "min duration") || !strings.Contains(err.Error(), "greater") {
		t.Errorf("error should mention min duration and greater: %v", err)
	}
}

func TestActivity_ValidateBasic(t *testing.T) {
	root := domain.BuildActivityForTest("A", "A", 10, nil)
	err := root.ValidateBasic()
	if err != nil {
		t.Errorf("ValidateBasic valid: want nil, got %v", err)
	}
}

func TestActivity_ValidateBasic_Nil(t *testing.T) {
	var root *domain.Activity
	err := root.ValidateBasic()
	if err == nil {
		t.Fatal("ValidateBasic nil: want error")
	}
}

func TestActivity_ValidateBasic_NegativeDuration(t *testing.T) {
	root := domain.BuildActivityForTest("A", "A", -5, nil)
	err := root.ValidateBasic()
	if err == nil {
		t.Fatal("ValidateBasic negative duration: want error")
	}
}
