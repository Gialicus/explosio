// Formatted printing of the activity and material tree (PrettyPrint and helpers).
package core

import (
	"explosio/core/material"
	"fmt"
)

// ANSI color codes for terminal output.
const (
	reset = "\033[0m"
	blue1 = "\033[34m" // standard blue for [] first value
	blue2 = "\033[94m" // bright blue for [] second value
	red1  = "\033[31m" // standard red for () first value
	red2  = "\033[91m" // bright red for () second value
)

// newConnector returns the "├── " or "└── " prefix for the last element.
func newConnector(showConnector bool, connector string, isLastItem bool) string {
	if showConnector {
		connector = "├── "
		if isLastItem {
			connector = "└── "
		}
	}
	return connector
}

// newChildPrefix appends "    " or "│   " to the prefix for children.
func newChildPrefix(isLastItem bool, childPrefix string) string {
	if isLastItem {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}
	return childPrefix
}

// prettyPrintComplexMaterials prints the list of complex materials with prefix and tree connectors.
func prettyPrintComplexMaterials(materials []*material.ComplexMaterial, prefix string, showConnector bool) {
	for i, m := range materials {
		isLastItem := i == len(materials)-1
		connector := newConnector(showConnector, "", isLastItem)
		price := fmt.Sprintf("%.2f %s", m.CalculatePrice(), m.Price.Currency)
		quantity := fmt.Sprintf("%.0f%s", m.MeasurableMaterial.Quantity.Value, m.MeasurableMaterial.Quantity.Unit)
		row := "📦 " + m.Name + " [" + blue1 + price + reset + " - " + blue2 + quantity + reset + "]" + " <" + m.MeasurableMaterial.Name + ">"
		fmt.Println(prefix + connector + row)
	}
}

// prettyPrintCountableMaterials prints the list of countable materials with prefix and tree connectors.
func prettyPrintCountableMaterials(materials []*material.CountableMaterial, prefix string, showConnector bool) {
	for i, m := range materials {
		isLastItem := i == len(materials)-1
		connector := newConnector(showConnector, "", isLastItem)
		price := fmt.Sprintf("%.2f %s", m.CalculatePrice(), m.Price.Currency)
		quantity := fmt.Sprintf("%d", m.Quantity)
		row := "🔢 " + m.Name + " [" + blue1 + price + reset + " - " + blue2 + quantity + reset + "]"
		fmt.Println(prefix + connector + row)
	}
}

// prettyPrintMeasurableMaterials prints the list of measurable materials with prefix and tree connectors.
func prettyPrintMeasurableMaterials(materials []*material.MeasurableMaterial, prefix string, showConnector bool) {
	for i, m := range materials {
		isLastItem := i == len(materials)-1
		connector := newConnector(showConnector, "", isLastItem)
		price := fmt.Sprintf("%.2f %s", m.CalculatePrice(), m.Price.Currency)
		quantity := fmt.Sprintf("%.0f%s", m.Quantity.Value, m.Quantity.Unit)
		row := "📏 " + m.Name + " [" + blue1 + price + reset + " - " + blue2 + quantity + reset + "]"
		fmt.Println(prefix + connector + row)
	}
}

// prettyPrintRecursive walks the tree in depth and prints activities (with critical path icon) and materials.
func prettyPrintRecursive(activities []*Activity, prefix string, showConnector bool, criticalSet map[*Activity]bool) {
	for i, activity := range activities {
		isLastItem := i == len(activities)-1
		connector := newConnector(showConnector, "", isLastItem)
		ownPrice := fmt.Sprintf("%.2f %s", activity.Price.Value, activity.Price.Currency)
		price := fmt.Sprintf("%.2f %s", activity.CalculatePrice(), activity.Price.Currency)
		ownDuration := fmt.Sprintf("%.0f %s", activity.Duration.Value, activity.Duration.Unit)
		duration := fmt.Sprintf("%.0f %s", activity.CalculateDuration(), activity.Duration.Unit)
		totalFmt := " (" + red1 + price + reset + " - " + red2 + duration + reset + ")"
		ownFmt := " [" + blue1 + ownPrice + reset + " - " + blue2 + ownDuration + reset + "]"
		icon := "🟢"
		if criticalSet != nil && criticalSet[activity] {
			icon = "🔴"
		}
		row := icon + " " + activity.Name + ownFmt + totalFmt
		fmt.Println(prefix + connector + row)
		childPrefix := newChildPrefix(isLastItem, prefix)
		prettyPrintComplexMaterials(activity.ComplexMaterials, childPrefix, true)
		prettyPrintCountableMaterials(activity.CountableMaterials, childPrefix, true)
		prettyPrintMeasurableMaterials(activity.MeasurableMaterials, childPrefix, true)
		prettyPrintRecursive(activity.Activities, childPrefix, true, criticalSet)
	}
}

// PrettyPrint prints the activity and material tree. criticalPath is the result of root.CalculateCriticalPath();
// activities on the path are shown with red, others with green. Pass nil for criticalPath to show all as non-critical.
func PrettyPrint(activities []*Activity, criticalPath []*Activity) {
	var criticalSet map[*Activity]bool
	if criticalPath != nil {
		criticalSet = make(map[*Activity]bool)
		for _, a := range criticalPath {
			criticalSet[a] = true
		}
	}
	fmt.Println("--------------------------------")
	fmt.Printf("Legend:\n")
	fmt.Println("🟢: Non-critical activity")
	fmt.Println("🔴: Critical activity")
	fmt.Println("📦: Complex material")
	fmt.Println("🔢: Countable material")
	fmt.Println("📏: Measurable material")
	fmt.Println("[]: Own price and duration (blue variants)")
	fmt.Println("(): Total price and duration (red variants)")
	fmt.Println("<>: Measurable material in complex material")
	fmt.Println("--------------------------------")
	fmt.Println("   Activity and Material Tree:")
	fmt.Println("--------------------------------")
	prettyPrintRecursive(activities, "", false, criticalSet)
}
