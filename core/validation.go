// Package core provides validation for activity trees.
package core

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation issue.
type ValidationError struct {
	Activity string
	Message  string
}

func (e ValidationError) Error() string {
	if e.Activity != "" {
		return fmt.Sprintf("%s: %s", e.Activity, e.Message)
	}
	return e.Message
}

// ValidationResult holds validation errors and warnings.
type ValidationResult struct {
	Errors   []ValidationError
	Warnings []ValidationError
}

// Valid returns true if there are no errors.
func (r *ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

// AddError adds an error.
func (r *ValidationResult) AddError(activity, msg string) {
	r.Errors = append(r.Errors, ValidationError{Activity: activity, Message: msg})
}

// AddWarning adds a warning.
func (r *ValidationResult) AddWarning(activity, msg string) {
	r.Warnings = append(r.Warnings, ValidationError{Activity: activity, Message: msg})
}

// Validate checks the activity tree for errors and warnings.
func (a *Activity) Validate() *ValidationResult {
	r := &ValidationResult{}
	all := make(map[*Activity]bool)
	collectActivities(a, all)

	// Check for circular dependencies
	if hasCycleInTree(all) {
		r.AddError(a.Name, "circular dependency detected in DependsOn")
	}

	// Check that DependsOn references exist in the tree
	for _, act := range allActivities(a) {
		for _, dep := range act.DependsOn {
			if !all[dep] {
				r.AddError(act.Name, fmt.Sprintf("DependsOn references activity %q not in tree", dep.Name))
			}
		}
	}

	// Warnings: activities without materials/resources
	for _, act := range allActivities(a) {
		if len(act.ComplexMaterials) == 0 && len(act.CountableMaterials) == 0 && len(act.MeasurableMaterials) == 0 &&
			len(act.HumanResources) == 0 && len(act.Assets) == 0 && len(act.Activities) > 0 {
			r.AddWarning(act.Name, "activity has sub-activities but no materials or resources")
		}
	}

	// Check currency consistency (all same currency)
	currencies := make(map[string]bool)
	for _, act := range allActivities(a) {
		if act.Price.Currency != "" {
			currencies[act.Price.Currency] = true
		}
	}
	if len(currencies) > 1 {
		var list []string
		for c := range currencies {
			list = append(list, c)
		}
		r.AddWarning(a.Name, "multiple currencies used: "+strings.Join(list, ", "))
	}

	return r
}

func collectActivities(a *Activity, m map[*Activity]bool) {
	m[a] = true
	for _, child := range a.Activities {
		collectActivities(child, m)
	}
}

func allActivities(a *Activity) []*Activity {
	return a.GetActivities()
}

// hasCycle performs DFS to detect cycles in the DependsOn graph.
func hasCycle(a *Activity, all map[*Activity]bool, visiting, visited map[*Activity]bool) bool {
	if visiting[a] {
		return true
	}
	if visited[a] {
		return false
	}
	visiting[a] = true
	for _, dep := range a.DependsOn {
		if all[dep] {
			if hasCycle(dep, all, visiting, visited) {
				return true
			}
		}
	}
	visiting[a] = false
	visited[a] = true
	return false
}

// hasCycleInTree checks if any activity in the tree has a cyclic DependsOn.
func hasCycleInTree(all map[*Activity]bool) bool {
	visiting := make(map[*Activity]bool)
	visited := make(map[*Activity]bool)
	for act := range all {
		if !visited[act] && hasCycle(act, all, visiting, visited) {
			return true
		}
	}
	return false
}
