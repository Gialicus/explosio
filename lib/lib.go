// Package lib è la facade per explosio: re-esporta tipi e funzioni dai sottomoduli
// (domain, tree, resources, engine, clone, whatif, serialize, mermaid, presenter)
// così che ui e cmd possano importare solo explosio/lib.
package lib

import (
	"explosio/lib/domain"
	"explosio/lib/engine"
	"explosio/lib/mermaid"
	"explosio/lib/presenter"
	"explosio/lib/serialize"
	"explosio/lib/whatif"
)

// Tipi di dominio (alias)
type (
	Activity         = domain.Activity
	Project          = domain.Project
	Resource         = domain.Resource
	HumanResource    = domain.HumanResource
	MaterialResource = domain.MaterialResource
	Asset            = domain.Asset
	ValidationErrors = domain.ValidationErrors
)

// Costanti temporali
const (
	MinutesPerHour  = domain.MinutesPerHour
	MinutesPerDay   = domain.MinutesPerDay
	MinutesPerWeek  = domain.MinutesPerWeek
	MinutesPerMonth = domain.MinutesPerMonth
	MinutesPerYear  = domain.MinutesPerYear
	DaysPerMonth    = domain.DaysPerMonth
	DaysPerYear     = domain.DaysPerYear
	HoursPerDay     = domain.HoursPerDay
)

// Errori
var (
	ErrInvalidPeriod   = domain.ErrInvalidPeriod
	ErrNegativeQuantity = domain.ErrNegativeQuantity
	ErrNegativeCost    = domain.ErrNegativeCost
	ErrInvalidActivity = domain.ErrInvalidActivity
)

// Costruttori dominio
var NewProject = domain.NewProject

// Engine e tipi di output
type (
	AnalysisEngine   = engine.AnalysisEngine
	CostBreakdown    = engine.CostBreakdown
	FinancialMetrics = engine.FinancialMetrics
	CPMSummary       = engine.CPMSummary
)

// Whatif
type (
	WhatIfEngine    = whatif.WhatIfEngine
	Scenario        = whatif.Scenario
	ScenarioResult  = whatif.ScenarioResult
	ActivityOverride = whatif.ActivityOverride
)

var NewWhatIfEngine = whatif.NewWhatIfEngine

// Serialize
func SerializeProject(project *Project) ([]byte, error) {
	return serialize.SerializeProject(project)
}

func DeserializeProject(data []byte) (*Project, error) {
	return serialize.DeserializeProject(data)
}

// Report e Mermaid
func PrintReport(a *Activity, level int, isLast bool, prefix string) {
	presenter.PrintReport(a, level, isLast, prefix)
}

func PrintReportTo(w interface {
	Write([]byte) (int, error)
}, a *Activity, level int, isLast bool, prefix string) {
	presenter.PrintReportTo(w, a, level, isLast, prefix)
}

func PrintMermaid(root *Activity) {
	mermaid.PrintMermaid(root)
}

func PrintMermaidTo(w interface {
	Write([]byte) (int, error)
}, root *Activity) {
	mermaid.PrintMermaidTo(w, root)
}

func GenerateMermaid(root *Activity) string {
	return mermaid.GenerateMermaid(root)
}

func WriteMermaidToFile(root *Activity, path string) error {
	return mermaid.WriteMermaidToFile(root, path)
}
