package components

import (
	"fmt"
	"strconv"
	"strings"

	"explosio/core"
	"explosio/core/unit"
	"explosio/gui/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const rootUID = ""

// treeUIDToActivityPath converts a tree node UID to the activity path.
// We use a virtual root: tree root "" has one child "0" = root activity.
// Path "0" -> activity "", "0-0" -> "0", "0-1" -> "1", "0-1-0" -> "1-0", etc.
func treeUIDToActivityPath(uid widget.TreeNodeID) string {
	if uid == "" || uid == "0" {
		return ""
	}
	return strings.TrimPrefix(uid, "0-")
}

// NewActivityTree creates a Tree widget bound to the given root activity.
// Uses a virtual root so the root activity (e.g. "Home Renovation") is visible as the first node.
// actSvc is used for path-to-activity navigation; onSelect is called when a node is selected.
func NewActivityTree(root *core.Activity, actSvc *app.ActivityService, onSelect func(*core.Activity)) *widget.Tree {
	if root == nil {
		root = core.NewActivity("Progetto", "", unit.Duration{Value: 0, Unit: unit.DurationUnitDay}, unit.Price{Value: 0, Currency: "EUR"})
	}

	childUIDs := func(uid widget.TreeNodeID) []widget.TreeNodeID {
		if uid == "" {
			// Virtual root: single child is the root activity
			return []widget.TreeNodeID{"0"}
		}
		activityPath := treeUIDToActivityPath(uid)
		a := actSvc.PathToActivity(root, activityPath)
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
		if uid == "" {
			return true
		}
		activityPath := treeUIDToActivityPath(uid)
		a := actSvc.PathToActivity(root, activityPath)
		return a != nil && len(a.Activities) > 0
	}

	createNode := func(branch bool) fyne.CanvasObject {
		return widget.NewLabel("")
	}

	updateNode := func(uid widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
		a := actSvc.PathToActivity(root, treeUIDToActivityPath(uid))
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
			onSelect(actSvc.PathToActivity(root, treeUIDToActivityPath(uid)))
		}
	}

	return tree
}

// TreeSelectActivityPath returns the tree UID for selecting the given activity.
// Used when the tree needs to select an activity by path (e.g. "1-0" -> "0-1-0").
func TreeSelectActivityPath(activityPath string) string {
	if activityPath == "" {
		return "0"
	}
	return "0-" + activityPath
}
