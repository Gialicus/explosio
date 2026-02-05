package main

import (
	"explosio/lib"
	"explosio/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	myApp := app.NewWithID("com.explosio.app")
	window := myApp.NewWindow("Explosio")
	window.Resize(fyne.NewSize(800, 600))
	window.CenterOnScreen()

	project := lib.NewProject()
	progettoContent := ui.NewProjectTab(project)
	simulazioneContent := container.NewStack()
	tabs := container.NewAppTabs(
		container.NewTabItem("Progetto", progettoContent),
		container.NewTabItem("Simulazione", simulazioneContent),
	)
	window.SetContent(tabs)

	window.ShowAndRun()
}
