package ui

import (
	"explosio/lib"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateTreeNodeTemplate crea il template di riga per l'albero (icona + label in HBox con padding).
func CreateTreeNodeTemplate() fyne.CanvasObject {
	icon := widget.NewIcon(theme.DocumentIcon())
	label := widget.NewLabel("")
	return container.NewPadded(container.NewHBox(icon, label))
}

// UpdateTreeNodeRow aggiorna il template di riga con icona e testo in base a id e activityMap.
func UpdateTreeNodeRow(id widget.TreeNodeID, branch bool, obj fyne.CanvasObject, activityMap map[string]*lib.Activity) {
	outer := obj.(*fyne.Container)
	box := outer.Objects[0].(*fyne.Container)
	icon := box.Objects[0].(*widget.Icon)
	l := box.Objects[1].(*widget.Label)
	l.Wrapping = fyne.TextWrapOff

	if isResourceNodeID(id) {
		actID, kind, idx, ok := parseResourceNodeID(id)
		if !ok {
			l.SetText("?")
			return
		}
		a := activityMap[actID]
		if a == nil {
			l.SetText("?")
			return
		}
		switch kind {
		case "human":
			icon.SetResource(theme.AccountIcon())
			if idx >= 0 && idx < len(a.Humans) {
				h := a.Humans[idx]
				l.SetText(fmt.Sprintf("Umano: %s (%.1f h, €%.2f/h)", h.Role, h.Quantity, h.CostPerH))
			} else {
				l.SetText("Umano ?")
			}
		case "material":
			icon.SetResource(theme.StorageIcon())
			if idx >= 0 && idx < len(a.Materials) {
				m := a.Materials[idx]
				l.SetText(fmt.Sprintf("Materiale: %s (qty %.1f, €%.2f/unit)", m.Name, m.Quantity, m.UnitCost))
			} else {
				l.SetText("Materiale ?")
			}
		case "asset":
			icon.SetResource(theme.ComputerIcon())
			if idx >= 0 && idx < len(a.Assets) {
				as := a.Assets[idx]
				l.SetText(fmt.Sprintf("Asset: %s (qty %.1f, €%.2f/use)", as.Name, as.Quantity, as.CostPerUse))
			} else {
				l.SetText("Asset ?")
			}
		default:
			l.SetText("?")
		}
		return
	}

	if branch {
		icon.SetResource(theme.FolderOpenIcon())
	} else {
		icon.SetResource(theme.DocumentIcon())
	}
	a, ok := activityMap[string(id)]
	if !ok || a == nil {
		l.SetText("?")
		return
	}
	n := len(a.Humans) + len(a.Materials) + len(a.Assets)
	if n > 0 {
		l.SetText(fmt.Sprintf("%s: %s (%d min) [%d risorse]", a.ID, a.Name, a.Duration, n))
	} else {
		l.SetText(fmt.Sprintf("%s: %s (%d min)", a.ID, a.Name, a.Duration))
	}
}
