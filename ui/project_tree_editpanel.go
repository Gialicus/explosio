package ui

import (
	"explosio/lib"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// BuildEditPanelContent restituisce i widget da mostrare nel pannello di modifica per il nodo selezionato.
// Se selectedID è vuoto restituisce nil (nessun contenuto). Altrimenti form + pulsante Applica.
func BuildEditPanelContent(selectedID widget.TreeNodeID, activityMap map[string]*lib.Activity, refresh func()) []fyne.CanvasObject {
	if selectedID == "" {
		return nil
	}
	if isResourceNodeID(selectedID) {
		actID, kind, idx, ok := parseResourceNodeID(selectedID)
		if !ok {
			return nil
		}
		a := activityMap[actID]
		if a == nil || idx < 0 {
			return nil
		}
		switch kind {
		case "human":
			if idx >= len(a.Humans) {
				return nil
			}
			h := &a.Humans[idx]
			nameEntry := widget.NewEntry()
			nameEntry.SetText(h.Role)
			nameEntry.PlaceHolder = "Ruolo"
			numEntry := widget.NewEntry()
			numEntry.SetText(fmt.Sprintf("%.2f", h.Quantity))
			numEntry.PlaceHolder = "Ore"
			costEntry := widget.NewEntry()
			costEntry.SetText(fmt.Sprintf("%.2f", h.CostPerH))
			costEntry.PlaceHolder = "€/h"
			apply := func() {
				fmt.Sscanf(numEntry.Text, "%f", &h.Quantity)
				fmt.Sscanf(costEntry.Text, "%f", &h.CostPerH)
				h.Role = nameEntry.Text
				refresh()
			}
			numEntry.OnSubmitted = func(string) { apply() }
			costEntry.OnSubmitted = func(string) { apply() }
			return []fyne.CanvasObject{
				widget.NewLabel("Risorsa umana"),
				container.NewGridWithColumns(2, widget.NewLabel("Ruolo"), nameEntry, widget.NewLabel("Ore"), numEntry, widget.NewLabel("€/h"), costEntry),
				widget.NewButton("Applica", apply),
			}
		case "material":
			if idx >= len(a.Materials) {
				return nil
			}
			m := &a.Materials[idx]
			nameEntry := widget.NewEntry()
			nameEntry.SetText(m.Name)
			numEntry := widget.NewEntry()
			numEntry.SetText(fmt.Sprintf("%.2f", m.Quantity))
			costEntry := widget.NewEntry()
			costEntry.SetText(fmt.Sprintf("%.2f", m.UnitCost))
			apply := func() {
				fmt.Sscanf(numEntry.Text, "%f", &m.Quantity)
				fmt.Sscanf(costEntry.Text, "%f", &m.UnitCost)
				m.Name = nameEntry.Text
				refresh()
			}
			return []fyne.CanvasObject{
				widget.NewLabel("Materiale"),
				container.NewGridWithColumns(2, widget.NewLabel("Nome"), nameEntry, widget.NewLabel("Q.tà"), numEntry, widget.NewLabel("€/un"), costEntry),
				widget.NewButton("Applica", apply),
			}
		case "asset":
			if idx >= len(a.Assets) {
				return nil
			}
			as := &a.Assets[idx]
			nameEntry := widget.NewEntry()
			nameEntry.SetText(as.Name)
			numEntry := widget.NewEntry()
			numEntry.SetText(fmt.Sprintf("%.2f", as.Quantity))
			costEntry := widget.NewEntry()
			costEntry.SetText(fmt.Sprintf("%.2f", as.CostPerUse))
			apply := func() {
				fmt.Sscanf(numEntry.Text, "%f", &as.Quantity)
				fmt.Sscanf(costEntry.Text, "%f", &as.CostPerUse)
				as.Name = nameEntry.Text
				refresh()
			}
			return []fyne.CanvasObject{
				widget.NewLabel("Asset"),
				container.NewGridWithColumns(2, widget.NewLabel("Nome"), nameEntry, widget.NewLabel("Q.tà"), numEntry, widget.NewLabel("€/uso"), costEntry),
				widget.NewButton("Applica", apply),
			}
		default:
			return nil
		}
	}
	a := activityMap[string(selectedID)]
	if a == nil {
		return nil
	}
	nameEntry := widget.NewEntry()
	nameEntry.SetText(a.Name)
	nameEntry.PlaceHolder = "Nome"
	durEntry := widget.NewEntry()
	durEntry.SetText(fmt.Sprintf("%d", a.Duration))
	durEntry.PlaceHolder = "Durata (min)"
	apply := func() {
		fmt.Sscanf(durEntry.Text, "%d", &a.Duration)
		if a.Duration < 1 {
			a.Duration = 1
		}
		a.Name = nameEntry.Text
		refresh()
	}
	return []fyne.CanvasObject{
		widget.NewLabel("Attività"),
		container.NewGridWithColumns(2, widget.NewLabel("Nome"), nameEntry, widget.NewLabel("Durata (min)"), durEntry),
		widget.NewButton("Applica", apply),
	}
}
