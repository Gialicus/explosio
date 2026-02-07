package core

import (
	"explosio/core/material"
	"fmt"
)

func newConnector(showConnector bool, connector string, isLastItem bool) string {
	if showConnector {
		connector = "├── "
		if isLastItem {
			connector = "└── "
		}
	}
	return connector
}

func newChildPrefix(isLastItem bool, childPrefix string) string {
	if isLastItem {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}
	return childPrefix
}

func prettyPrintComplexMaterials(materials []*material.ComplexMaterial, prefix string, showConnector bool) {
	for i, m := range materials {
		isLastItem := i == len(materials)-1
		connector := newConnector(showConnector, "", isLastItem)
		price := fmt.Sprintf("%.2f %s", m.CalculatePrice(), m.Price.Currency)
		quantity := fmt.Sprintf("%.0f%s", m.MeasurableMaterial.Quantity.Value, m.MeasurableMaterial.Quantity.Unit)
		row := "📦 " + m.Name + " (" + price + " - " + quantity + ")" + " [" + m.MeasurableMaterial.Name + "]"
		fmt.Println(prefix + connector + row)
	}
}

func prettyPrintCountableMaterials(materials []*material.CountableMaterial, prefix string, showConnector bool) {
	for i, m := range materials {
		isLastItem := i == len(materials)-1
		connector := newConnector(showConnector, "", isLastItem)
		price := fmt.Sprintf("%.2f %s", m.CalculatePrice(), m.Price.Currency)
		quantity := fmt.Sprintf("%d", m.Quantity)
		row := "🔢 " + m.Name + " (" + price + " - " + quantity + ")"
		fmt.Println(prefix + connector + row)
	}
}

func prettyPrintMeasurableMaterials(materials []*material.MeasurableMaterial, prefix string, showConnector bool) {
	for i, m := range materials {
		isLastItem := i == len(materials)-1
		connector := newConnector(showConnector, "", isLastItem)
		price := fmt.Sprintf("%.2f %s", m.CalculatePrice(), m.Price.Currency)
		quantity := fmt.Sprintf("%.0f%s", m.Quantity.Value, m.Quantity.Unit)
		row := "📏 " + m.Name + " (" + price + " - " + quantity + ")"
		fmt.Println(prefix + connector + row)
	}
}

func prettyPrintRecursive(activities []*Activity, prefix string, showConnector bool, criticalSet map[*Activity]bool) {
	for i, activity := range activities {
		isLastItem := i == len(activities)-1
		connector := newConnector(showConnector, "", isLastItem)
		ownPrice := fmt.Sprintf("%.2f %s", activity.Price.Value, activity.Price.Currency)
		price := fmt.Sprintf("%.2f %s", activity.CalculatePrice(), activity.Price.Currency)
		ownDuration := fmt.Sprintf("%.0f %s", activity.Duration.Value, activity.Duration.Unit)
		duration := fmt.Sprintf("%.0f %s", activity.CalculateDuration(), activity.Duration.Unit)
		totalFmt := " (" + price + " - " + duration + ")"
		ownFmt := " [" + ownPrice + " - " + ownDuration + "]"
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

// PrettyPrint prints the activity tree. criticalPath is the result of root.CalculateCriticalPath();
// activities on the path are shown with 🔴, others with 🟢. Pass nil to show all as 🟢.
func PrettyPrint(activities []*Activity, criticalPath []*Activity) {
	var criticalSet map[*Activity]bool
	if criticalPath != nil {
		criticalSet = make(map[*Activity]bool)
		for _, a := range criticalPath {
			criticalSet[a] = true
		}
	}
	fmt.Println("Activity Tree:")
	prettyPrintRecursive(activities, "", false, criticalSet)
}
