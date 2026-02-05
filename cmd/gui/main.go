package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.NewWithID("com.explosio.app")
	window := myApp.NewWindow("Explosio")
	window.Resize(fyne.NewSize(800, 600))
	window.CenterOnScreen()
	window.SetContent(container.NewCenter(widget.NewLabel("Explosio")))
	window.ShowAndRun()
}
