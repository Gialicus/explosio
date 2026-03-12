package components

import (
	"fmt"

	"explosio/gui/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StatusBar displays project summary (total price, duration, cost breakdown).
type StatusBar struct {
	label *widget.Label
}

// NewStatusBar creates a new StatusBar.
func NewStatusBar() *StatusBar {
	return &StatusBar{
		label: widget.NewLabel(""),
	}
}

// Update updates the displayed summary.
func (s *StatusBar) Update(summary app.ProjectSummary) {
	cb := summary.CostBreakdown
	s.label.SetText(fmt.Sprintf(
		"Totale: %.0f %s | Durata: %.0f %s | Attività: %.0f | Materiali: %.0f | Risorse: %.0f | Asset: %.0f",
		summary.TotalPrice, summary.Currency,
		summary.TotalDuration, summary.DurationUnit,
		cb.Activities, cb.Materials, cb.Human, cb.Assets,
	))
}

// Content returns the status bar container for layout.
func (s *StatusBar) Content() fyne.CanvasObject {
	return container.NewVBox(widget.NewSeparator(), s.label)
}
