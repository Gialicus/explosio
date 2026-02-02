// Package lib fornisce funzionalità per l'analisi di progetti strutturati come alberi di attività.
// Questo modulo gestisce il calcolo CPM (Critical Path Method) e l'identificazione del cammino critico.
package lib

import "sort"

// countActivities conta il numero totale di attività nell'albero
func countActivities(root *Activity) int {
	if root == nil {
		return 0
	}
	count := 1
	for _, sub := range root.SubActivities {
		count += countActivities(sub)
	}
	return count
}

// ComputeCPM esegue il calcolo dei tempi e identifica il cammino critico
func (e *AnalysisEngine) ComputeCPM(root *Activity) {
	if root == nil {
		return
	}
	e.forwardPass(root)
	e.backwardPass(root, root.EF)
}

func (e *AnalysisEngine) forwardPass(a *Activity) int {
	if len(a.SubActivities) == 0 {
		a.ES = 0
		a.EF = a.Duration
		return a.EF
	}
	maxEF := 0
	for _, sub := range a.SubActivities {
		ef := e.forwardPass(sub)
		if ef > maxEF {
			maxEF = ef
		}
	}
	a.ES = maxEF
	a.EF = a.ES + a.Duration
	return a.EF
}

func (e *AnalysisEngine) backwardPass(a *Activity, parentLS int) {
	a.LF = parentLS
	a.LS = a.LF - a.Duration
	a.Slack = a.LS - a.ES
	for _, sub := range a.SubActivities {
		e.backwardPass(sub, a.LS)
	}
}

// GetTotalDuration restituisce la durata totale del progetto (root.EF). Richiede che ComputeCPM sia già stato eseguito. Restituisce 0 se root è nil.
func (e *AnalysisEngine) GetTotalDuration(root *Activity) int {
	if root == nil {
		return 0
	}
	return root.EF
}

// GetCriticalPath restituisce le attività sul cammino critico (Slack == 0) in ordine pre-order. Richiede che ComputeCPM sia già stato eseguito.
func (e *AnalysisEngine) GetCriticalPath(root *Activity) []*Activity {
	if root == nil {
		return nil
	}
	// Pre-alloca la slice con una stima (tutte le attività potrebbero essere critiche)
	estimatedSize := countActivities(root)
	path := make([]*Activity, 0, estimatedSize)
	e.collectCritical(root, &path)
	return path
}

func (e *AnalysisEngine) collectCritical(a *Activity, path *[]*Activity) {
	if a == nil {
		return
	}
	if a.Slack == 0 {
		*path = append(*path, a)
	}
	for _, sub := range a.SubActivities {
		e.collectCritical(sub, path)
	}
}

// Walk attraversa l'albero in pre-order chiamando f su ogni attività.
func (e *AnalysisEngine) Walk(root *Activity, f func(*Activity)) {
	if root == nil {
		return
	}
	f(root)
	for _, sub := range root.SubActivities {
		e.Walk(sub, f)
	}
}

// CPMSummary contiene il riepilogo del CPM. ComputeCPM deve essere già stato eseguito prima di chiamare GetCPMSummary.
type CPMSummary struct {
	TotalDuration int
	CriticalPath  []*Activity
	ActivityCount int
}

// GetCPMSummary restituisce TotalDuration (root.EF), CriticalPath e ActivityCount. ComputeCPM deve essere già stato eseguito.
func (e *AnalysisEngine) GetCPMSummary(root *Activity) CPMSummary {
	if root == nil {
		return CPMSummary{}
	}
	var count int
	e.Walk(root, func(*Activity) { count++ })
	return CPMSummary{
		TotalDuration: root.EF,
		CriticalPath:  e.GetCriticalPath(root),
		ActivityCount: count,
	}
}

// ActivitiesByES restituisce tutte le attività ordinate per ES (e per EF in caso di pari ES). Richiede che ComputeCPM sia già stato eseguito.
func (e *AnalysisEngine) ActivitiesByES(root *Activity) []*Activity {
	if root == nil {
		return nil
	}
	// Pre-alloca la slice con il numero esatto di attività
	totalCount := countActivities(root)
	list := make([]*Activity, 0, totalCount)
	e.Walk(root, func(a *Activity) { list = append(list, a) })
	sort.Slice(list, func(i, j int) bool {
		if list[i].ES != list[j].ES {
			return list[i].ES < list[j].ES
		}
		if list[i].EF != list[j].EF {
			return list[i].EF < list[j].EF
		}
		return list[i].ID < list[j].ID
	})
	return list
}
