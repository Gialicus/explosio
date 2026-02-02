package lib

import "sort"

// GetMaxCrashPotential calcola quanto tempo e quanto costo extra comporterebbe il crashing totale
func (e *AnalysisEngine) GetMaxCrashPotential(a *Activity) (timeSaved int, totalExtraCost float64) {
	timeSaved = a.Duration - a.MinDuration
	totalExtraCost = float64(timeSaved) * a.CrashCostStep

	for _, sub := range a.SubActivities {
		sTime, sCost := e.GetMaxCrashPotential(sub)
		timeSaved += sTime
		totalExtraCost += sCost
	}
	return
}

// crashableActivity tiene un'attività sul cammino critico che può essere crashata.
type crashableActivity struct {
	*Activity
	timeSave int
	cost     float64
}

// collectCrashableCritical raccoglie le attività critiche crashabili (Slack==0, Duration > MinDuration).
func (e *AnalysisEngine) collectCrashableCritical(root *Activity) []crashableActivity {
	var list []crashableActivity
	e.Walk(root, func(a *Activity) {
		if a.Slack == 0 && a.Duration > a.MinDuration {
			ts := a.Duration - a.MinDuration
			cost := float64(ts) * a.CrashCostStep
			list = append(list, crashableActivity{Activity: a, timeSave: ts, cost: cost})
		}
	})
	// Ordina per CrashCostStep crescente (costo per minuto; prima le più economiche).
	sort.Slice(list, func(i, j int) bool {
		return list[i].CrashCostStep < list[j].CrashCostStep
	})
	return list
}

// CrashWithBudget calcola quanto tempo si può risparmiare spendendo al massimo maxExtraCost, crashando solo attività sul cammino critico (tutto o niente per attività). Restituisce timeSaved e actualCost.
func (e *AnalysisEngine) CrashWithBudget(root *Activity, maxExtraCost float64) (timeSaved int, actualCost float64) {
	list := e.collectCrashableCritical(root)
	for _, ca := range list {
		if actualCost+ca.cost <= maxExtraCost {
			timeSaved += ca.timeSave
			actualCost += ca.cost
		}
	}
	return timeSaved, actualCost
}

// CrashToSaveTime calcola il costo extra minimo per risparmiare almeno targetMinutes minuti, crashando solo attività sul cammino critico (greedy per costo per minuto). Restituisce extraCost e achieved (true se targetMinutes raggiunto).
func (e *AnalysisEngine) CrashToSaveTime(root *Activity, targetMinutes int) (extraCost float64, achieved bool) {
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
