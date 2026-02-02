package ui

import (
	"explosio/lib"
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// MermaidView mostra il diagramma Mermaid
type MermaidView struct {
	appState  *AppState
	window    fyne.Window
	textArea  *widget.Entry
}

// NewMermaidView crea una nuova vista Mermaid
func NewMermaidView(appState *AppState, window fyne.Window) fyne.CanvasObject {
	view := &MermaidView{
		appState: appState,
		window:   window,
	}

	// Area di testo per il diagramma
	view.textArea = widget.NewMultiLineEntry()
	view.textArea.SetText("Nessun progetto caricato. Crea o carica un progetto per vedere il diagramma.")

	// Pulsante per aggiornare
	refreshBtn := widget.NewButton("Aggiorna Diagramma", func() {
		view.updateDiagram()
	})

	// Pulsante per esportare
	exportBtn := widget.NewButton("Esporta in File", func() {
		view.exportDiagram()
	})

	// Layout
	content := container.NewBorder(
		container.NewHBox(refreshBtn, exportBtn),
		nil,
		nil,
		nil,
		container.NewScroll(view.textArea),
	)

	// Aggiorna quando cambia il progetto
	appState.OnProjectChanged(func(*lib.Project) {
		view.updateDiagram()
	})

	// Aggiorna inizialmente
	view.updateDiagram()

	return content
}

// updateDiagram aggiorna il diagramma Mermaid
func (v *MermaidView) updateDiagram() {
	project := v.appState.GetProject()
	if project == nil || project.Root == nil {
		v.textArea.SetText("Nessun progetto caricato.")
		return
	}

	// Calcola CPM se necessario
	engine := v.appState.GetEngine()
	engine.ComputeCPM(project.Root)

	// Genera il diagramma Mermaid
	diagram := lib.GenerateMermaid(project.Root)
	v.textArea.SetText(diagram)
}

// exportDiagram esporta il diagramma in un file
func (v *MermaidView) exportDiagram() {
	project := v.appState.GetProject()
	if project == nil || project.Root == nil {
		dialog.ShowInformation("Errore", "Nessun progetto caricato", v.window)
		return
	}

	// Apri dialog per salvare file
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, v.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		// Genera il diagramma
		engine := v.appState.GetEngine()
		engine.ComputeCPM(project.Root)
		diagram := lib.GenerateMermaid(project.Root)

		// Scrivi nel file
		if _, err := writer.Write([]byte(diagram)); err != nil {
			dialog.ShowError(fmt.Errorf("errore durante il salvataggio: %w", err), v.window)
			return
		}

		dialog.ShowInformation("Successo", "Diagramma esportato con successo", v.window)
	}, v.window)
}
