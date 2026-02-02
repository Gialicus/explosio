package lib

// ActivityOverride contiene override opzionali per un'attività (solo i campi impostati vengono applicati).
type ActivityOverride struct {
	Duration          *int
	MinDuration       *int
	CrashCostStep     *float64
	HumanCostFactor   float64 // moltiplicatore per CostPerH (1.0 = nessun cambio)
	MaterialCostFactor float64 // moltiplicatore per UnitCost (1.0 = nessun cambio)
	AssetCostFactor   float64 // moltiplicatore per CostPerUse (1.0 = nessun cambio)
}

// Scenario definisce uno scenario what-if: nome, prezzo di vendita e overrides per attività.
type Scenario struct {
	Name      string
	SellPrice float64
	Overrides map[string]ActivityOverride
}

// ScenarioResult contiene il risultato dell'esecuzione di uno scenario (confrontabile).
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

// applyOverridesRec applica gli overrides del scenario al clone (visita ricorsiva).
func applyOverridesRec(a *Activity, overrides map[string]ActivityOverride) {
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
// Il progetto base non viene modificato.
func ApplyScenario(root *Activity, scenario Scenario) *Activity {
	if root == nil {
		return nil
	}
	clone := CloneActivity(root)
	if scenario.Overrides != nil {
		applyOverridesRec(clone, scenario.Overrides)
	}
	return clone
}

// WhatIfEngine esegue scenari what-if (clone + overrides + CPM + risultato).
type WhatIfEngine struct {
	engine *AnalysisEngine
}

// NewWhatIfEngine crea un motore what-if che usa un AnalysisEngine interno.
func NewWhatIfEngine() *WhatIfEngine {
	return &WhatIfEngine{engine: &AnalysisEngine{}}
}

// RunScenario applica lo scenario al base, esegue CPM sul clone e restituisce clone e risultato.
func (w *WhatIfEngine) RunScenario(base *Activity, scenario Scenario) (clone *Activity, result ScenarioResult) {
	clone = ApplyScenario(base, scenario)
	if clone == nil {
		return nil, ScenarioResult{ScenarioName: scenario.Name}
	}
	w.engine.ComputeCPM(clone)
	totalCost := w.engine.GetTotalCost(clone)
	sellPrice := scenario.SellPrice
	fin := w.engine.GetFinancials(clone, sellPrice)
	criticalPath := w.engine.GetCriticalPath(clone)
	timeSaved, extraCrash := w.engine.GetMaxCrashPotential(clone)
	result = ScenarioResult{
		ScenarioName:      scenario.Name,
		TotalDuration:     clone.EF,
		TotalCost:         totalCost,
		Margin:            fin.Margin,
		Markup:            fin.Markup,
		IsViable:          fin.IsViable,
		CriticalPathCount: len(criticalPath),
		TimeSaved:         timeSaved,
		ExtraCrashCost:    extraCrash,
	}
	return clone, result
}

// RunScenarios esegue più scenari sul progetto base e restituisce i risultati nello stesso ordine.
func (w *WhatIfEngine) RunScenarios(base *Activity, scenarios []Scenario) []ScenarioResult {
	results := make([]ScenarioResult, 0, len(scenarios))
	for _, s := range scenarios {
		_, res := w.RunScenario(base, s)
		results = append(results, res)
	}
	return results
}
