package presenter

import (
	"fmt"
	"io"
	"os"

	"explosio/lib/domain"
	"explosio/lib/engine"
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
		supplierInfo := ""
		if s := r.GetSupplier(); s != nil {
			supplierInfo = fmt.Sprintf(" (validato da %s: max %.1f/%s)", s.Name, s.AvailableQuantity, s.Period.String())
		}
		switch x := r.(type) {
		case domain.HumanResource:
			fmt.Fprintf(w, "%s    • Human: %.0f %s (Cost/h: €%.2f)%s\n", newPrefix, x.Quantity, x.Role, x.CostPerH, supplierInfo)
		case domain.MaterialResource:
			fmt.Fprintf(w, "%s    • Material: %.1f %s (Tot: €%.2f)%s\n", newPrefix, x.Quantity, x.Name, x.UnitCost*x.Quantity, supplierInfo)
		case domain.Asset:
			fmt.Fprintf(w, "%s    • Asset: %.1f %s (Tot: €%.2f)%s\n", newPrefix, x.Quantity, x.Name, x.CostPerUse*x.Quantity, supplierInfo)
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

// PrintSupplierRequirementsTo scrive i requisiti di fornitori su w in modo leggibile
func PrintSupplierRequirementsTo(w io.Writer, requirements []engine.SupplierRequirement) {
	if len(requirements) == 0 {
		fmt.Fprintf(w, "Nessun requisito di fornitore trovato.\n")
		return
	}
	fmt.Fprintf(w, "\n=== REQUISITI FORNITORI ===\n")
	fmt.Fprintf(w, "%-30s %15s %15s %15s %10s\n", "Fornitore", "Quantità Richiesta", "Periodo", "Fornitori Necessari", "Fattibile")
	fmt.Fprintf(w, "--------------------------------------------------------------------------------------------\n")
	for _, req := range requirements {
		feasible := "Sì"
		if !req.IsFeasible {
			feasible = "No"
		}
		fmt.Fprintf(w, "%-30s %15.2f %15s %15.2f %10s\n",
			req.SupplierName,
			req.RequiredQuantity,
			req.SupplierPeriod.String(),
			req.SuppliersNeeded,
			feasible)
	}
}

// PrintSupplierRequirements scrive i requisiti di fornitori su os.Stdout
func PrintSupplierRequirements(requirements []engine.SupplierRequirement) {
	PrintSupplierRequirementsTo(os.Stdout, requirements)
}
