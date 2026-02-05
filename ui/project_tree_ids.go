package ui

import (
	"explosio/lib"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/widget"
)

// BuildActivityMap costruisce la mappa id -> *Activity visitando ricorsivamente project.Root.
func BuildActivityMap(project *lib.Project) map[string]*lib.Activity {
	m := make(map[string]*lib.Activity)
	if project.Root != nil {
		buildMapRec(project.Root, m)
	}
	return m
}

func buildMapRec(a *lib.Activity, m map[string]*lib.Activity) {
	if a == nil {
		return
	}
	m[a.ID] = a
	for _, sub := range a.SubActivities {
		buildMapRec(sub, m)
	}
}

func isResourceNodeID(id widget.TreeNodeID) bool {
	return strings.Contains(string(id), "|")
}

func parseResourceNodeID(id widget.TreeNodeID) (activityID, kind string, index int, ok bool) {
	parts := strings.Split(string(id), "|")
	if len(parts) != 3 {
		return "", "", 0, false
	}
	idx, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", 0, false
	}
	return parts[0], parts[1], idx, true
}
