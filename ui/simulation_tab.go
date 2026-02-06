package ui

import (
	"explosio/lib"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// scenarioRow rappresenta uno scenario definito dall'utente (nome, prezzo, override).
type scenarioRow struct {
	Name      string
	SellPrice float64
	Overrides map[string]lib.ActivityOverride
}

// NewSimulationTab restituisce il contenuto del tab Simulazione: stato vuoto o due card (Scenari, Risultati).
func NewSimulationTab(project *lib.Project) fyne.CanvasObject {
	if project == nil || project.Root == nil {
		return container.NewCenter(
			widget.NewLabel("Definisci un progetto nel tab Progetto."),
		)
	}

	var scenarios []scenarioRow
	var results []lib.ScenarioResult

	cardScenari := widget.NewCard("Scenari what-if", "", nil)
	cardRisultati := widget.NewCard("Risultati", "", nil)

	var refreshScenari, refreshRisultati func()
	refreshScenari = func() {
		content := container.NewVBox()
		for i := range scenarios {
			s := &scenarios[i]
			row := container.NewBorder(nil, nil,
				widget.NewLabel(fmt.Sprintf("%s — %.2f €", s.Name, s.SellPrice)),
				widget.NewButton("Rimuovi", func() {
					for j := range scenarios {
						if &scenarios[j] == s {
							scenarios = append(scenarios[:j], scenarios[j+1:]...)
							break
						}
					}
					refreshScenari()
				}),
			)
			content.Add(row)
		}
		runBtn := widget.NewButton("Esegui simulazione", func() {
			if len(scenarios) == 0 {
				return
			}
			libScenarios := make([]lib.Scenario, 0, len(scenarios))
			for _, s := range scenarios {
				libScenarios = append(libScenarios, lib.Scenario{
					Name:      s.Name,
					SellPrice: s.SellPrice,
					Overrides: s.Overrides,
				})
			}
			engine := lib.NewWhatIfEngine()
			results = engine.RunScenarios(project.Root, libScenarios)
			refreshRisultati()
		})
		if len(scenarios) == 0 {
			runBtn.Disable()
		} else {
			runBtn.Enable()
		}
		content.Add(container.NewHBox(
			widget.NewButton("Aggiungi scenario", func() { openAddScenarioDialog(project, &scenarios, refreshScenari) }),
			runBtn,
		))
		cardScenari.SetContent(content)
	}
	refreshRisultati = func() {
		if len(results) == 0 {
			cardRisultati.SetContent(container.NewCenter(
				widget.NewLabel("Esegui la simulazione per vedere i risultati."),
			))
			return
		}
		rows := container.NewVBox()
		// Intestazione
		header := container.NewHBox(
			widget.NewLabel("Scenario"),
			widget.NewLabel("Durata"),
			widget.NewLabel("Costo"),
			widget.NewLabel("Margine"),
			widget.NewLabel("Markup %"),
		)
		rows.Add(header)
		for _, r := range results {
			row := container.NewHBox(
				widget.NewLabel(r.ScenarioName),
				widget.NewLabel(fmt.Sprintf("%d min", r.TotalDuration)),
				widget.NewLabel(fmt.Sprintf("%.2f €", r.TotalCost)),
				widget.NewLabel(fmt.Sprintf("%.2f €", r.Margin)),
				widget.NewLabel(fmt.Sprintf("%.1f%%", r.Markup)),
			)
			rows.Add(row)
		}
		cardRisultati.SetContent(container.NewPadded(container.NewScroll(rows)))
	}

	refreshScenari()
	refreshRisultati()

	// Split: Scenari in alto (scrollabile), Risultati in basso che occupa tutto lo spazio verticale restante
	scenariScroll := container.NewScroll(cardScenari)
	scenariScroll.Direction = container.ScrollVerticalOnly
	split := container.NewVSplit(scenariScroll, cardRisultati)
	split.SetOffset(0.35) // 35% scenari, 65% tabella risultati (sempre visibile)
	return container.NewPadded(split)
}

// openAddScenarioDialog apre il popup per aggiungere uno scenario (nome, prezzo, override per attività).
func openAddScenarioDialog(project *lib.Project, scenarios *[]scenarioRow, refreshScenari func()) {
	win := getWindow()
	if win == nil {
		return
	}
	activityMap := BuildActivityMap(project)
	options := make([]string, 0, len(activityMap))
	optionToID := make(map[string]string)
	for id, a := range activityMap {
		label := fmt.Sprintf("%s - %s", id, a.Name)
		options = append(options, label)
		optionToID[label] = id
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nome scenario")
	priceEntry := widget.NewEntry()
	priceEntry.SetText("0")
	priceEntry.SetPlaceHolder("Prezzo vendita (es. 4.50)")

	var overrideRowsContainer *fyne.Container
	overrideRowsContainer = container.NewVBox()

	addOverrideRow := func() {
		if len(options) == 0 {
			return
		}
		sel := widget.NewSelect(options, nil)
		sel.SetSelected(options[0])
		humanEntry := widget.NewEntry()
		humanEntry.SetText("1.0")
		materialEntry := widget.NewEntry()
		materialEntry.SetText("1.0")
		assetEntry := widget.NewEntry()
		assetEntry.SetText("1.0")
		row := container.NewHBox(
			widget.NewLabel("Attività"),
			sel,
			widget.NewLabel("Umano"),
			humanEntry,
			widget.NewLabel("Materiali"),
			materialEntry,
			widget.NewLabel("Asset"),
			assetEntry,
			widget.NewButton("Rimuovi", func() {
			}),
		)
		removeBtn := row.Objects[len(row.Objects)-1].(*widget.Button)
		removeBtn.OnTapped = func() {
			overrideRowsContainer.Remove(row)
			overrideRowsContainer.Refresh()
		}
		overrideRowsContainer.Add(row)
		overrideRowsContainer.Refresh()
	}

	addOverrideBtn := widget.NewButton("Aggiungi override", addOverrideRow)

	dialogContent := container.NewVBox(
		widget.NewLabel("Nome"),
		nameEntry,
		widget.NewLabel("Prezzo di vendita (€)"),
		priceEntry,
		widget.NewLabel("Override per attività"),
		addOverrideBtn,
		overrideRowsContainer,
	)

	var pop *widget.PopUp
	submit := func() {
		name := nameEntry.Text
		if name == "" {
			return
		}
		var price float64
		if _, err := fmt.Sscanf(priceEntry.Text, "%f", &price); err != nil {
			return
		}
		overrides := make(map[string]lib.ActivityOverride)
		for _, obj := range overrideRowsContainer.Objects {
			box, ok := obj.(*fyne.Container)
			if !ok || len(box.Objects) < 9 {
				continue
			}
			sel, ok := box.Objects[1].(*widget.Select)
			if !ok {
				continue
			}
			selected := sel.Selected
			if selected == "" {
				continue
			}
			actID := optionToID[selected]
			humanEntry := box.Objects[3].(*widget.Entry)
			materialEntry := box.Objects[5].(*widget.Entry)
			assetEntry := box.Objects[7].(*widget.Entry)
			var h, m, a float64
			fmt.Sscanf(humanEntry.Text, "%f", &h)
			fmt.Sscanf(materialEntry.Text, "%f", &m)
			fmt.Sscanf(assetEntry.Text, "%f", &a)
			if h == 0 {
				h = 1.0
			}
			if m == 0 {
				m = 1.0
			}
			if a == 0 {
				a = 1.0
			}
			if h != 1.0 || m != 1.0 || a != 1.0 {
				overrides[actID] = lib.ActivityOverride{
					HumanCostFactor:    h,
					MaterialCostFactor: m,
					AssetCostFactor:   a,
				}
			}
		}
		var ov map[string]lib.ActivityOverride
		if len(overrides) > 0 {
			ov = overrides
		}
		*scenarios = append(*scenarios, scenarioRow{Name: name, SellPrice: price, Overrides: ov})
		refreshScenari()
		if pop != nil {
			pop.Hide()
		}
	}
	cancel := func() {
		if pop != nil {
			pop.Hide()
		}
	}

	buttons := container.NewHBox(
		widget.NewButton("Annulla", cancel),
		widget.NewButton("Aggiungi", submit),
	)
	dialogContent.Add(buttons)
	pop = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Aggiungi scenario"),
			dialogContent,
		),
		win.Canvas(),
	)
	pop.Resize(fyne.NewSize(520, 400))
	pop.Show()
}

func getWindow() fyne.Window {
	w := fyne.CurrentApp().Driver().AllWindows()
	if len(w) > 0 {
		return w[0]
	}
	return nil
}
