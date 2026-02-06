package engine

import (
	"explosio/lib/domain"
	"explosio/lib/resources"
	"fmt"
)

// Validate controlla che l'albero di attività sia valido (durata, MinDuration, assenza di cicli, risorse e fornitori).
func (e *AnalysisEngine) Validate(root *domain.Activity) error {
	if root == nil {
		return fmt.Errorf("root activity is nil")
	}
	seen := make(map[string]bool)
	if err := e.validateRec(root, seen); err != nil {
		return err
	}
	var ve domain.ValidationErrors
	resources.WalkResources(root, func(a *domain.Activity, r domain.Resource) {
		if err := validateResource(r); err != nil {
			ve.Add(fmt.Errorf("activity %s: %w", a.ID, err))
		}
	})
	if ve.HasErrors() {
		return &ve
	}
	return nil
}

func (e *AnalysisEngine) validateRec(a *domain.Activity, seen map[string]bool) error {
	if a == nil {
		return nil
	}
	if seen[a.ID] {
		return fmt.Errorf("activity %s: cycle detected", a.ID)
	}
	seen[a.ID] = true
	if err := a.ValidateBasic(); err != nil {
		return err
	}
	for _, sub := range a.SubActivities {
		if err := e.validateRec(sub, seen); err != nil {
			return err
		}
	}
	return nil
}

func validateResource(r domain.Resource) error {
	switch x := r.(type) {
	case domain.HumanResource:
		return x.Validate()
	case domain.MaterialResource:
		return x.Validate()
	case domain.Asset:
		return x.Validate()
	default:
		return nil
	}
}
