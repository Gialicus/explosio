package engine

import (
	"explosio/lib/domain"
	"explosio/lib/resources"
	"fmt"
)

// ProductionRequirement rappresenta un requisito di produzione (quantità per periodo)
type ProductionRequirement struct {
	ResourceName string
	Quantity     float64
	Period       domain.PeriodType
}

// SupplierRequirement rappresenta il requisito di fornitori per soddisfare una produzione
type SupplierRequirement struct {
	SupplierName     string
	RequiredQuantity float64
	SupplierPeriod   domain.PeriodType
	SuppliersNeeded  float64
	IsFeasible       bool
}

// CalculateSupplierRequirements calcola quanti fornitori sono necessari per soddisfare una produzione target
func (e *AnalysisEngine) CalculateSupplierRequirements(root *domain.Activity, productionTarget float64, targetPeriod domain.PeriodType) []SupplierRequirement {
	if root == nil {
		return nil
	}
	supplierMap := make(map[string]*SupplierRequirement)
	supplierRefs := make(map[string]*domain.Supplier)
	targetMinutes := targetPeriod.ToMinutes()
	resources.WalkResources(root, func(a *domain.Activity, r domain.Resource) {
		s := r.GetSupplier()
		if s == nil {
			return
		}
		supplierRefs[s.Name] = s
		quantityPerUnit := r.GetQuantity()
		totalQuantity := quantityPerUnit * productionTarget
		supplierMinutes := s.Period.ToMinutes()
		if targetMinutes <= 0 || supplierMinutes <= 0 {
			return
		}
		quantityInSupplierPeriod := (totalQuantity / float64(targetMinutes)) * float64(supplierMinutes)
		if req, exists := supplierMap[s.Name]; exists {
			req.RequiredQuantity += quantityInSupplierPeriod
		} else {
			supplierMap[s.Name] = &SupplierRequirement{
				SupplierName:     s.Name,
				RequiredQuantity: quantityInSupplierPeriod,
				SupplierPeriod:   s.Period,
			}
		}
	})
	for _, req := range supplierMap {
		if supplier, exists := supplierRefs[req.SupplierName]; exists {
			req.SuppliersNeeded = req.RequiredQuantity / supplier.AvailableQuantity
			req.IsFeasible = supplier.AvailableQuantity > 0
		} else {
			req.IsFeasible = false
		}
	}
	requirements := make([]SupplierRequirement, 0, len(supplierMap))
	for _, req := range supplierMap {
		requirements = append(requirements, *req)
	}
	return requirements
}

// ValidateSupplierUsage valida che le quantità utilizzate non superino la capacità disponibile dei fornitori.
func (e *AnalysisEngine) ValidateSupplierUsage(root *domain.Activity) error {
	var ve domain.ValidationErrors
	resources.WalkResources(root, func(a *domain.Activity, r domain.Resource) {
		s := r.GetSupplier()
		if s == nil || r.GetQuantity() <= s.AvailableQuantity {
			return
		}
		ve.Add(fmt.Errorf(
			"activity %s: resource %s uses %.1f units but supplier %s only provides %.1f/%s",
			a.ID, resources.ResourceDisplayName(r), r.GetQuantity(), s.Name,
			s.AvailableQuantity, s.Period.String()))
	})
	if ve.HasErrors() {
		return &ve
	}
	return nil
}
