package engine

import (
	"explosio/lib/domain"
	"explosio/lib/tree"
	"sort"
)

// GetMaxCrashPotential calcola quanto tempo e quanto costo extra comporterebbe il crashing totale
func (e *AnalysisEngine) GetMaxCrashPotential(a *domain.Activity) (timeSaved int, totalExtraCost float64) {
	timeSaved = a.Duration - a.MinDuration
	totalExtraCost = float64(timeSaved) * a.CrashCostStep
	for _, sub := range a.SubActivities {
		sTime, sCost := e.GetMaxCrashPotential(sub)
		timeSaved += sTime
		totalExtraCost += sCost
	}
	return
}

type crashableActivity struct {
	*domain.Activity
	timeSave int
	cost     float64
}

func (e *AnalysisEngine) collectCrashableCritical(root *domain.Activity) []crashableActivity {
	estimatedSize := tree.CountActivities(root)
	list := make([]crashableActivity, 0, estimatedSize)
	e.Walk(root, func(a *domain.Activity) {
		if a.Slack == 0 && a.Duration > a.MinDuration {
			ts := a.Duration - a.MinDuration
			cost := float64(ts) * a.CrashCostStep
			list = append(list, crashableActivity{Activity: a, timeSave: ts, cost: cost})
		}
	})
	sort.Slice(list, func(i, j int) bool {
		return list[i].CrashCostStep < list[j].CrashCostStep
	})
	return list
}

// CrashWithBudget calcola quanto tempo si può risparmiare spendendo al massimo maxExtraCost.
func (e *AnalysisEngine) CrashWithBudget(root *domain.Activity, maxExtraCost float64) (timeSaved int, actualCost float64) {
	list := e.collectCrashableCritical(root)
	for _, ca := range list {
		if actualCost+ca.cost <= maxExtraCost {
			timeSaved += ca.timeSave
			actualCost += ca.cost
		}
	}
	return timeSaved, actualCost
}

// CrashToSaveTime calcola il costo extra minimo per risparmiare almeno targetMinutes minuti.
func (e *AnalysisEngine) CrashToSaveTime(root *domain.Activity, targetMinutes int) (extraCost float64, achieved bool) {
	list := e.collectCrashableCritical(root)
	var timeSaved int
	for _, ca := range list {
		if timeSaved >= targetMinutes {
			break
		}
		extraCost += ca.cost
		timeSaved += ca.timeSave
	}
	return extraCost, timeSaved >= targetMinutes
}
