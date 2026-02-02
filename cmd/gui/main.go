package main

import (
	"explosio/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	// Crea l'applicazione Fyne
	myApp := app.NewWithID("com.explosio.app")
	// Usa il tema di default di Fyne (rimosso ExplosioTheme per evitare ricorsione)

	// Crea la finestra principale
	window := myApp.NewWindow("Explosio - Analisi Progetti")
	window.Resize(ui.WindowSize)
	window.CenterOnScreen()

	// Crea lo stato dell'applicazione
	appState := ui.NewAppState()

	// Crea la dashboard principale
	content := ui.NewDashboard(window, appState)
	window.SetContent(content)

	// Mostra la finestra
	window.ShowAndRun()
}
