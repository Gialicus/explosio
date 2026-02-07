package core

import (
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

func prettyPrintRecursive(activities []*Activity, prefix string, showConnector bool) {
	for i, activity := range activities {
		isLastItem := i == len(activities)-1
		connector := newConnector(showConnector, "", isLastItem)
		ownPrice := fmt.Sprintf("%.2f %s", activity.Price.Value, activity.Price.Currency)
		price := fmt.Sprintf("%.2f %s", activity.CalculatePrice(), activity.Price.Currency)
		ownDuration := fmt.Sprintf("%.0f %s", activity.Duration.Value, activity.Duration.Unit)
		duration := fmt.Sprintf("%.0f %s", activity.CalculateDuration(), activity.Duration.Unit)
		totalFmt := " (" + price + " - " + duration + ")"
		ownFmt := " [" + ownPrice + " - " + ownDuration + "]"
		row := activity.Name + ownFmt + totalFmt
		fmt.Println(prefix + connector + row)
		childPrefix := newChildPrefix(isLastItem, prefix)
		prettyPrintRecursive(activity.Activities, childPrefix, true)
	}
}

func PrettyPrint(activities []*Activity) {
	fmt.Println("Activity Tree:")
	prettyPrintRecursive(activities, "", false)
}
