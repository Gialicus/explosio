package clone

import (
	"explosio/lib/domain"
	"testing"
)

func TestCloneActivity_Nil(t *testing.T) {
	got := CloneActivity(nil)
	if got != nil {
		t.Errorf("CloneActivity(nil) want nil, got %v", got)
	}
}

func TestCloneActivity_SingleNode(t *testing.T) {
	root := domain.BuildActivityForTest("A", "Root", 5, nil)
	cl := CloneActivity(root)
	if cl == nil {
		t.Fatal("CloneActivity: want non-nil clone")
	}
	if cl.ID != root.ID || cl.Name != root.Name || cl.Duration != root.Duration {
		t.Errorf("clone: ID=%q Name=%q Duration=%d, want same as root", cl.ID, cl.Name, cl.Duration)
	}
	if cl == root {
		t.Error("clone must be a different pointer from root")
	}
}

func TestCloneActivity_TwoLevels(t *testing.T) {
	child := domain.BuildActivityForTest("B", "Child", 3, nil)
	root := domain.BuildActivityForTest("A", "Root", 2, []*domain.Activity{child})
	cl := CloneActivity(root)
	if cl == nil || cl.SubActivities == nil || len(cl.SubActivities) != 1 {
		t.Fatalf("clone: want root with one sub, got %v", cl)
	}
	sub := cl.SubActivities[0]
	if sub.ID != "B" || sub.Name != "Child" || sub.Duration != 3 {
		t.Errorf("clone child: ID=%q Name=%q Duration=%d", sub.ID, sub.Name, sub.Duration)
	}
	if sub == child {
		t.Error("clone child must be a different pointer from original child")
	}
}

func TestCloneActivity_Independent(t *testing.T) {
	child := domain.BuildActivityForTest("B", "Child", 3, nil)
	root := domain.BuildActivityForTest("A", "Root", 2, []*domain.Activity{child})
	cl := CloneActivity(root)
	cl.Duration = 99
	cl.SubActivities[0].Duration = 77
	if root.Duration != 2 {
		t.Errorf("modifying clone must not change original root.Duration: want 2, got %d", root.Duration)
	}
	if child.Duration != 3 {
		t.Errorf("modifying clone must not change original child.Duration: want 3, got %d", child.Duration)
	}
}
