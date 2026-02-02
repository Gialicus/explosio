package lib

// AnalysisEngine fornisce funzionalità di analisi per progetti strutturati come alberi di attività.
// Supporta:
//   - Calcolo CPM (Critical Path Method) e identificazione del cammino critico
//   - Calcolo costi totali e breakdown per categoria
//   - Analisi finanziaria (margini, markup, break-even)
//   - Analisi fornitori e calcolo requisiti per scenari di produzione
//   - Validazione di attività e fornitori
//   - Ottimizzazione tramite crashing (compressione tempi)
//
// I metodi sono organizzati in moduli separati per migliorare la manutenibilità:
//   - cpm.go: calcolo CPM e cammino critico
//   - cost.go: calcolo costi
//   - financial.go: analisi finanziaria
//   - supplier.go: analisi fornitori
//   - validation.go: validazione
//   - crashing.go: ottimizzazione crashing
type AnalysisEngine struct{}
