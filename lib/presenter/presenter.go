package presenter

import (
	"fmt"
	"io"
	"os"

	"explosio/lib/domain"
	"explosio/lib/resources"
)

// PrintReportTo scrive il report gerarchico su w.
func PrintReportTo(w io.Writer, a *domain.Activity, level int, isLast bool, prefix string) {
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
	resources.ForEachResource(a, func(r domain.Resource) {
		switch x := r.(type) {
		case domain.HumanResource:
			fmt.Fprintf(w, "%s    • Human: %.0f %s (Cost/h: €%.2f)\n", newPrefix, x.Quantity, x.Role, x.CostPerH)
		case domain.MaterialResource:
			fmt.Fprintf(w, "%s    • Material: %.1f %s (Tot: €%.2f)\n", newPrefix, x.Quantity, x.Name, x.UnitCost*x.Quantity)
		case domain.Asset:
			fmt.Fprintf(w, "%s    • Asset: %.1f %s (Tot: €%.2f)\n", newPrefix, x.Quantity, x.Name, x.CostPerUse*x.Quantity)
		}
	})
	for i, sub := range a.SubActivities {
		PrintReportTo(w, sub, level+1, i == len(a.SubActivities)-1, newPrefix)
	}
}

// PrintReport scrive il report su os.Stdout.
func PrintReport(a *domain.Activity, level int, isLast bool, prefix string) {
	PrintReportTo(os.Stdout, a, level, isLast, prefix)
}
