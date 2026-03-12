package app

import (
	"strconv"
	"strings"

	"explosio/core"
	"explosio/core/unit"
)

// ActivityService handles activity tree operations.
type ActivityService struct{}

// NewActivityService creates a new ActivityService.
func NewActivityService() *ActivityService {
	return &ActivityService{}
}

// PathToActivity navigates the activity tree using path indices.
// Path "" = root, "0" = first child, "0-1" = second child of first, etc.
func (s *ActivityService) PathToActivity(root *core.Activity, path string) *core.Activity {
	if root == nil || path == "" {
		return root
	}
	parts := strings.Split(path, "-")
	a := root
	for _, p := range parts {
		idx, err := strconv.Atoi(p)
		if err != nil || idx < 0 || idx >= len(a.Activities) {
			return nil
		}
		a = a.Activities[idx]
		if a == nil {
			return nil
		}
	}
	return a
}

// ActivityToPath returns the path string for the given activity within the tree.
// Returns empty string if activity is root, or path like "0", "0-1", etc.
func (s *ActivityService) ActivityToPath(root *core.Activity, target *core.Activity) string {
	if root == nil || target == nil || root == target {
		return ""
	}
	var find func(a *core.Activity, prefix string) string
	find = func(a *core.Activity, prefix string) string {
		for i, child := range a.Activities {
			p := prefix + strconv.Itoa(i)
			if prefix != "" {
				p = prefix + "-" + strconv.Itoa(i)
			}
			if child == target {
				return p
			}
			if found := find(child, p); found != "" {
				return found
			}
		}
		return ""
	}
	return find(root, "")
}

// GetParent returns the parent activity and index of the given activity.
// For root, returns (nil, -1).
func (s *ActivityService) GetParent(root *core.Activity, target *core.Activity) (*core.Activity, int) {
	if root == nil || target == nil || root == target {
		return nil, -1
	}
	var find func(a *core.Activity) (*core.Activity, int)
	find = func(a *core.Activity) (*core.Activity, int) {
		for i, child := range a.Activities {
			if child == target {
				return a, i
			}
			if parent, idx := find(child); parent != nil {
				return parent, idx
			}
		}
		return nil, -1
	}
	return find(root)
}

// AddChild creates a new activity and adds it as a child of parent.
func (s *ActivityService) AddChild(parent *core.Activity, name string) *core.Activity {
	if parent == nil {
		return nil
	}
	if name == "" {
		name = "Nuova attività"
	}
	child := core.NewActivity(
		name,
		"",
		unit.Duration{Value: 0, Unit: unit.DurationUnitDay},
		unit.Price{Value: 0, Currency: "EUR"},
	)
	parent.AddActivity(child)
	return child
}

// DeleteActivity removes the target activity from the tree.
// Returns false if target is root (cannot delete root).
func (s *ActivityService) DeleteActivity(root *core.Activity, target *core.Activity) bool {
	if root == nil || target == nil || root == target {
		return false
	}
	parent, idx := s.GetParent(root, target)
	if parent == nil || idx < 0 {
		return false
	}
	parent.Activities = append(parent.Activities[:idx], parent.Activities[idx+1:]...)
	return true
}

// UpdateFields holds the editable fields of an activity for update.
type UpdateFields struct {
	Name        string
	Description string
	Duration    unit.Duration
	Price       unit.Price
}

// UpdateActivity updates the given activity with the provided fields.
func (s *ActivityService) UpdateActivity(a *core.Activity, fields UpdateFields) {
	if a == nil {
		return
	}
	a.Name = fields.Name
	a.Description = fields.Description
	a.Duration = fields.Duration
	a.Price = fields.Price
	if a.Price.Currency == "" {
		a.Price.Currency = "EUR"
	}
}
