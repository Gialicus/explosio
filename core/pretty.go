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
		price := activity.Price.String()
		duration := activity.Duration.String()
		row := activity.Name + " (" + price + " - " + duration + ")"
		fmt.Println(prefix + connector + row)
		childPrefix := newChildPrefix(isLastItem, prefix)
		prettyPrintRecursive(activity.Activities, childPrefix, true)
	}
}

func PrettyPrint(activities []*Activity) {
	fmt.Println("Activity Tree:")
	prettyPrintRecursive(activities, "", false)
}
