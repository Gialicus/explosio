package ui

import (
	"explosio/lib"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// WindowSize definisce la dimensione predefinita della finestra
var WindowSize = fyne.NewSize(1200, 800)

// ExplosioTheme rimosso - causava ricorsione infinita
// Usa il tema di default di Fyne

// Dashboard rappresenta la schermata principale dell'applicazione
type Dashboard struct {
	window   fyne.Window
	appState *AppState
	content  *container.AppTabs
	sidebar  *widget.List
}

// NewDashboard crea una nuova dashboard
func NewDashboard(window fyne.Window, appState *AppState) fyne.CanvasObject {
	d := &Dashboard{
		window:   window,
		appState: appState,
	}

	// Crea la barra degli strumenti superiore
	toolbar := d.createToolbar()

	// Crea il contenuto principale con tabs
	d.content = container.NewAppTabs(
		container.NewTabItem("Progetto", d.createProjectTab()),
		container.NewTabItem("Fornitori", d.createSuppliersTab()),
		container.NewTabItem("Risultati", d.createResultsTab()),
		container.NewTabItem("What-If", d.createWhatIfTab()),
		container.NewTabItem("Diagramma", d.createMermaidTab()),
	)

	// Layout principale: toolbar in alto, contenuto al centro
	mainContent := container.NewBorder(toolbar, nil, nil, nil, d.content)

	return container.NewBorder(nil, nil, nil, nil, mainContent)
}

// createToolbar crea la barra degli strumenti superiore
func (d *Dashboard) createToolbar() fyne.CanvasObject {
	newBtn := widget.NewButton("Nuovo", func() {
		d.appState.SetProject(lib.NewProject())
		d.refreshAllTabs()
	})

	openBtn := widget.NewButton("Apri", func() {
		d.openProject()
	})

	saveBtn := widget.NewButton("Salva", func() {
		d.saveProject()
	})

	computeBtn := widget.NewButton("Calcola CPM", func() {
		d.appState.ComputeCPM()
		d.refreshAllTabs()
	})

	toolbar := container.NewHBox(
		newBtn,
		openBtn,
		saveBtn,
		widget.NewSeparator(),
		computeBtn,
	)

	return toolbar
}

// createProjectTab crea il tab per l'editing del progetto
func (d *Dashboard) createProjectTab() fyne.CanvasObject {
	return NewProjectEditor(d.appState, func() {
		d.refreshAllTabs()
	})
}

// createSuppliersTab crea il tab per la gestione dei fornitori
func (d *Dashboard) createSuppliersTab() fyne.CanvasObject {
	return NewSupplierManager(d.appState, func() {
		d.refreshAllTabs()
	})
}

// createResultsTab crea il tab per i risultati
func (d *Dashboard) createResultsTab() fyne.CanvasObject {
	return NewResultsView(d.appState)
}

// createWhatIfTab crea il tab per l'analisi what-if
func (d *Dashboard) createWhatIfTab() fyne.CanvasObject {
	return NewWhatIfView(d.appState)
}

// createMermaidTab crea il tab per il diagramma Mermaid
func (d *Dashboard) createMermaidTab() fyne.CanvasObject {
	return NewMermaidView(d.appState, d.window)
}

// refreshAllTabs aggiorna tutti i tab
func (d *Dashboard) refreshAllTabs() {
	// Ricrea i tab con contenuto aggiornato - usa SelectTabIndex per cambiare tab attivo
	// e aggiorna il contenuto dei tab esistenti
	if len(d.content.Items) >= 5 {
		d.content.Items[0].Content = d.createProjectTab()
		d.content.Items[1].Content = d.createSuppliersTab()
		d.content.Items[2].Content = d.createResultsTab()
		d.content.Items[3].Content = d.createWhatIfTab()
		d.content.Items[4].Content = d.createMermaidTab()
		d.content.Refresh()
	}
}

// openProject apre un progetto da file
func (d *Dashboard) openProject() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, d.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		// Leggi i dati
		data := make([]byte, 0)
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				data = append(data, buf[:n]...)
			}
			if err != nil {
				break
			}
		}

		// Deserializza il progetto
		project, err := lib.DeserializeProject(data)
		if err != nil {
			dialog.ShowError(fmt.Errorf("errore durante il caricamento: %w", err), d.window)
			return
		}

		// Estrai i fornitori dal progetto serializzato (se presenti)
		// Per ora, i fornitori vengono gestiti separatamente
		// In una versione completa, potremmo salvare anche i fornitori nel file

		d.appState.SetProject(project)
		d.refreshAllTabs()
	}, d.window)
}

// saveProject salva il progetto corrente
func (d *Dashboard) saveProject() {
	project := d.appState.GetProject()
	if project == nil || project.Root == nil {
		dialog.ShowInformation("Errore", "Nessun progetto da salvare", d.window)
		return
	}

	// Apri dialog per salvare file
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, d.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		// Serializza il progetto
		data, err := lib.SerializeProject(project)
		if err != nil {
			dialog.ShowError(fmt.Errorf("errore durante la serializzazione: %w", err), d.window)
			return
		}

		// Scrivi nel file
		if _, err := writer.Write(data); err != nil {
			dialog.ShowError(fmt.Errorf("errore durante il salvataggio: %w", err), d.window)
			return
		}

		dialog.ShowInformation("Successo", "Progetto salvato con successo", d.window)
	}, d.window)
}
