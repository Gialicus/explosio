package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func runGUI() {
	a := app.New()
	w := a.NewWindow("Explosio")
	w.SetContent(widget.NewLabel("Explosio - Activity tree modelling"))
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}
