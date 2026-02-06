package mermaid

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"explosio/lib/domain"
	"explosio/lib/resources"
	"explosio/lib/tree"
)

func mermaidLabel(s string) string {
	s = strings.ReplaceAll(s, "\"", "#quot;")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)
	return s
}

func mermaidSafeID(activityID string) string {
	return strings.ReplaceAll(activityID, "-", "_")
}

func printMermaidRec(w io.Writer, a *domain.Activity) {
	if a == nil {
		return
	}
	label := mermaidLabel(a.Name + " (" + strconv.Itoa(a.Duration) + " min)")
	fmt.Fprintf(w, "  %s[\"%s\"]\n", a.ID, label)
	prefix := mermaidSafeID(a.ID)
	var hi, mi, ai int
	resources.ForEachResource(a, func(r domain.Resource) {
		switch x := r.(type) {
		case domain.HumanResource:
			resID := fmt.Sprintf("%s_H_%d", prefix, hi)
			hi++
			lbl := mermaidLabel(fmt.Sprintf("%s (€%.0f/h)", x.Role, x.CostPerH))
			fmt.Fprintf(w, "  %s([\"%s\"])\n", resID, lbl)
			fmt.Fprintf(w, "  %s --> %s\n", a.ID, resID)
			if x.Supplier != nil {
				supplierID := fmt.Sprintf("%s_SUP_H_%d", prefix, hi-1)
				supplierLbl := mermaidLabel(fmt.Sprintf("%s: max %.1f/%s [validatore]", x.Supplier.Name, x.Supplier.AvailableQuantity, x.Supplier.Period.String()))
				fmt.Fprintf(w, "  %s{\"%s\"}\n", supplierID, supplierLbl)
				fmt.Fprintf(w, "  %s -.-> %s\n", resID, supplierID)
			}
		case domain.MaterialResource:
			resID := fmt.Sprintf("%s_M_%d", prefix, mi)
			mi++
			lbl := mermaidLabel(fmt.Sprintf("%s (€%.2f)", x.Name, x.UnitCost*x.Quantity))
			fmt.Fprintf(w, "  %s[/\"%s\"/]\n", resID, lbl)
			fmt.Fprintf(w, "  %s --> %s\n", a.ID, resID)
			if x.Supplier != nil {
				supplierID := fmt.Sprintf("%s_SUP_M_%d", prefix, mi-1)
				supplierLbl := mermaidLabel(fmt.Sprintf("%s: max %.1f/%s [validatore]", x.Supplier.Name, x.Supplier.AvailableQuantity, x.Supplier.Period.String()))
				fmt.Fprintf(w, "  %s{\"%s\"}\n", supplierID, supplierLbl)
				fmt.Fprintf(w, "  %s -.-> %s\n", resID, supplierID)
			}
		case domain.Asset:
			resID := fmt.Sprintf("%s_A_%d", prefix, ai)
			ai++
			lbl := mermaidLabel(fmt.Sprintf("%s (€%.2f)", x.Name, x.CostPerUse*x.Quantity))
			fmt.Fprintf(w, "  %s[(\"%s\")]\n", resID, lbl)
			fmt.Fprintf(w, "  %s --> %s\n", a.ID, resID)
			if x.Supplier != nil {
				supplierID := fmt.Sprintf("%s_SUP_A_%d", prefix, ai-1)
				supplierLbl := mermaidLabel(fmt.Sprintf("%s: max %.1f/%s [validatore]", x.Supplier.Name, x.Supplier.AvailableQuantity, x.Supplier.Period.String()))
				fmt.Fprintf(w, "  %s{\"%s\"}\n", supplierID, supplierLbl)
				fmt.Fprintf(w, "  %s -.-> %s\n", resID, supplierID)
			}
		}
	})
	for _, sub := range a.SubActivities {
		fmt.Fprintf(w, "  %s --> %s\n", a.ID, sub.ID)
	}
	for _, sub := range a.SubActivities {
		printMermaidRec(w, sub)
	}
}

// PrintMermaidTo scrive l'albero delle attività in formato Mermaid flowchart su w.
func PrintMermaidTo(w io.Writer, root *domain.Activity) {
	if root == nil {
		return
	}
	fmt.Fprintf(w, "flowchart TB\n")
	printMermaidRec(w, root)
	var critical []string
	tree.Walk(root, func(a *domain.Activity) {
		if a.Slack == 0 {
			critical = append(critical, a.ID)
		}
	})
	for _, id := range critical {
		fmt.Fprintf(w, "  style %s fill:#f96,stroke:#333\n", id)
	}
}

// PrintMermaid scrive il diagramma Mermaid su os.Stdout.
func PrintMermaid(root *domain.Activity) {
	PrintMermaidTo(os.Stdout, root)
}

// GenerateMermaid genera il diagramma Mermaid come stringa
func GenerateMermaid(root *domain.Activity) string {
	if root == nil {
		return ""
	}
	var buf strings.Builder
	PrintMermaidTo(&buf, root)
	return buf.String()
}

// WriteMermaidToFile scrive il diagramma Mermaid su file.
func WriteMermaidToFile(root *domain.Activity, path string) error {
	if root == nil {
		return nil
	}
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	PrintMermaidTo(f, root)
	if err := f.Close(); err != nil {
		return fmt.Errorf("close file: %w", err)
	}
	return nil
}
