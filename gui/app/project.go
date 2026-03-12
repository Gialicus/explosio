// Package app provides application services for the Explosio GUI.
package app

import (
	"io"
	"strings"

	"explosio/core"
	"explosio/core/unit"
)

// ProjectFormat indicates the file format for project persistence.
type ProjectFormat string

const (
	FormatJSON ProjectFormat = "json"
	FormatYAML ProjectFormat = "yaml"
)

// ProjectService handles project loading, saving, and validation.
type ProjectService struct{}

// NewProjectService creates a new ProjectService.
func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// LoadProject loads a project from the given reader.
// format is "json", "yaml", or "yml".
func (s *ProjectService) LoadProject(r io.Reader, format string) (*core.Project, error) {
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "yaml" || format == "yml" {
		return core.ReadYAML(r)
	}
	return core.ReadJSON(r)
}

// SaveProject saves the project to the given writer in JSON format.
func (s *ProjectService) SaveProject(proj *core.Project, w io.Writer) error {
	return proj.WriteJSON(w)
}

// Validate validates the project and returns errors and warnings.
func (s *ProjectService) Validate(proj *core.Project) *core.ValidationResult {
	if proj == nil || proj.Root == nil {
		return &core.ValidationResult{}
	}
	return proj.Root.Validate()
}

// ProjectSummary holds aggregated project statistics for display.
type ProjectSummary struct {
	TotalPrice   float64
	Currency     string
	TotalDuration float64
	DurationUnit  string
	CostBreakdown core.CostBreakdown
}

// GetSummary returns the project summary for the given root activity.
func (s *ProjectService) GetSummary(root *core.Activity) ProjectSummary {
	if root == nil {
		return ProjectSummary{}
	}
	totalPrice := root.CalculatePrice()
	totalDur := root.CalculateDuration()
	cb := root.CostBreakdown()
	unitStr := string(root.Duration.Unit)
	if unitStr == "" {
		unitStr = "day"
	}
	return ProjectSummary{
		TotalPrice:     totalPrice,
		Currency:       root.Price.Currency,
		TotalDuration:  totalDur,
		DurationUnit:   unitStr,
		CostBreakdown:  cb,
	}
}

// NewProject creates a new empty project with a minimal default root activity.
func (s *ProjectService) NewProject() *core.Project {
	root := core.NewActivity(
		"Nuovo progetto",
		"",
		unit.Duration{Value: 0, Unit: unit.DurationUnitDay},
		unit.Price{Value: 0, Currency: "EUR"},
	)
	return core.NewProject(root)
}
