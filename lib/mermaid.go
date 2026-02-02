package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// mermaidLabel restituisce una stringa sicura per l'uso come label Mermaid dentro ["..."].
func mermaidLabel(s string) string {
	s = strings.ReplaceAll(s, "\"", "#quot;")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)
	return s
}

// mermaidSafeID restituisce l'ID attività con '-' sostituito da '_' per usarlo come prefisso negli ID dei nodi risorsa.
func mermaidSafeID(activityID string) string {
	return strings.ReplaceAll(activityID, "-", "_")
}

// printMermaidRec visita l'albero e scrive nodi (attività e risorse) e archi in formato Mermaid flowchart.
func printMermaidRec(w io.Writer, a *Activity) {
	if a == nil {
		return
	}
	label := mermaidLabel(a.Name + " (" + strconv.Itoa(a.Duration) + " min)")
	fmt.Fprintf(w, "  %s[\"%s\"]\n", a.ID, label)
	prefix := mermaidSafeID(a.ID)
	for i, h := range a.Humans {
		resID := fmt.Sprintf("%s_H_%d", prefix, i)
		lbl := mermaidLabel(fmt.Sprintf("%s (€%.0f/h)", h.Role, h.CostPerH))
		fmt.Fprintf(w, "  %s([\"%s\"])\n", resID, lbl)
		fmt.Fprintf(w, "  %s --> %s\n", a.ID, resID)
	}
	for i, m := range a.Materials {
		resID := fmt.Sprintf("%s_M_%d", prefix, i)
		lbl := mermaidLabel(fmt.Sprintf("%s (€%.2f)", m.Name, m.UnitCost*m.Quantity))
		fmt.Fprintf(w, "  %s[/\"%s\"/]\n", resID, lbl)
		fmt.Fprintf(w, "  %s --> %s\n", a.ID, resID)
	}
	for i, as := range a.Assets {
		resID := fmt.Sprintf("%s_A_%d", prefix, i)
		lbl := mermaidLabel(fmt.Sprintf("%s (€%.2f)", as.Name, as.CostPerUse*as.Quantity))
		fmt.Fprintf(w, "  %s[(\"%s\")]\n", resID, lbl)
		fmt.Fprintf(w, "  %s --> %s\n", a.ID, resID)
	}
	for _, sub := range a.SubActivities {
		fmt.Fprintf(w, "  %s --> %s\n", a.ID, sub.ID)
	}
	for _, sub := range a.SubActivities {
		printMermaidRec(w, sub)
	}
}

// collectCriticalIDs raccoglie gli ID delle attività con Slack == 0 (cammino critico).
func collectCriticalIDs(a *Activity, ids *[]string) {
	if a == nil {
		return
	}
	if a.Slack == 0 {
		*ids = append(*ids, a.ID)
	}
	for _, sub := range a.SubActivities {
		collectCriticalIDs(sub, ids)
	}
}

// PrintMermaidTo scrive l'albero delle attività in formato Mermaid flowchart su w.
// Se root è nil non scrive nulla. Se ComputeCPM è già stato eseguito, evidenzia i nodi del cammino critico (Slack == 0).
func PrintMermaidTo(w io.Writer, root *Activity) {
	if root == nil {
		return
	}
	fmt.Fprintf(w, "flowchart TB\n")
	printMermaidRec(w, root)
	var critical []string
	collectCriticalIDs(root, &critical)
	for _, id := range critical {
		fmt.Fprintf(w, "  style %s fill:#f96,stroke:#333\n", id)
	}
}

// PrintMermaid scrive il diagramma Mermaid su os.Stdout.
func PrintMermaid(root *Activity) {
	PrintMermaidTo(os.Stdout, root)
}

// WriteMermaidToFile scrive il diagramma Mermaid su file. Crea la directory se non esiste.
// Se root è nil non scrive nulla e restituisce nil.
func WriteMermaidToFile(root *Activity, path string) error {
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
