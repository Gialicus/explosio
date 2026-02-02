package lib

import (
	"fmt"
	"sort"
)

type AnalysisEngine struct{}

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

// GetTotalCost calcola il costo totale ricorsivo della filiera
func (e *AnalysisEngine) GetTotalCost(a *Activity) float64 {
	cost := 0.0
	for _, h := range a.Humans {
		cost += (h.CostPerH / 60.0) * float64(a.Duration) * h.Quantity
	}
	for _, m := range a.Materials {
		cost += m.UnitCost * m.Quantity
	}
	for _, as := range a.Assets {
		cost += as.CostPerUse * as.Quantity
	}
	for _, sub := range a.SubActivities {
		cost += e.GetTotalCost(sub)
	}
	return cost
}

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

// FinancialReport genera metriche basate su un prezzo di vendita target
type FinancialMetrics struct {
	TotalCost float64
	Margin    float64
	Markup    float64
	IsViable  bool
}

func (e *AnalysisEngine) GetFinancials(root *Activity, sellPrice float64) FinancialMetrics {
	totalCost := e.GetTotalCost(root)
	margin := sellPrice - totalCost
	markup := 0.0
	if totalCost != 0 {
		markup = (margin / totalCost) * 100
	}
	return FinancialMetrics{
		TotalCost: totalCost,
		Margin:    margin,
		Markup:    markup,
		IsViable:  margin > 0,
	}
}

// Validate controlla che l'albero di attività sia valido (durata, MinDuration, assenza di cicli).
func (e *AnalysisEngine) Validate(root *Activity) error {
	if root == nil {
		return fmt.Errorf("root activity is nil")
	}
	seen := make(map[string]bool)
	return e.validateRec(root, seen)
}

func (e *AnalysisEngine) validateRec(a *Activity, seen map[string]bool) error {
	if a == nil {
		return nil
	}
	if seen[a.ID] {
		return fmt.Errorf("activity %s: cycle detected", a.ID)
	}
	seen[a.ID] = true
	if a.Duration < 0 {
		return fmt.Errorf("activity %s: duration negative", a.ID)
	}
	if a.MinDuration < 0 {
		return fmt.Errorf("activity %s: min duration negative", a.ID)
	}
	if a.MinDuration > a.Duration {
		return fmt.Errorf("activity %s: min duration greater than duration", a.ID)
	}
	for _, sub := range a.SubActivities {
		if err := e.validateRec(sub, seen); err != nil {
			return err
		}
	}
	return nil
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
	var path []*Activity
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
	var list []*Activity
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

// CostBreakdown ripartisce il costo totale per categoria (Human, Material, Asset).
type CostBreakdown struct {
	Human    float64
	Material float64
	Asset    float64
}

// GetCostBreakdown calcola il costo totale per categoria sull'albero radicato in a.
func (e *AnalysisEngine) GetCostBreakdown(a *Activity) CostBreakdown {
	if a == nil {
		return CostBreakdown{}
	}
	var human, material, asset float64
	for _, h := range a.Humans {
		human += (h.CostPerH / 60.0) * float64(a.Duration) * h.Quantity
	}
	for _, m := range a.Materials {
		material += m.UnitCost * m.Quantity
	}
	for _, as := range a.Assets {
		asset += as.CostPerUse * as.Quantity
	}
	for _, sub := range a.SubActivities {
		subB := e.GetCostBreakdown(sub)
		human += subB.Human
		material += subB.Material
		asset += subB.Asset
	}
	return CostBreakdown{Human: human, Material: material, Asset: asset}
}

// GetBreakEvenPrice restituisce il prezzo di pareggio (costo totale). Restituisce 0 se root è nil.
func (e *AnalysisEngine) GetBreakEvenPrice(root *Activity) float64 {
	if root == nil {
		return 0
	}
	return e.GetTotalCost(root)
}

// GetFinancialsForPrices restituisce le metriche finanziarie per ogni prezzo in prices.
func (e *AnalysisEngine) GetFinancialsForPrices(root *Activity, prices []float64) []FinancialMetrics {
	out := make([]FinancialMetrics, 0, len(prices))
	for _, p := range prices {
		out = append(out, e.GetFinancials(root, p))
	}
	return out
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
