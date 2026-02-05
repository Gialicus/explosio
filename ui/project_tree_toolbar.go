package ui

import (
	"explosio/lib"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewProjectTreeToolbar crea la toolbar con i pulsanti per aggiungere sotto-attività e risorse.
// getActivityMap restituisce la mappa attività corrente (così dopo refresh() la toolbar usa la mappa aggiornata).
// selectedID è un puntatore così le azioni leggono sempre il valore aggiornato dalla selezione nel tree.
func NewProjectTreeToolbar(project *lib.Project, getActivityMap func() map[string]*lib.Activity, selectedID *widget.TreeNodeID, refresh func(), getWindow func() fyne.Window) fyne.CanvasObject {
	addSubActivity := func() {
		if selectedID == nil || isResourceNodeID(*selectedID) {
			return
		}
		activityMap := getActivityMap()
		a := activityMap[string(*selectedID)]
		if a == nil {
			return
		}
		win := getWindow()
		if win == nil {
			return
		}
		nameEntry := widget.NewEntry()
		nameEntry.SetText("Nuova attività")
		durEntry := widget.NewEntry()
		durEntry.SetText("1")
		form := widget.NewForm(
			widget.NewFormItem("Nome", nameEntry),
			widget.NewFormItem("Durata (min)", durEntry),
		)
		var pop *widget.PopUp
		form.OnSubmit = func() {
			dur := 1
			fmt.Sscanf(durEntry.Text, "%d", &dur)
			if dur < 1 {
				dur = 1
			}
			child := project.Node(nameEntry.Text, "", dur)
			a.SubActivities = append(a.SubActivities, child)
			child.Next = append(child.Next, a.ID)
			refresh()
			if pop != nil {
				pop.Hide()
			}
		}
		form.OnCancel = func() {
			if pop != nil {
				pop.Hide()
			}
		}
		pop = widget.NewModalPopUp(container.NewVBox(widget.NewLabel("Aggiungi sotto-attività"), form), win.Canvas())
		pop.Resize(fyne.NewSize(400, 200))
		pop.Show()
	}

	addHuman := func() {
		if selectedID == nil || isResourceNodeID(*selectedID) {
			return
		}
		activityMap := getActivityMap()
		a := activityMap[string(*selectedID)]
		if a == nil {
			return
		}
		a.WithHuman("Nuovo ruolo", "", 15, 1)
		refresh()
	}

	addMaterial := func() {
		if selectedID == nil || isResourceNodeID(*selectedID) {
			return
		}
		activityMap := getActivityMap()
		a := activityMap[string(*selectedID)]
		if a == nil {
			return
		}
		a.WithMaterial("Nuovo materiale", "", 1, 1)
		refresh()
	}

	addAsset := func() {
		if selectedID == nil || isResourceNodeID(*selectedID) {
			return
		}
		activityMap := getActivityMap()
		a := activityMap[string(*selectedID)]
		if a == nil {
			return
		}
		a.WithAsset("Nuovo asset", "", 1, 1)
		refresh()
	}

	return container.NewHBox(
		widget.NewButton("Aggiungi sotto-attività", addSubActivity),
		widget.NewButton("Aggiungi risorsa umana", addHuman),
		widget.NewButton("Aggiungi materiale", addMaterial),
		widget.NewButton("Aggiungi asset", addAsset),
	)
}
