package ui

import (
	"explosio/lib"
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ProjectEditor gestisce l'editing del progetto
type ProjectEditor struct {
	appState     *AppState
	onChanged    func()
	rootForm     *widget.Form
	activityTree *widget.Tree
	activityData map[string]*ActivityNode
}

// ActivityNode rappresenta un nodo nell'albero delle attività
type ActivityNode struct {
	Activity *lib.Activity
	Parent   *ActivityNode
	Children []*ActivityNode
}

// NewProjectEditor crea un nuovo editor di progetto
func NewProjectEditor(appState *AppState, onChanged func()) fyne.CanvasObject {
	editor := &ProjectEditor{
		appState:     appState,
		onChanged:    onChanged,
		activityData: make(map[string]*ActivityNode),
	}

	// Form per il progetto root
	editor.rootForm = editor.createRootForm()

	// Tree view per le attività
	editor.activityTree = editor.createActivityTree()

	// Layout: form a sinistra, tree a destra
	split := container.NewHSplit(
		container.NewScroll(editor.rootForm),
		container.NewScroll(editor.activityTree),
	)
	split.SetOffset(0.3)

	return split
}

// createRootForm crea il form per modificare il progetto root
func (e *ProjectEditor) createRootForm() *widget.Form {
	project := e.appState.GetProject()
	
	nameEntry := widget.NewEntry()
	descEntry := widget.NewMultiLineEntry()
	durationEntry := widget.NewEntry()

	if project != nil && project.Root != nil {
		nameEntry.SetText(project.Root.Name)
		descEntry.SetText(project.Root.Description)
		durationEntry.SetText(fmt.Sprintf("%d", project.Root.Duration))
	}

	// Callback per salvare le modifiche
	onSubmit := func() {
		project := e.appState.GetProject()
		if project == nil {
			project = lib.NewProject()
		}

		duration := 0
		fmt.Sscanf(durationEntry.Text, "%d", &duration)
		if duration <= 0 {
			duration = 1
		}

		if project.Root == nil {
			// Usa Start per creare la root con ID corretto
			project.Start(nameEntry.Text, descEntry.Text, duration)
		} else {
			project.Root.Name = nameEntry.Text
			project.Root.Description = descEntry.Text
			project.Root.Duration = duration
			if project.Root.MinDuration > duration {
				project.Root.MinDuration = duration
			}
		}

		e.appState.SetProject(project)
		if e.onChanged != nil {
			e.onChanged()
		}
	}

	form := widget.NewForm(
		widget.NewFormItem("Nome", nameEntry),
		widget.NewFormItem("Descrizione", descEntry),
		widget.NewFormItem("Durata (min)", durationEntry),
	)
	form.OnSubmit = onSubmit

	return form
}

// createActivityTree crea l'albero delle attività
func (e *ProjectEditor) createActivityTree() *widget.Tree {
	tree := widget.NewTree(
		e.treeChildUIDs,
		e.treeIsBranch,
		e.treeCreate,
		e.treeUpdate,
	)

	tree.OnSelected = func(id widget.TreeNodeID) {
		// Quando si seleziona un nodo, apri l'editor dell'attività
		if node, ok := e.activityData[id]; ok {
			e.openActivityEditor(node.Activity)
		}
	}

	return tree
}

// treeChildUIDs restituisce gli ID dei figli di un nodo
func (e *ProjectEditor) treeChildUIDs(id widget.TreeNodeID) []widget.TreeNodeID {
	project := e.appState.GetProject()
	if project == nil || project.Root == nil {
		return []widget.TreeNodeID{}
	}

	// Costruisci la struttura ad albero se necessario
	if len(e.activityData) == 0 {
		e.buildActivityTree(project.Root, nil)
	}

	// Se è la root, restituisci il suo ID
	if id == "" {
		return []widget.TreeNodeID{project.Root.ID}
	}

	// Altrimenti, restituisci i figli del nodo
	if node, ok := e.activityData[id]; ok {
		children := make([]widget.TreeNodeID, len(node.Children))
		for i, child := range node.Children {
			children[i] = child.Activity.ID
		}
		return children
	}

	return []widget.TreeNodeID{}
}

// buildActivityTree costruisce ricorsivamente l'albero delle attività
func (e *ProjectEditor) buildActivityTree(activity *lib.Activity, parent *ActivityNode) *ActivityNode {
	if activity == nil {
		return nil
	}

	node := &ActivityNode{
		Activity: activity,
		Parent:   parent,
		Children: make([]*ActivityNode, 0),
	}
	e.activityData[activity.ID] = node

	// Aggiungi i figli
	for _, sub := range activity.SubActivities {
		childNode := e.buildActivityTree(sub, node)
		if childNode != nil {
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}

// treeIsBranch verifica se un nodo è un branch (ha figli)
func (e *ProjectEditor) treeIsBranch(id widget.TreeNodeID) bool {
	if node, ok := e.activityData[id]; ok {
		return len(node.Children) > 0
	}
	return false
}

// treeCreate crea il widget per un nodo
func (e *ProjectEditor) treeCreate(isBranch bool) fyne.CanvasObject {
	return widget.NewLabel("")
}

// treeUpdate aggiorna il widget di un nodo
func (e *ProjectEditor) treeUpdate(id widget.TreeNodeID, isBranch bool, obj fyne.CanvasObject) {
	if node, ok := e.activityData[id]; ok && node.Activity != nil {
		label := obj.(*widget.Label)
		label.SetText(fmt.Sprintf("%s: %s", node.Activity.ID, node.Activity.Name))
	}
}

// openActivityEditor apre l'editor per un'attività
func (e *ProjectEditor) openActivityEditor(activity *lib.Activity) {
	if activity == nil {
		return
	}

	editor := NewActivityEditor(e.appState, activity, func() {
		e.onChanged()
		// Ricostruisci l'albero
		e.activityData = make(map[string]*ActivityNode)
		e.activityTree.Refresh()
	})

	// Ottieni la window corrente per il dialog
	var canvas fyne.Canvas
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}

	// Crea una finestra dialog per l'editor
	var popup *widget.PopUp
	popup = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Modifica Attività: %s", activity.Name)),
			editor,
			widget.NewButton("Chiudi", func() {
				if popup != nil {
					popup.Hide()
				}
			}),
		),
		canvas,
	)
	popup.Resize(fyne.NewSize(600, 500))
	popup.Show()
}

