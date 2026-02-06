package tree

import (
	"explosio/lib/domain"
	"testing"
)

func TestCountActivities_Nil(t *testing.T) {
	if n := CountActivities(nil); n != 0 {
		t.Errorf("CountActivities(nil) = %d, want 0", n)
	}
}

func TestCountActivities_Single(t *testing.T) {
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	if n := CountActivities(root); n != 1 {
		t.Errorf("CountActivities(single) = %d, want 1", n)
	}
}

func TestCountActivities_Tree(t *testing.T) {
	c1 := domain.BuildActivityForTest("C1", "C1", 1, nil)
	c2 := domain.BuildActivityForTest("C2", "C2", 1, nil)
	root := domain.BuildActivityForTest("R", "R", 1, []*domain.Activity{c1, c2})
	if n := CountActivities(root); n != 3 {
		t.Errorf("CountActivities(tree) = %d, want 3", n)
	}
}

func TestWalk_Nil(t *testing.T) {
	Walk(nil, func(*domain.Activity) { t.Error("callback should not be called") })
}

func TestWalk_Order(t *testing.T) {
	c := domain.BuildActivityForTest("C", "C", 1, nil)
	root := domain.BuildActivityForTest("R", "R", 1, []*domain.Activity{c})
	var ids []string
	Walk(root, func(a *domain.Activity) { ids = append(ids, a.ID) })
	if len(ids) != 2 || ids[0] != "R" || ids[1] != "C" {
		t.Errorf("Walk pre-order: got ids %v, want [R C]", ids)
	}
}
