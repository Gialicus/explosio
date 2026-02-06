package engine

import "explosio/lib/domain"

// FinancialMetrics genera metriche basate su un prezzo di vendita target
type FinancialMetrics struct {
	TotalCost float64
	Margin    float64
	Markup    float64
	IsViable  bool
}

// GetFinancials calcola le metriche finanziarie per un prezzo di vendita
func (e *AnalysisEngine) GetFinancials(root *domain.Activity, sellPrice float64) FinancialMetrics {
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

// GetBreakEvenPrice restituisce il prezzo di pareggio (costo totale).
func (e *AnalysisEngine) GetBreakEvenPrice(root *domain.Activity) float64 {
	if root == nil {
		return 0
	}
	return e.GetTotalCost(root)
}

// GetFinancialsForPrices restituisce le metriche finanziarie per ogni prezzo in prices.
func (e *AnalysisEngine) GetFinancialsForPrices(root *domain.Activity, prices []float64) []FinancialMetrics {
	out := make([]FinancialMetrics, 0, len(prices))
	for _, p := range prices {
		out = append(out, e.GetFinancials(root, p))
	}
	return out
}
