package core

import (
	"bytes"
	"explosio/core/unit"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	root := NewActivity("Root", "Test root", *unit.NewDuration(5, unit.DurationUnitDay), *unit.NewPrice(100, "EUR"))
	child := NewActivity("Child", "Test child", *unit.NewDuration(2, unit.DurationUnitDay), *unit.NewPrice(50, "EUR"))
	root.AddActivity(child)
	criticalPath := root.CalculateCriticalPath()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe: %v", err)
	}
	os.Stdout = w

	PrettyPrint([]*Activity{root}, criticalPath)

	os.Stdout = old
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "Root") {
		t.Error("PrettyPrint output should contain Root")
	}
	if !strings.Contains(out, "Child") {
		t.Error("PrettyPrint output should contain Child")
	}
	if !strings.Contains(out, "100.00 EUR") {
		t.Error("PrettyPrint output should contain price")
	}
	if !strings.Contains(out, "Legend:") {
		t.Error("PrettyPrint output should contain Legend")
	}
}

func TestPrettyPrint_NilCriticalPath(t *testing.T) {
	root := NewActivity("Solo", "", *unit.NewDuration(1, unit.DurationUnitDay), *unit.NewPrice(10, "EUR"))

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe: %v", err)
	}
	os.Stdout = w

	PrettyPrint([]*Activity{root}, nil)

	os.Stdout = old
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "Solo") {
		t.Error("PrettyPrint with nil criticalPath should still print activity")
	}
}
