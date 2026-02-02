package ui

import (
	"explosio/lib"
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// WhatIfView gestisce l'analisi what-if
type WhatIfView struct {
	appState  *AppState
	window    fyne.Window
	scenarios []lib.Scenario
	results   []lib.ScenarioResult
	table     *widget.Table
}

// NewWhatIfView crea una nuova vista what-if
func NewWhatIfView(appState *AppState) fyne.CanvasObject {
	view := &WhatIfView{
		appState:  appState,
		scenarios: make([]lib.Scenario, 0),
		results:   make([]lib.ScenarioResult, 0),
	}
	
	// Ottieni la window corrente
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		view.window = windows[0]
	}

	// Tabella risultati
	view.table = widget.NewTable(
		func() (int, int) {
			return len(view.results) + 1, 7 // +1 per header
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id.Row == 0 {
				// Header
				headers := []string{"Scenario", "Durata", "Costo", "Margine", "Markup%", "Tempo Risparmiabile", "Costo Extra"}
				if id.Col < len(headers) {
					label.SetText(headers[id.Col])
				}
			} else {
				// Dati
				row := id.Row - 1
				if row < len(view.results) {
					result := view.results[row]
					switch id.Col {
					case 0:
						label.SetText(result.ScenarioName)
					case 1:
						label.SetText(fmt.Sprintf("%d min", result.TotalDuration))
					case 2:
						label.SetText(fmt.Sprintf("€%.2f", result.TotalCost))
					case 3:
						label.SetText(fmt.Sprintf("€%.2f", result.Margin))
					case 4:
						label.SetText(fmt.Sprintf("%.1f%%", result.Markup))
					case 5:
						label.SetText(fmt.Sprintf("%d min", result.TimeSaved))
					case 6:
						label.SetText(fmt.Sprintf("€%.2f", result.ExtraCrashCost))
					}
				}
			}
		},
	)

	// Pulsante per aggiungere scenario
	addBtn := widget.NewButton("Aggiungi Scenario", func() {
		view.addScenario()
	})

	// Pulsante per eseguire analisi
	runBtn := widget.NewButton("Esegui Analisi", func() {
		view.runAnalysis()
	})

	// Layout
	content := container.NewBorder(
		container.NewHBox(addBtn, runBtn),
		nil,
		nil,
		nil,
		container.NewScroll(view.table),
	)

	return content
}

// addScenario aggiunge un nuovo scenario
func (v *WhatIfView) addScenario() {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(fmt.Sprintf("Scenario %d", len(v.scenarios)+1))
	priceEntry := widget.NewEntry()
	priceEntry.SetText("0.00")

	form := widget.NewForm(
		widget.NewFormItem("Nome Scenario", nameEntry),
		widget.NewFormItem("Prezzo di Vendita (€)", priceEntry),
	)

	var dialog *widget.PopUp
	form.OnSubmit = func() {
		var sellPrice float64
		fmt.Sscanf(priceEntry.Text, "%f", &sellPrice)

		scenario := lib.Scenario{
			Name:      nameEntry.Text,
			SellPrice: sellPrice,
			Overrides: make(map[string]lib.ActivityOverride),
		}

		v.scenarios = append(v.scenarios, scenario)
		dialog.Hide()
	}

	form.OnCancel = func() {
		dialog.Hide()
	}

	// Per il dialog, usiamo la window se disponibile
	var canvas fyne.Canvas
	if v.window != nil {
		canvas = v.window.Canvas()
	} else if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}
	
	dialog = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Nuovo Scenario"),
			form,
		),
		canvas,
	)
	dialog.Resize(fyne.NewSize(400, 200))
	dialog.Show()
}

// runAnalysis esegue l'analisi what-if
func (v *WhatIfView) runAnalysis() {
	project := v.appState.GetProject()
	if project == nil || project.Root == nil {
		return
	}

	whatIfEngine := lib.NewWhatIfEngine()
	v.results = whatIfEngine.RunScenarios(project.Root, v.scenarios)
	v.table.Refresh()
}
