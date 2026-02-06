package engine

import (
	"explosio/lib/domain"
	"explosio/lib/tree"
)

// AnalysisEngine fornisce funzionalità di analisi per progetti strutturati come alberi di attività.
// Supporta calcolo CPM, costi, analisi finanziaria, fornitori, validazione e crashing.
type AnalysisEngine struct{}

// Walk attraversa l'albero in pre-order chiamando f su ogni attività.
func (e *AnalysisEngine) Walk(root *domain.Activity, f func(*domain.Activity)) {
	tree.Walk(root, f)
}
