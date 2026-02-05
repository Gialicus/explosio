package main

import (
	"explosio/lib"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// newProjectTab restituisce il contenuto del tab Progetto: albero attività/risorse o stato vuoto.
func newProjectTab(project *lib.Project) fyne.CanvasObject {
	mainContent := container.NewStack()

	emptyState := container.NewCenter(container.NewVBox(
		widget.NewLabel("Nessun progetto. Aggiungi un'attività radice per iniziare."),
		widget.NewButton("Aggiungi attività radice", func() {
			project.Start("Nuova attività", "", 1)
			treeContent := buildTreeContent(project, mainContent)
			mainContent.RemoveAll()
			mainContent.Add(treeContent)
			mainContent.Refresh()
		}),
	))

	if project.Root == nil {
		mainContent.Add(emptyState)
		return mainContent
	}

	treeContent := buildTreeContent(project, mainContent)
	mainContent.Add(treeContent)
	return mainContent
}

func buildTreeContent(project *lib.Project, mainContent *fyne.Container) fyne.CanvasObject {
	var activityMap map[string]*lib.Activity
	var tree *widget.Tree
	var selectedID widget.TreeNodeID

	buildMap := func() {
		activityMap = make(map[string]*lib.Activity)
		if project.Root != nil {
			buildMapRec(project.Root, activityMap)
		}
	}
	buildMap()

	refresh := func() {
		buildMap()
		if tree != nil && project.Root != nil {
			tree.Root = widget.TreeNodeID(project.Root.ID)
			tree.Refresh()
		}
	}

	childUIDs := func(id widget.TreeNodeID) []widget.TreeNodeID {
		if isResourceNodeID(id) {
			return nil
		}
		a, ok := activityMap[string(id)]
		if !ok || a == nil {
			return nil
		}
		var out []widget.TreeNodeID
		for _, sub := range a.SubActivities {
			out = append(out, widget.TreeNodeID(sub.ID))
		}
		for i := range a.Humans {
			out = append(out, widget.TreeNodeID(fmt.Sprintf("%s|human|%d", a.ID, i)))
		}
		for i := range a.Materials {
			out = append(out, widget.TreeNodeID(fmt.Sprintf("%s|material|%d", a.ID, i)))
		}
		for i := range a.Assets {
			out = append(out, widget.TreeNodeID(fmt.Sprintf("%s|asset|%d", a.ID, i)))
		}
		return out
	}

	isBranch := func(id widget.TreeNodeID) bool {
		if isResourceNodeID(id) {
			return false
		}
		a, ok := activityMap[string(id)]
		if !ok || a == nil {
			return false
		}
		return len(a.SubActivities) > 0 || len(a.Humans) > 0 || len(a.Materials) > 0 || len(a.Assets) > 0
	}

	createNode := func(bool) fyne.CanvasObject {
		return widget.NewLabel("")
	}

	updateNode := func(id widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
		l := obj.(*widget.Label)
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
				if idx >= 0 && idx < len(a.Humans) {
					h := a.Humans[idx]
					l.SetText(fmt.Sprintf("Umano: %s (%.1f h, €%.2f/h)", h.Role, h.Quantity, h.CostPerH))
				} else {
					l.SetText("Umano ?")
				}
			case "material":
				if idx >= 0 && idx < len(a.Materials) {
					m := a.Materials[idx]
					l.SetText(fmt.Sprintf("Materiale: %s (qty %.1f, €%.2f/unit)", m.Name, m.Quantity, m.UnitCost))
				} else {
					l.SetText("Materiale ?")
				}
			case "asset":
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

	tree = widget.NewTree(childUIDs, isBranch, createNode, updateNode)

	tree.Root = widget.TreeNodeID(project.Root.ID)

	// Pannello sotto l'albero: mostra campi editabili per il nodo selezionato
	editPanel := container.NewVBox()
	updateEditPanel := func() {
		editPanel.RemoveAll()
		if selectedID == "" {
			editPanel.Refresh()
			return
		}
		if isResourceNodeID(selectedID) {
			actID, kind, idx, ok := parseResourceNodeID(selectedID)
			if !ok {
				editPanel.Refresh()
				return
			}
			a := activityMap[actID]
			if a == nil || idx < 0 {
				editPanel.Refresh()
				return
			}
			var nameEntry, numEntry, costEntry *widget.Entry
			switch kind {
			case "human":
				if idx >= len(a.Humans) {
					editPanel.Refresh()
					return
				}
				h := &a.Humans[idx]
				nameEntry = widget.NewEntry()
				nameEntry.SetText(h.Role)
				nameEntry.PlaceHolder = "Ruolo"
				numEntry = widget.NewEntry()
				numEntry.SetText(fmt.Sprintf("%.2f", h.Quantity))
				numEntry.PlaceHolder = "Ore"
				costEntry = widget.NewEntry()
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
				editPanel.Add(widget.NewLabel("Risorsa umana"))
				editPanel.Add(container.NewGridWithColumns(2, widget.NewLabel("Ruolo"), nameEntry, widget.NewLabel("Ore"), numEntry, widget.NewLabel("€/h"), costEntry))
				editPanel.Add(widget.NewButton("Applica", apply))
			case "material":
				if idx >= len(a.Materials) {
					editPanel.Refresh()
					return
				}
				m := &a.Materials[idx]
				nameEntry = widget.NewEntry()
				nameEntry.SetText(m.Name)
				numEntry = widget.NewEntry()
				numEntry.SetText(fmt.Sprintf("%.2f", m.Quantity))
				costEntry = widget.NewEntry()
				costEntry.SetText(fmt.Sprintf("%.2f", m.UnitCost))
				apply := func() {
					fmt.Sscanf(numEntry.Text, "%f", &m.Quantity)
					fmt.Sscanf(costEntry.Text, "%f", &m.UnitCost)
					m.Name = nameEntry.Text
					refresh()
				}
				editPanel.Add(widget.NewLabel("Materiale"))
				editPanel.Add(container.NewGridWithColumns(2, widget.NewLabel("Nome"), nameEntry, widget.NewLabel("Q.tà"), numEntry, widget.NewLabel("€/un"), costEntry))
				editPanel.Add(widget.NewButton("Applica", apply))
			case "asset":
				if idx >= len(a.Assets) {
					editPanel.Refresh()
					return
				}
				as := &a.Assets[idx]
				nameEntry = widget.NewEntry()
				nameEntry.SetText(as.Name)
				numEntry = widget.NewEntry()
				numEntry.SetText(fmt.Sprintf("%.2f", as.Quantity))
				costEntry = widget.NewEntry()
				costEntry.SetText(fmt.Sprintf("%.2f", as.CostPerUse))
				apply := func() {
					fmt.Sscanf(numEntry.Text, "%f", &as.Quantity)
					fmt.Sscanf(costEntry.Text, "%f", &as.CostPerUse)
					as.Name = nameEntry.Text
					refresh()
				}
				editPanel.Add(widget.NewLabel("Asset"))
				editPanel.Add(container.NewGridWithColumns(2, widget.NewLabel("Nome"), nameEntry, widget.NewLabel("Q.tà"), numEntry, widget.NewLabel("€/uso"), costEntry))
				editPanel.Add(widget.NewButton("Applica", apply))
			default:
				editPanel.Refresh()
				return
			}
		} else {
			a := activityMap[string(selectedID)]
			if a == nil {
				editPanel.Refresh()
				return
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
			editPanel.Add(widget.NewLabel("Attività"))
			editPanel.Add(container.NewGridWithColumns(2, widget.NewLabel("Nome"), nameEntry, widget.NewLabel("Durata (min)"), durEntry))
			editPanel.Add(widget.NewButton("Applica", apply))
		}
		editPanel.Refresh()
	}

	tree.OnSelected = func(id widget.TreeNodeID) {
		selectedID = id
		updateEditPanel()
	}

	getWindow := func() fyne.Window {
		w := fyne.CurrentApp().Driver().AllWindows()
		if len(w) > 0 {
			return w[0]
		}
		return nil
	}

	addSubActivity := func() {
		if !isResourceNodeID(selectedID) {
			a := activityMap[string(selectedID)]
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
	}

	addHuman := func() {
		if isResourceNodeID(selectedID) {
			return
		}
		a := activityMap[string(selectedID)]
		if a == nil {
			return
		}
		a.WithHuman("Nuovo ruolo", "", 15, 1)
		refresh()
	}

	addMaterial := func() {
		if isResourceNodeID(selectedID) {
			return
		}
		a := activityMap[string(selectedID)]
		if a == nil {
			return
		}
		a.WithMaterial("Nuovo materiale", "", 1, 1)
		refresh()
	}

	addAsset := func() {
		if isResourceNodeID(selectedID) {
			return
		}
		a := activityMap[string(selectedID)]
		if a == nil {
			return
		}
		a.WithAsset("Nuovo asset", "", 1, 1)
		refresh()
	}

	toolbar := container.NewHBox(
		widget.NewButton("Aggiungi sotto-attività", addSubActivity),
		widget.NewButton("Aggiungi risorsa umana", addHuman),
		widget.NewButton("Aggiungi materiale", addMaterial),
		widget.NewButton("Aggiungi asset", addAsset),
	)

	treeScroll := container.NewScroll(tree)
	// Pannello edit sotto l'albero (vuoto finché non si seleziona un nodo)
	updateEditPanel()
	return container.NewBorder(toolbar, editPanel, nil, nil, treeScroll)
}

func buildMapRec(a *lib.Activity, m map[string]*lib.Activity) {
	if a == nil {
		return
	}
	m[a.ID] = a
	for _, sub := range a.SubActivities {
		buildMapRec(sub, m)
	}
}

func isResourceNodeID(id widget.TreeNodeID) bool {
	return strings.Contains(string(id), "|")
}

func parseResourceNodeID(id widget.TreeNodeID) (activityID, kind string, index int, ok bool) {
	parts := strings.Split(string(id), "|")
	if len(parts) != 3 {
		return "", "", 0, false
	}
	idx, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", 0, false
	}
	return parts[0], parts[1], idx, true
}
