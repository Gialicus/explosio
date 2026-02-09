package core

// This file contains calculation methods on Activity (price, duration, quantity, critical path).

// CalculatePrice returns the total price (activity plus all materials and sub-activities).
func (a *Activity) CalculatePrice() float64 {
	price := a.Price.Value
	for _, complexMaterial := range a.ComplexMaterials {
		price += complexMaterial.CalculatePrice()
	}
	for _, countableMaterial := range a.CountableMaterials {
		price += countableMaterial.CalculatePrice()
	}
	for _, measurableMaterial := range a.MeasurableMaterials {
		price += measurableMaterial.CalculatePrice()
	}
	for _, activity := range a.Activities {
		price += activity.CalculatePrice()
	}
	for _, humanResource := range a.HumanResources {
		price += humanResource.CalculatePrice()
	}
	for _, asset := range a.Assets {
		price += asset.CalculatePrice()
	}
	return price
}

// CalculateDuration returns the total duration (activity plus sub-activities) in the same unit as the activity's Duration (e.g. days).
func (a *Activity) CalculateDuration() float64 {
	duration := a.Duration.Value
	for _, activity := range a.Activities {
		duration += activity.CalculateDuration()
	}
	return duration
}

// CalculateCriticalPath returns the critical path (classic CPM interpretation): the longest path from root to leaf when sub-activities are in parallel. Any delay on this path delays the whole project.
func (a *Activity) CalculateCriticalPath() []*Activity {
	path, _ := a.criticalPathAndDuration()
	return path
}

// criticalPathAndDuration returns the critical path and its total duration in hours.
// For a leaf (no children), the path is just this activity.
// For a node with children, it picks the child whose subtree has the longest
// total duration and appends that child's critical path to this activity.
func (a *Activity) criticalPathAndDuration() ([]*Activity, float64) {
	myHours := a.Duration.ToHours()
	if len(a.Activities) == 0 {
		return []*Activity{a}, myHours
	}
	var bestPath []*Activity
	bestHours := -1.0
	for _, child := range a.Activities {
		childPath, childHours := child.criticalPathAndDuration()
		total := myHours + childHours
		if total > bestHours {
			bestHours = total
			bestPath = childPath
		}
	}
	return append([]*Activity{a}, bestPath...), bestHours
}
