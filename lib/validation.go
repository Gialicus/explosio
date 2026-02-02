package lib

import "fmt"

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
