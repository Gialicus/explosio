package core

// This file contains calculation methods on Activity (price, duration, quantity, critical path).

// pricers returns all direct price contributors (materials, sub-activities, human resources, assets).
func (a *Activity) pricers() []Pricer {
	var p []Pricer
	for _, m := range a.ComplexMaterials {
		p = append(p, m)
	}
	for _, m := range a.CountableMaterials {
		p = append(p, m)
	}
	for _, m := range a.MeasurableMaterials {
		p = append(p, m)
	}
	for _, child := range a.Activities {
		p = append(p, child)
	}
	for _, h := range a.HumanResources {
		p = append(p, h)
	}
	for _, as := range a.Assets {
		p = append(p, as)
	}
	return p
}

// CalculatePrice returns the total price (activity plus all materials and sub-activities).
func (a *Activity) CalculatePrice() float64 {
	price := a.Price.Value
	for _, p := range a.pricers() {
		price += p.CalculatePrice()
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

// calculateDurationHours returns the total subtree duration in hours (for filtering).
func (a *Activity) calculateDurationHours() float64 {
	hours := a.Duration.ToHours()
	for _, child := range a.Activities {
		hours += child.calculateDurationHours()
	}
	return hours
}

// CalculateCriticalPath returns the critical path. With explicit DependsOn, uses full CPM; otherwise uses tree-based longest path.
func (a *Activity) CalculateCriticalPath() []*Activity {
	if a.hasExplicitDependencies() {
		m, _ := a.cpmForwardBackward()
		return criticalPathFromSlack(m)
	}
	path, _ := a.criticalPathAndDuration()
	return path
}

// CostBreakdown holds the price breakdown by category for an activity tree.
type CostBreakdown struct {
	Activities float64 // Direct activity prices (own Price.Value)
	Materials  float64 // Complex + countable + measurable materials
	Human      float64 // Human resources
	Assets     float64 // Assets
}

// Total returns the sum of all categories.
func (c *CostBreakdown) Total() float64 {
	return c.Activities + c.Materials + c.Human + c.Assets
}

// CostBreakdown returns the price breakdown by category for this activity and all descendants.
func (a *Activity) CostBreakdown() CostBreakdown {
	var cb CostBreakdown
	cb.Activities = a.Price.Value
	for _, m := range a.ComplexMaterials {
		cb.Materials += m.CalculatePrice()
	}
	for _, m := range a.CountableMaterials {
		cb.Materials += m.CalculatePrice()
	}
	for _, m := range a.MeasurableMaterials {
		cb.Materials += m.CalculatePrice()
	}
	for _, h := range a.HumanResources {
		cb.Human += h.CalculatePrice()
	}
	for _, as := range a.Assets {
		cb.Assets += as.CalculatePrice()
	}
	for _, child := range a.Activities {
		childCB := child.CostBreakdown()
		cb.Activities += childCB.Activities
		cb.Materials += childCB.Materials
		cb.Human += childCB.Human
		cb.Assets += childCB.Assets
	}
	return cb
}

// SlackInfo holds ES, EF, LS, LF and Slack for an activity (in hours).
type SlackInfo struct {
	ES    float64 // Early Start
	EF    float64 // Early Finish
	LS    float64 // Late Start
	LF    float64 // Late Finish
	Slack float64 // LS - ES (or LF - EF)
}

// CalculateSlack returns a map of activity -> SlackInfo for all activities in the tree.
// Activities on the critical path have Slack = 0.
func (a *Activity) CalculateSlack() map[*Activity]SlackInfo {
	if a.hasExplicitDependencies() {
		m, _ := a.cpmForwardBackward()
		return m
	}
	_, projectEnd := a.criticalPathAndDuration()
	m := make(map[*Activity]SlackInfo)
	a.forwardPass(0, m)
	a.backwardPass(projectEnd, m)
	return m
}

func (a *Activity) forwardPass(parentEF float64, m map[*Activity]SlackInfo) {
	myHours := a.Duration.ToHours()
	es := parentEF
	ef := es + myHours

	// Children start when this activity's own work finishes (parallel). Pass the same start time to all.
	childStart := ef
	if len(a.Activities) > 0 {
		for _, child := range a.Activities {
			child.forwardPass(childStart, m)
			childInfo := m[child]
			if childInfo.EF > ef {
				ef = childInfo.EF
			}
		}
	}

	m[a] = SlackInfo{ES: es, EF: ef}
}

func (a *Activity) backwardPass(projectEnd float64, m map[*Activity]SlackInfo) {
	info, ok := m[a]
	if !ok {
		return
	}

	var lf float64
	if len(a.Activities) == 0 {
		lf = projectEnd
	} else {
		lf = projectEnd
		for _, child := range a.Activities {
			child.backwardPass(projectEnd, m)
			childInfo := m[child]
			if childInfo.LS < lf {
				lf = childInfo.LS
			}
		}
	}

	myHours := a.Duration.ToHours()
	ls := lf - myHours
	slack := ls - info.ES
	if slack < 0 {
		slack = 0
	}
	m[a] = SlackInfo{
		ES:    info.ES,
		EF:    info.EF,
		LS:    ls,
		LF:    lf,
		Slack: slack,
	}
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
