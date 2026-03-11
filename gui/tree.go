// Package gui provides the Fyne-based GUI for Explosio activity tree editing.
package gui

import (
	"explosio/core"
	"explosio/core/unit"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const rootUID = ""

// pathToActivity navigates the activity tree using path indices.
// Path "" = root, "0" = first child, "0-1" = second child of first, etc.
func pathToActivity(root *core.Activity, path string) *core.Activity {
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

// activityToPath returns the path string for the given activity within the tree.
// Returns empty string if activity is root, or path like "0", "0-1", etc.
func activityToPath(root *core.Activity, target *core.Activity) string {
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

// activityParent returns the parent activity and index of the given activity.
// For root, returns (nil, -1).
func activityParent(root *core.Activity, target *core.Activity) (*core.Activity, int) {
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

// NewActivityTree creates a Tree widget bound to the given root activity.
// onSelect is called when a node is selected, with the selected activity (nil if root for "").
func NewActivityTree(root *core.Activity, onSelect func(*core.Activity)) *widget.Tree {
	if root == nil {
		root = core.NewActivity("Progetto", "", unit.Duration{Value: 0, Unit: unit.DurationUnitDay}, unit.Price{Value: 0, Currency: "EUR"})
	}

	childUIDs := func(uid widget.TreeNodeID) []widget.TreeNodeID {
		a := pathToActivity(root, uid)
		if a == nil || len(a.Activities) == 0 {
			return nil
		}
		prefix := uid
		if prefix != "" {
			prefix += "-"
		}
		ids := make([]widget.TreeNodeID, len(a.Activities))
		for i := range a.Activities {
			ids[i] = prefix + strconv.Itoa(i)
		}
		return ids
	}

	isBranch := func(uid widget.TreeNodeID) bool {
		a := pathToActivity(root, uid)
		return a != nil && len(a.Activities) > 0
	}

	createNode := func(branch bool) fyne.CanvasObject {
		return widget.NewLabel("")
	}

	updateNode := func(uid widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
		a := pathToActivity(root, uid)
		if a == nil {
			return
		}
		label := obj.(*widget.Label)
		dur := fmt.Sprintf("%.0f %s", a.Duration.Value, a.Duration.Unit)
		price := fmt.Sprintf("%.0f %s", a.CalculatePrice(), a.Price.Currency)
		label.SetText(fmt.Sprintf("%s [%s, %s]", a.Name, dur, price))
	}

	tree := widget.NewTree(childUIDs, isBranch, createNode, updateNode)
	tree.Root = rootUID

	tree.OnSelected = func(uid widget.TreeNodeID) {
		if onSelect != nil {
			onSelect(pathToActivity(root, uid))
		}
	}

	return tree
}
