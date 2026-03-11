// Package core provides CPM (Critical Path Method) with explicit dependencies.
package core

// buildCPMGraph builds predecessor map for all activities. Recursively includes all descendants.
// Predecessors = parent (implicit) + DependsOn (explicit), deduplicated.
func (a *Activity) buildCPMGraph(parent *Activity, all map[*Activity]bool, preds map[*Activity][]*Activity) {
	all[a] = true
	seen := make(map[*Activity]bool)
	var p []*Activity
	if parent != nil && !seen[parent] {
		seen[parent] = true
		p = append(p, parent)
	}
	for _, dep := range a.DependsOn {
		if !seen[dep] {
			seen[dep] = true
			p = append(p, dep)
		}
	}
	preds[a] = p

	for _, child := range a.Activities {
		child.buildCPMGraph(a, all, preds)
	}
}

// topoOrder returns activities in topological order (all predecessors before each activity).
func topoOrder(all map[*Activity]bool, preds map[*Activity][]*Activity) []*Activity {
	inDegree := make(map[*Activity]int)
	for a := range all {
		inDegree[a] = 0
	}
	for a, predList := range preds {
		seen := make(map[*Activity]bool)
		for _, p := range predList {
			if all[p] && !seen[p] {
				seen[p] = true
				inDegree[a]++
			}
		}
	}

	var queue []*Activity
	for a := range all {
		if inDegree[a] == 0 {
			queue = append(queue, a)
		}
	}

	var order []*Activity
	for len(queue) > 0 {
		a := queue[0]
		queue = queue[1:]
		order = append(order, a)
		for other := range all {
			for _, p := range preds[other] {
				if p == a {
					inDegree[other]--
					if inDegree[other] == 0 {
						queue = append(queue, other)
					}
					break
				}
			}
		}
	}
	return order
}

// cpmForwardBackward runs full CPM with dependencies. Returns slack map and project end.
func (a *Activity) cpmForwardBackward() (map[*Activity]SlackInfo, float64) {
	all := make(map[*Activity]bool)
	preds := make(map[*Activity][]*Activity)
	a.buildCPMGraph(nil, all, preds)

	order := topoOrder(all, preds)
	if len(order) == 0 {
		return nil, 0
	}

	// Forward pass
	m := make(map[*Activity]SlackInfo)
	for _, act := range order {
		myHours := act.Duration.ToHours()
		es := 0.0
		for _, p := range preds[act] {
			if info, ok := m[p]; ok && info.EF > es {
				es = info.EF
			}
		}
		ef := es + myHours
		m[act] = SlackInfo{ES: es, EF: ef}
	}

	projectEnd := 0.0
	for _, info := range m {
		if info.EF > projectEnd {
			projectEnd = info.EF
		}
	}

	// Build successors: for each activity, who has it as predecessor?
	succs := make(map[*Activity][]*Activity)
	for act, predList := range preds {
		for _, p := range predList {
			succs[p] = append(succs[p], act)
		}
	}

	// Backward pass (reverse order)
	for i := len(order) - 1; i >= 0; i-- {
		act := order[i]
		info := m[act]
		myHours := act.Duration.ToHours()

		var lf float64
		if len(succs[act]) == 0 {
			lf = projectEnd
		} else {
			lf = projectEnd
			for _, s := range succs[act] {
				if sinfo, ok := m[s]; ok && sinfo.LS < lf {
					lf = sinfo.LS
				}
			}
		}
		ls := lf - myHours
		slack := ls - info.ES
		if slack < 0 {
			slack = 0
		}
		m[act] = SlackInfo{
			ES:    info.ES,
			EF:    info.EF,
			LS:    ls,
			LF:    lf,
			Slack: slack,
		}
	}
	return m, projectEnd
}

// criticalPathFromSlack returns activities with slack 0, ordered by ES (early start).
func criticalPathFromSlack(m map[*Activity]SlackInfo) []*Activity {
	var path []*Activity
	for a, info := range m {
		if info.Slack == 0 {
			path = append(path, a)
		}
	}
	// Sort by ES
	for i := 0; i < len(path); i++ {
		for j := i + 1; j < len(path); j++ {
			if m[path[j]].ES < m[path[i]].ES {
				path[i], path[j] = path[j], path[i]
			}
		}
	}
	return path
}

// hasExplicitDependencies returns true if any activity in the tree has DependsOn.
func (a *Activity) hasExplicitDependencies() bool {
	if len(a.DependsOn) > 0 {
		return true
	}
	for _, child := range a.Activities {
		if child.hasExplicitDependencies() {
			return true
		}
	}
	return false
}
