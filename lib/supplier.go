// Package lib fornisce funzionalità per l'analisi di progetti strutturati come alberi di attività.
// Questo modulo gestisce l'analisi dei fornitori e il calcolo dei requisiti per scenari di produzione.
package lib

import "fmt"

// ProductionRequirement rappresenta un requisito di produzione (quantità per periodo)
type ProductionRequirement struct {
	ResourceName string
	Quantity     float64
	Period       PeriodType
}

// SupplierRequirement rappresenta il requisito di fornitori per soddisfare una produzione
type SupplierRequirement struct {
	SupplierName     string
	RequiredQuantity float64
	SupplierPeriod   PeriodType
	SuppliersNeeded  float64 // numero di fornitori (può essere frazionario)
	IsFeasible       bool    // true se la capacità è sufficiente
}

// CalculateSupplierRequirements calcola quanti fornitori sono necessari per soddisfare una produzione target
func (e *AnalysisEngine) CalculateSupplierRequirements(root *Activity, productionTarget float64, targetPeriod PeriodType) []SupplierRequirement {
	if root == nil {
		return nil
	}

	// Mappa per raggruppare i requisiti per fornitore
	supplierMap := make(map[string]*SupplierRequirement)
	// Mappa per salvare i fornitori originali
	supplierRefs := make(map[string]*Supplier)

	// Visita ricorsiva per raccogliere tutti i requisiti
	e.collectSupplierRequirementsRec(root, productionTarget, targetPeriod, supplierMap, supplierRefs)

	// Calcola il numero di fornitori necessari per ogni requisito
	for _, req := range supplierMap {
		if supplier, exists := supplierRefs[req.SupplierName]; exists {
			req.SuppliersNeeded = req.RequiredQuantity / supplier.AvailableQuantity
			req.IsFeasible = supplier.AvailableQuantity > 0
		} else {
			req.IsFeasible = false
		}
	}

	// Converti la mappa in slice
	requirements := make([]SupplierRequirement, 0, len(supplierMap))
	for _, req := range supplierMap {
		requirements = append(requirements, *req)
	}

	return requirements
}

// collectSupplierRequirementsRec raccoglie ricorsivamente i requisiti di fornitori
func (e *AnalysisEngine) collectSupplierRequirementsRec(a *Activity, productionTarget float64, targetPeriod PeriodType, supplierMap map[string]*SupplierRequirement, supplierRefs map[string]*Supplier) {
	if a == nil {
		return
	}

	// Processa materiali con fornitori
	for _, m := range a.Materials {
		if m.Supplier != nil {
			// Salva il riferimento al fornitore
			supplierRefs[m.Supplier.Name] = m.Supplier

			// Quantità per unità di prodotto
			quantityPerUnit := m.Quantity
			// Quantità totale necessaria per la produzione target
			totalQuantity := quantityPerUnit * productionTarget

			// Converti la quantità totale nel periodo del fornitore
			targetMinutes := targetPeriod.ToMinutes()
			supplierMinutes := m.Supplier.Period.ToMinutes()
			if targetMinutes > 0 && supplierMinutes > 0 {
				// Calcola la quantità nel periodo del fornitore
				quantityInSupplierPeriod := (totalQuantity / float64(targetMinutes)) * float64(supplierMinutes)

				// Aggiorna o crea il requisito per questo fornitore
				if req, exists := supplierMap[m.Supplier.Name]; exists {
					req.RequiredQuantity += quantityInSupplierPeriod
				} else {
					supplierMap[m.Supplier.Name] = &SupplierRequirement{
						SupplierName:     m.Supplier.Name,
						RequiredQuantity: quantityInSupplierPeriod,
						SupplierPeriod:   m.Supplier.Period,
					}
				}
			}
		}
	}

	// Processa risorse umane con fornitori
	for _, h := range a.Humans {
		if h.Supplier != nil {
			// Salva il riferimento al fornitore
			supplierRefs[h.Supplier.Name] = h.Supplier

			// Per le risorse umane, la quantità è già per unità di prodotto
			quantityPerUnit := h.Quantity
			totalQuantity := quantityPerUnit * productionTarget

			targetMinutes := targetPeriod.ToMinutes()
			supplierMinutes := h.Supplier.Period.ToMinutes()
			if targetMinutes > 0 && supplierMinutes > 0 {
				quantityInSupplierPeriod := (totalQuantity / float64(targetMinutes)) * float64(supplierMinutes)

				if req, exists := supplierMap[h.Supplier.Name]; exists {
					req.RequiredQuantity += quantityInSupplierPeriod
				} else {
					supplierMap[h.Supplier.Name] = &SupplierRequirement{
						SupplierName:     h.Supplier.Name,
						RequiredQuantity: quantityInSupplierPeriod,
						SupplierPeriod:   h.Supplier.Period,
					}
				}
			}
		}
	}

	// Processa asset con fornitori
	for _, as := range a.Assets {
		if as.Supplier != nil {
			// Salva il riferimento al fornitore
			supplierRefs[as.Supplier.Name] = as.Supplier

			quantityPerUnit := as.Quantity
			totalQuantity := quantityPerUnit * productionTarget

			targetMinutes := targetPeriod.ToMinutes()
			supplierMinutes := as.Supplier.Period.ToMinutes()
			if targetMinutes > 0 && supplierMinutes > 0 {
				quantityInSupplierPeriod := (totalQuantity / float64(targetMinutes)) * float64(supplierMinutes)

				if req, exists := supplierMap[as.Supplier.Name]; exists {
					req.RequiredQuantity += quantityInSupplierPeriod
				} else {
					supplierMap[as.Supplier.Name] = &SupplierRequirement{
						SupplierName:     as.Supplier.Name,
						RequiredQuantity: quantityInSupplierPeriod,
						SupplierPeriod:   as.Supplier.Period,
					}
				}
			}
		}
	}

	// Visita ricorsivamente le sotto-attività
	for _, sub := range a.SubActivities {
		e.collectSupplierRequirementsRec(sub, productionTarget, targetPeriod, supplierMap, supplierRefs)
	}
}

// ValidateSupplierUsage valida che le quantità utilizzate non superino la capacità disponibile dei fornitori.
// Restituisce un ValidationErrors se ci sono errori, nil altrimenti.
func (e *AnalysisEngine) ValidateSupplierUsage(root *Activity) error {
	var ve ValidationErrors
	e.validateSupplierUsageRec(root, &ve)
	if ve.HasErrors() {
		return &ve
	}
	return nil
}

func (e *AnalysisEngine) validateSupplierUsageRec(a *Activity, ve *ValidationErrors) {
	if a == nil {
		return
	}
	for i, h := range a.Humans {
		if h.Supplier != nil && h.Quantity > h.Supplier.AvailableQuantity {
			ve.Add(fmt.Errorf(
				"activity %s: human resource %d (%s) uses %.1f units but supplier %s only provides %.1f/%s",
				a.ID, i, h.Role, h.Quantity, h.Supplier.Name,
				h.Supplier.AvailableQuantity, h.Supplier.Period.String()))
		}
	}
	for i, m := range a.Materials {
		if m.Supplier != nil && m.Quantity > m.Supplier.AvailableQuantity {
			ve.Add(fmt.Errorf(
				"activity %s: material resource %d (%s) uses %.1f units but supplier %s only provides %.1f/%s",
				a.ID, i, m.Name, m.Quantity, m.Supplier.Name,
				m.Supplier.AvailableQuantity, m.Supplier.Period.String()))
		}
	}
	for i, as := range a.Assets {
		if as.Supplier != nil && as.Quantity > as.Supplier.AvailableQuantity {
			ve.Add(fmt.Errorf(
				"activity %s: asset resource %d (%s) uses %.1f units but supplier %s only provides %.1f/%s",
				a.ID, i, as.Name, as.Quantity, as.Supplier.Name,
				as.Supplier.AvailableQuantity, as.Supplier.Period.String()))
		}
	}
	for _, sub := range a.SubActivities {
		e.validateSupplierUsageRec(sub, ve)
	}
}
