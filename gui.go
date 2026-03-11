package main

import (
	"explosio/gui"

	"fyne.io/fyne/v2/app"
)

func runGUI() {
	app.New() // Inizializza Fyne prima di gui.Run
	gui.Run(BuildDemoTree())
}
