package engine

import (
	"explosio/lib/domain"
	"explosio/lib/resources"
)

func calculateResourceCost(r domain.Resource, duration int) float64 {
	return r.GetCost(duration)
}

// GetTotalCost calcola il costo totale ricorsivo della filiera
func (e *AnalysisEngine) GetTotalCost(a *domain.Activity) float64 {
	cost := 0.0
	resources.ForEachResource(a, func(r domain.Resource) {
		cost += calculateResourceCost(r, a.Duration)
	})
	for _, sub := range a.SubActivities {
		cost += e.GetTotalCost(sub)
	}
	return cost
}

// CostBreakdown ripartisce il costo totale per categoria (Human, Material, Asset).
type CostBreakdown struct {
	Human    float64
	Material float64
	Asset    float64
}

// GetCostBreakdown calcola il costo totale per categoria sull'albero radicato in a.
func (e *AnalysisEngine) GetCostBreakdown(a *domain.Activity) CostBreakdown {
	if a == nil {
		return CostBreakdown{}
	}
	var human, material, asset float64
	resources.ForEachResource(a, func(r domain.Resource) {
		c := calculateResourceCost(r, a.Duration)
		switch r.(type) {
		case domain.HumanResource:
			human += c
		case domain.MaterialResource:
			material += c
		case domain.Asset:
			asset += c
		}
	})
	for _, sub := range a.SubActivities {
		subB := e.GetCostBreakdown(sub)
		human += subB.Human
		material += subB.Material
		asset += subB.Asset
	}
	return CostBreakdown{Human: human, Material: material, Asset: asset}
}
