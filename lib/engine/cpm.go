package engine

import (
	"explosio/lib/domain"
	"explosio/lib/tree"
	"sort"
)

// ComputeCPM esegue il calcolo dei tempi e identifica il cammino critico
func (e *AnalysisEngine) ComputeCPM(root *domain.Activity) {
	if root == nil {
		return
	}
	e.forwardPass(root)
	e.backwardPass(root, root.EF)
}

func (e *AnalysisEngine) forwardPass(a *domain.Activity) int {
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

func (e *AnalysisEngine) backwardPass(a *domain.Activity, parentLS int) {
	a.LF = parentLS
	a.LS = a.LF - a.Duration
	a.Slack = a.LS - a.ES
	for _, sub := range a.SubActivities {
		e.backwardPass(sub, a.LS)
	}
}

// GetTotalDuration restituisce la durata totale del progetto (root.EF).
func (e *AnalysisEngine) GetTotalDuration(root *domain.Activity) int {
	if root == nil {
		return 0
	}
	return root.EF
}

// GetCriticalPath restituisce le attività sul cammino critico (Slack == 0) in ordine pre-order.
func (e *AnalysisEngine) GetCriticalPath(root *domain.Activity) []*domain.Activity {
	if root == nil {
		return nil
	}
	estimatedSize := tree.CountActivities(root)
	path := make([]*domain.Activity, 0, estimatedSize)
	e.collectCritical(root, &path)
	return path
}

func (e *AnalysisEngine) collectCritical(a *domain.Activity, path *[]*domain.Activity) {
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

// CPMSummary contiene il riepilogo del CPM.
type CPMSummary struct {
	TotalDuration int
	CriticalPath  []*domain.Activity
	ActivityCount int
}

// GetCPMSummary restituisce TotalDuration, CriticalPath e ActivityCount.
func (e *AnalysisEngine) GetCPMSummary(root *domain.Activity) CPMSummary {
	if root == nil {
		return CPMSummary{}
	}
	var count int
	e.Walk(root, func(*domain.Activity) { count++ })
	return CPMSummary{
		TotalDuration: root.EF,
		CriticalPath:  e.GetCriticalPath(root),
		ActivityCount: count,
	}
}

// ActivitiesByES restituisce tutte le attività ordinate per ES (e per EF in caso di pari ES).
func (e *AnalysisEngine) ActivitiesByES(root *domain.Activity) []*domain.Activity {
	if root == nil {
		return nil
	}
	totalCount := tree.CountActivities(root)
	list := make([]*domain.Activity, 0, totalCount)
	e.Walk(root, func(a *domain.Activity) { list = append(list, a) })
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
