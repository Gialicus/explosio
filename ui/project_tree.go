package ui

import (
	"explosio/lib"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewProjectTab restituisce il contenuto del tab Progetto: albero attività/risorse o stato vuoto.
// Se onProjectChange non è nil, viene chiamato quando l'utente aggiunge un'attività radice (per aggiornare altri tab, es. Simulazione).
func NewProjectTab(project *lib.Project, onProjectChange func()) fyne.CanvasObject {
	mainContent := container.NewStack()

	emptyState := container.NewCenter(container.NewVBox(
		widget.NewLabel("Nessun progetto. Aggiungi un'attività radice per iniziare."),
		widget.NewButton("Aggiungi attività radice", func() {
			project.Start("Nuova attività", "", 1)
			treeContent := buildTreeContent(project, mainContent)
			mainContent.RemoveAll()
			mainContent.Add(treeContent)
			mainContent.Refresh()
			if onProjectChange != nil {
				onProjectChange()
			}
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

	activityMap = BuildActivityMap(project)

	refresh := func() {
		activityMap = BuildActivityMap(project)
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
		return CreateTreeNodeTemplate()
	}
	updateNode := func(id widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
		UpdateTreeNodeRow(id, branch, obj, activityMap)
	}

	tree = widget.NewTree(childUIDs, isBranch, createNode, updateNode)

	tree.Root = widget.TreeNodeID(project.Root.ID)

	// Pannello sotto l'albero: Card "Modifica nodo" con form o messaggio se nessuna selezione
	editCard := widget.NewCard("Modifica nodo", "", container.NewCenter(
		widget.NewLabel("Seleziona un'attività o una risorsa nell'albero per modificarla."),
	))
	updateEditPanel := func() {
		objs := BuildEditPanelContent(selectedID, activityMap, refresh)
		if len(objs) == 0 {
			editCard.SetContent(container.NewCenter(
				widget.NewLabel("Seleziona un'attività o una risorsa nell'albero per modificarla."),
			))
		} else {
			editCard.SetContent(container.NewVBox(objs...))
		}
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

	getActivityMap := func() map[string]*lib.Activity { return activityMap }
	toolbar := NewProjectTreeToolbar(project, getActivityMap, &selectedID, refresh, getWindow)

	treeScroll := container.NewScroll(tree)
	updateEditPanel()
	editPanelWithPadding := container.NewPadded(editCard)
	return container.NewBorder(toolbar, editPanelWithPadding, nil, nil, treeScroll)
}
