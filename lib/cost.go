package lib

// Metodi helper per calcolo costi delle singole risorse

// calculateHumanCost calcola il costo di una risorsa umana
func (e *AnalysisEngine) calculateHumanCost(h HumanResource, duration int) float64 {
	return h.GetCost(duration)
}

// calculateMaterialCost calcola il costo di un materiale
func (e *AnalysisEngine) calculateMaterialCost(m MaterialResource) float64 {
	return m.GetCost(0) // duration non usato per materiali
}

// calculateAssetCost calcola il costo di un asset
func (e *AnalysisEngine) calculateAssetCost(as Asset) float64 {
	return as.GetCost(0) // duration non usato per asset
}

// GetTotalCost calcola il costo totale ricorsivo della filiera
func (e *AnalysisEngine) GetTotalCost(a *Activity) float64 {
	cost := 0.0
	for _, h := range a.Humans {
		cost += e.calculateHumanCost(h, a.Duration)
	}
	for _, m := range a.Materials {
		cost += e.calculateMaterialCost(m)
	}
	for _, as := range a.Assets {
		cost += e.calculateAssetCost(as)
	}
	for _, sub := range a.SubActivities {
		cost += e.GetTotalCost(sub)
	}
	return cost
}

// CostBreakdown ripartisce il costo totale per categoria (Human, Material, Asset).
// I fornitori non contribuiscono ai costi, servono solo per validare la capacità.
type CostBreakdown struct {
	Human    float64
	Material float64
	Asset    float64
}

// GetCostBreakdown calcola il costo totale per categoria sull'albero radicato in a.
// I fornitori non contribuiscono ai costi, servono solo per validare la capacità.
func (e *AnalysisEngine) GetCostBreakdown(a *Activity) CostBreakdown {
	if a == nil {
		return CostBreakdown{}
	}
	var human, material, asset float64
	for _, h := range a.Humans {
		human += e.calculateHumanCost(h, a.Duration)
	}
	for _, m := range a.Materials {
		material += e.calculateMaterialCost(m)
	}
	for _, as := range a.Assets {
		asset += e.calculateAssetCost(as)
	}
	for _, sub := range a.SubActivities {
		subB := e.GetCostBreakdown(sub)
		human += subB.Human
		material += subB.Material
		asset += subB.Asset
	}
	return CostBreakdown{Human: human, Material: material, Asset: asset}
}
