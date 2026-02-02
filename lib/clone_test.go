package lib

import "testing"

func TestCloneActivity_Nil(t *testing.T) {
	got := CloneActivity(nil)
	if got != nil {
		t.Errorf("CloneActivity(nil) want nil, got %v", got)
	}
}

func TestCloneActivity_SingleNode(t *testing.T) {
	root := buildActivity("A", "Root", 5, nil)
	clone := CloneActivity(root)
	if clone == nil {
		t.Fatal("CloneActivity: want non-nil clone")
	}
	if clone.ID != root.ID || clone.Name != root.Name || clone.Duration != root.Duration {
		t.Errorf("clone: ID=%q Name=%q Duration=%d, want same as root", clone.ID, clone.Name, clone.Duration)
	}
	if clone == root {
		t.Error("clone must be a different pointer from root")
	}
}

func TestCloneActivity_TwoLevels(t *testing.T) {
	child := buildActivity("B", "Child", 3, nil)
	root := buildActivity("A", "Root", 2, []*Activity{child})
	clone := CloneActivity(root)
	if clone == nil || clone.SubActivities == nil || len(clone.SubActivities) != 1 {
		t.Fatalf("clone: want root with one sub, got %v", clone)
	}
	sub := clone.SubActivities[0]
	if sub.ID != "B" || sub.Name != "Child" || sub.Duration != 3 {
		t.Errorf("clone child: ID=%q Name=%q Duration=%d", sub.ID, sub.Name, sub.Duration)
	}
	if sub == child {
		t.Error("clone child must be a different pointer from original child")
	}
}

func TestCloneActivity_Independent(t *testing.T) {
	child := buildActivity("B", "Child", 3, nil)
	root := buildActivity("A", "Root", 2, []*Activity{child})
	clone := CloneActivity(root)
	clone.Duration = 99
	clone.SubActivities[0].Duration = 77
	if root.Duration != 2 {
		t.Errorf("modifying clone must not change original root.Duration: want 2, got %d", root.Duration)
	}
	if child.Duration != 3 {
		t.Errorf("modifying clone must not change original child.Duration: want 3, got %d", child.Duration)
	}
}
