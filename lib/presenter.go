package lib

import (
	"fmt"
	"io"
	"os"
)

// PrintReportTo scrive il report gerarchico su w.
func PrintReportTo(w io.Writer, a *Activity, level int, isLast bool, prefix string) {
	marker := "├── "
	if level == 0 {
		marker = ""
	} else if isLast {
		marker = "└── "
	}

	tag := "[ ]"
	if a.Slack == 0 {
		tag = "[CRITICAL]"
	}

	fmt.Fprintf(w, "%s%s%s %s: %s (%d min)\n", prefix, marker, tag, a.Name, a.Description, a.Duration)

	newPrefix := prefix
	if level > 0 {
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
	}

	for _, h := range a.Humans {
		fmt.Fprintf(w, "%s    • Umano: %.0f %s (Cost/h: €%.2f)\n", newPrefix, h.Quantity, h.Role, h.CostPerH)
	}
	for _, m := range a.Materials {
		fmt.Fprintf(w, "%s    • Mat: %.1f %s (Tot: €%.2f)\n", newPrefix, m.Quantity, m.Name, m.UnitCost*m.Quantity)
	}
	for _, as := range a.Assets {
		fmt.Fprintf(w, "%s    • Asset: %.1f %s (Tot: €%.2f)\n", newPrefix, as.Quantity, as.Name, as.CostPerUse*as.Quantity)
	}

	for i, sub := range a.SubActivities {
		PrintReportTo(w, sub, level+1, i == len(a.SubActivities)-1, newPrefix)
	}
}

// PrintReport scrive il report su os.Stdout.
func PrintReport(a *Activity, level int, isLast bool, prefix string) {
	PrintReportTo(os.Stdout, a, level, isLast, prefix)
}
