package ui

import (
	"bytes"
	"explosio/lib"
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ResultsView mostra i risultati dell'analisi
type ResultsView struct {
	appState *AppState
	tabs     *container.AppTabs
}

// NewResultsView crea una nuova vista dei risultati
func NewResultsView(appState *AppState) fyne.CanvasObject {
	view := &ResultsView{
		appState: appState,
	}

	// Crea tabs per diverse visualizzazioni
	view.tabs = container.NewAppTabs(
		container.NewTabItem("Report", view.createReportTab()),
		container.NewTabItem("Metriche Finanziarie", view.createFinancialTab()),
		container.NewTabItem("Requisiti Fornitori", view.createSupplierRequirementsTab()),
	)

	// Aggiorna quando cambia il progetto
	appState.OnProjectChanged(func(*lib.Project) {
		view.refresh()
	})

	return view.tabs
}

// refresh aggiorna tutti i tab
func (v *ResultsView) refresh() {
	if len(v.tabs.Items) >= 3 {
		v.tabs.Items[0].Content = v.createReportTab()
		v.tabs.Items[1].Content = v.createFinancialTab()
		v.tabs.Items[2].Content = v.createSupplierRequirementsTab()
		v.tabs.Refresh()
	}
}

// createReportTab crea il tab con il report gerarchico
func (v *ResultsView) createReportTab() fyne.CanvasObject {
	project := v.appState.GetProject()
	if project == nil || project.Root == nil {
		return widget.NewLabel("Nessun progetto caricato")
	}

	// Calcola CPM se non già fatto
	engine := v.appState.GetEngine()
	engine.ComputeCPM(project.Root)

	// Genera il report
	var buf bytes.Buffer
	lib.PrintReportTo(&buf, project.Root, 0, true, "")

	reportText := widget.NewRichTextFromMarkdown("```\n" + buf.String() + "\n```")
	scroll := container.NewScroll(reportText)

	return scroll
}

// createFinancialTab crea il tab con le metriche finanziarie
func (v *ResultsView) createFinancialTab() fyne.CanvasObject {
	project := v.appState.GetProject()
	if project == nil || project.Root == nil {
		return widget.NewLabel("Nessun progetto caricato")
	}

	engine := v.appState.GetEngine()
	engine.ComputeCPM(project.Root)

	// Calcola costi
	totalCost := engine.GetTotalCost(project.Root)
	
	// Form per inserire prezzo di vendita
	priceEntry := widget.NewEntry()
	priceEntry.SetText("0.00")
	
	metricsCard := widget.NewCard("Metriche", "", nil)
	
	updateMetrics := func() {
		var sellPrice float64
		fmt.Sscanf(priceEntry.Text, "%f", &sellPrice)
		
		fin := engine.GetFinancials(project.Root, sellPrice)
		
		metricsText := fmt.Sprintf(
			"Costo di Produzione: €%.2f\n"+
			"Prezzo di Vendita: €%.2f\n"+
			"Margine Unitario: €%.2f\n"+
			"Markup: %.1f%%\n"+
			"Fattibilità: %s",
			fin.TotalCost,
			sellPrice,
			fin.Margin,
			fin.Markup,
			map[bool]string{true: "Sì", false: "No"}[fin.IsViable],
		)
		
		metricsCard.SetContent(widget.NewRichTextFromMarkdown(metricsText))
	}
	
	priceEntry.OnChanged = func(string) {
		updateMetrics()
	}
	
	updateMetrics()

	// Analisi crashing
	timeSaved, extraCost := engine.GetMaxCrashPotential(project.Root)
	crashingText := fmt.Sprintf(
		"Tempo Risparmiabile: %d min\n"+
		"Costo Extra Stimato: €%.2f\n"+
		"Costo Nuovo (Accelerato): €%.2f",
		timeSaved,
		extraCost,
		totalCost+extraCost,
	)
	crashingCard := widget.NewCard("Ottimizzazione (Crashing)", "", widget.NewRichTextFromMarkdown(crashingText))

	content := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Prezzo di Vendita (€)", priceEntry),
		),
		metricsCard,
		crashingCard,
	)

	return container.NewScroll(content)
}

// createSupplierRequirementsTab crea il tab con i requisiti dei fornitori
func (v *ResultsView) createSupplierRequirementsTab() fyne.CanvasObject {
	project := v.appState.GetProject()
	if project == nil || project.Root == nil {
		return widget.NewLabel("Nessun progetto caricato")
	}

	engine := v.appState.GetEngine()
	engine.ComputeCPM(project.Root)

	// Form per inserire target di produzione
	productionEntry := widget.NewEntry()
	productionEntry.SetText("1000")
	
	periodSelect := widget.NewSelect([]string{
		string(lib.PeriodMinute),
		string(lib.PeriodHour),
		string(lib.PeriodDay),
		string(lib.PeriodWeek),
		string(lib.PeriodMonth),
		string(lib.PeriodYear),
	}, nil)
	periodSelect.SetSelected(string(lib.PeriodDay))

	requirementsText := widget.NewRichTextFromMarkdown("")
	
	updateRequirements := func() {
		var productionTarget float64
		fmt.Sscanf(productionEntry.Text, "%f", &productionTarget)
		
		if productionTarget <= 0 {
			requirementsText.ParseMarkdown("Inserisci un target di produzione valido")
			return
		}

		period := lib.PeriodType(periodSelect.Selected)
		if !period.IsValid() {
			return
		}

		requirements := engine.CalculateSupplierRequirements(project.Root, productionTarget, period)
		
		if len(requirements) == 0 {
			requirementsText.ParseMarkdown("Nessun requisito di fornitore trovato.")
			return
		}

		var buf bytes.Buffer
		buf.WriteString("### Requisiti Fornitori\n\n")
		buf.WriteString("| Fornitore | Quantità Richiesta | Periodo | Fornitori Necessari | Fattibile |\n")
		buf.WriteString("|-----------|-------------------|---------|-------------------|----------|\n")
		
		for _, req := range requirements {
			feasible := "Sì"
			if !req.IsFeasible {
				feasible = "No"
			}
			buf.WriteString(fmt.Sprintf("| %s | %.2f | %s | %.2f | %s |\n",
				req.SupplierName,
				req.RequiredQuantity,
				req.SupplierPeriod.String(),
				req.SuppliersNeeded,
				feasible,
			))
		}
		
		requirementsText.ParseMarkdown(buf.String())
	}
	
	productionEntry.OnChanged = func(string) {
		updateRequirements()
	}
	
	periodSelect.OnChanged = func(string) {
		updateRequirements()
	}
	
	updateRequirements()

	content := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Target Produzione", productionEntry),
			widget.NewFormItem("Periodo", periodSelect),
		),
		container.NewScroll(requirementsText),
	)

	return content
}
