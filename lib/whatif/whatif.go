package whatif

import (
	"explosio/lib/clone"
	"explosio/lib/domain"
	"explosio/lib/engine"
)

// ActivityOverride contiene override opzionali per un'attività.
type ActivityOverride struct {
	Duration            *int
	MinDuration         *int
	CrashCostStep       *float64
	HumanCostFactor     float64
	MaterialCostFactor  float64
	AssetCostFactor     float64
}

// Scenario definisce uno scenario what-if: nome, prezzo di vendita e overrides per attività.
type Scenario struct {
	Name      string
	SellPrice float64
	Overrides map[string]ActivityOverride
}

// ScenarioResult contiene il risultato dell'esecuzione di uno scenario.
type ScenarioResult struct {
	ScenarioName      string
	TotalDuration     int
	TotalCost         float64
	Margin            float64
	Markup            float64
	IsViable          bool
	CriticalPathCount int
	TimeSaved         int
	ExtraCrashCost    float64
}

func applyOverridesRec(a *domain.Activity, overrides map[string]ActivityOverride) {
	if a == nil {
		return
	}
	if o, ok := overrides[a.ID]; ok {
		if o.Duration != nil {
			a.Duration = *o.Duration
		}
		if o.MinDuration != nil {
			a.MinDuration = *o.MinDuration
		}
		if o.CrashCostStep != nil {
			a.CrashCostStep = *o.CrashCostStep
		}
		if o.HumanCostFactor != 1.0 {
			for i := range a.Humans {
				a.Humans[i].CostPerH *= o.HumanCostFactor
			}
		}
		if o.MaterialCostFactor != 1.0 {
			for i := range a.Materials {
				a.Materials[i].UnitCost *= o.MaterialCostFactor
			}
		}
		if o.AssetCostFactor != 1.0 {
			for i := range a.Assets {
				a.Assets[i].CostPerUse *= o.AssetCostFactor
			}
		}
	}
	for _, sub := range a.SubActivities {
		applyOverridesRec(sub, overrides)
	}
}

// ApplyScenario restituisce un clone del progetto base con gli overrides dello scenario applicati.
func ApplyScenario(root *domain.Activity, scenario Scenario) *domain.Activity {
	if root == nil {
		return nil
	}
	cl := clone.CloneActivity(root)
	if scenario.Overrides != nil {
		applyOverridesRec(cl, scenario.Overrides)
	}
	return cl
}

// WhatIfEngine esegue scenari what-if (clone + overrides + CPM + risultato).
type WhatIfEngine struct {
	engine *engine.AnalysisEngine
}

// NewWhatIfEngine crea un motore what-if che usa un AnalysisEngine interno.
func NewWhatIfEngine() *WhatIfEngine {
	return &WhatIfEngine{engine: &engine.AnalysisEngine{}}
}

// RunScenario applica lo scenario al base, esegue CPM sul clone e restituisce clone e risultato.
func (w *WhatIfEngine) RunScenario(base *domain.Activity, scenario Scenario) (cl *domain.Activity, result ScenarioResult) {
	cl = ApplyScenario(base, scenario)
	if cl == nil {
		return nil, ScenarioResult{ScenarioName: scenario.Name}
	}
	w.engine.ComputeCPM(cl)
	totalCost := w.engine.GetTotalCost(cl)
	sellPrice := scenario.SellPrice
	fin := w.engine.GetFinancials(cl, sellPrice)
	criticalPath := w.engine.GetCriticalPath(cl)
	timeSaved, extraCrash := w.engine.GetMaxCrashPotential(cl)
	result = ScenarioResult{
		ScenarioName:      scenario.Name,
		TotalDuration:     cl.EF,
		TotalCost:         totalCost,
		Margin:            fin.Margin,
		Markup:            fin.Markup,
		IsViable:          fin.IsViable,
		CriticalPathCount: len(criticalPath),
		TimeSaved:         timeSaved,
		ExtraCrashCost:    extraCrash,
	}
	return cl, result
}

// RunScenarios esegue più scenari sul progetto base e restituisce i risultati nello stesso ordine.
func (w *WhatIfEngine) RunScenarios(base *domain.Activity, scenarios []Scenario) []ScenarioResult {
	results := make([]ScenarioResult, 0, len(scenarios))
	for _, s := range scenarios {
		_, res := w.RunScenario(base, s)
		results = append(results, res)
	}
	return results
}
