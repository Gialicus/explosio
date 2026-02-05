package ui

import (
	"explosio/lib"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ProjectEditor gestisce l'editing del progetto
type ProjectEditor struct {
	appState     *AppState
	onChanged    func()
	rootForm     fyne.CanvasObject
	activityTree *widget.Tree
	activityMap  map[string]*lib.Activity // Mappa semplificata per lookup veloce
	selectedID   string                   // ID dell'attività selezionata per operazioni contestuali
	contextPopup *widget.PopUp            // Popup per menu contestuale
}

// NewProjectEditor crea un nuovo editor di progetto
func NewProjectEditor(appState *AppState, onChanged func()) fyne.CanvasObject {
	editor := &ProjectEditor{
		appState:    appState,
		onChanged:   onChanged,
		activityMap: make(map[string]*lib.Activity),
	}

	// Form per il progetto root
	editor.rootForm = editor.createRootForm()

	// Tree view per le attività
	editor.activityTree = editor.createActivityTree()

	// Aggiorna l'albero quando cambia il progetto
	appState.OnProjectChanged(func(*lib.Project) {
		editor.refreshTree()
	})

	// Pulsanti per gestire le attività
	addActivityBtn := widget.NewButton("+ Aggiungi Attività", func() {
		editor.addActivityToRoot()
	})

	editSelectedBtn := widget.NewButton("✏️ Modifica Selezionata", func() {
		if editor.selectedID != "" {
			if activity, ok := editor.activityMap[editor.selectedID]; ok {
				editor.openActivityEditor(activity)
			}
		} else {
			dialog.ShowInformation("Attenzione", "Seleziona un'attività dall'albero", editor.getWindow())
		}
	})

	deleteSelectedBtn := widget.NewButton("🗑️ Elimina Selezionata", func() {
		if editor.selectedID != "" {
			if activity, ok := editor.activityMap[editor.selectedID]; ok {
				editor.deleteActivity(activity)
			}
		} else {
			dialog.ShowInformation("Attenzione", "Seleziona un'attività dall'albero", editor.getWindow())
		}
	})

	refreshBtn := widget.NewButton("Aggiorna Albero", func() {
		editor.refreshTree()
	})

	// Wrappare l'albero in uno scroll per gestire lo scroll interno
	treeScroll := container.NewScroll(editor.activityTree)

	// Layout verticale: toolbar in alto, tree in basso
	treeContainer := container.NewBorder(
		container.NewHBox(
			addActivityBtn,
			editSelectedBtn,
			deleteSelectedBtn,
			widget.NewSeparator(),
			refreshBtn,
		),
		nil,
		nil,
		nil,
		treeScroll,
	)

	// Layout principale: form in alto (compatto), albero in basso (espandibile)
	// Il form deve essere limitato in altezza per non occupare tutto lo spazio
	formWithSeparator := container.NewVBox(
		editor.rootForm,
		widget.NewSeparator(),
	)

	// Border layout: top = form, center = tree (non bottom!)
	// Il center si espande per occupare tutto lo spazio rimanente
	content := container.NewBorder(
		formWithSeparator,
		nil,           // bottom = nil
		nil,           // left = nil
		nil,           // right = nil
		treeContainer, // center = tree (si espande)
	)

	return content
}

// refreshTree ricostruisce e aggiorna l'albero
func (e *ProjectEditor) refreshTree() {
	// Pulisci la mappa esistente
	e.activityMap = make(map[string]*lib.Activity)

	// Ricostruisci la mappa
	project := e.appState.GetProject()
	if project != nil && project.Root != nil {
		e.buildActivityMap(project.Root)
		if e.activityTree != nil {
			e.activityTree.Root = project.Root.ID
			e.activityTree.Refresh()
		}
	}
}

// buildActivityMap costruisce ricorsivamente la mappa delle attività
func (e *ProjectEditor) buildActivityMap(activity *lib.Activity) {
	if activity == nil {
		return
	}
	e.activityMap[activity.ID] = activity
	for _, sub := range activity.SubActivities {
		e.buildActivityMap(sub)
	}
}

// createRootForm crea il form per modificare il progetto root
func (e *ProjectEditor) createRootForm() fyne.CanvasObject {
	project := e.appState.GetProject()

	nameEntry := widget.NewEntry()
	descEntry := widget.NewMultiLineEntry()
	descEntry.Wrapping = fyne.TextWrapWord
	// Limita l'altezza del campo descrizione per non far espandere troppo il form
	// Usa un'altezza fissa approssimativa (circa 3 righe)
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
		e.refreshTree()
	}

	// Salva automaticamente quando cambiano i valori
	nameEntry.OnChanged = func(string) {
		onSubmit()
	}
	descEntry.OnChanged = func(string) {
		onSubmit()
	}
	durationEntry.OnChanged = func(string) {
		onSubmit()
	}

	form := widget.NewForm(
		widget.NewFormItem("Nome Progetto", nameEntry),
		widget.NewFormItem("Descrizione", descEntry),
		widget.NewFormItem("Durata (min)", durationEntry),
	)
	form.OnSubmit = onSubmit

	infoLabel := widget.NewRichTextFromMarkdown("**Nota:** Modifica i campi per creare o aggiornare la root activity del progetto.")

	formContent := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		form,
	)

	// Wrappare in un container con altezza limitata per evitare che si espanda troppo
	// e nasconda l'albero. Il Border layout darà al form solo lo spazio necessario.
	return formContent
}

// createActivityTree crea l'albero delle attività
func (e *ProjectEditor) createActivityTree() *widget.Tree {
	project := e.appState.GetProject()

	// Costruisci la mappa
	if project != nil && project.Root != nil {
		e.buildActivityMap(project.Root)
	}

	tree := widget.NewTree(
		e.treeChildUIDs,
		e.treeIsBranch,
		e.treeCreate,
		e.treeUpdate,
	)

	// Imposta la root del tree
	if project != nil && project.Root != nil {
		tree.Root = project.Root.ID
	}

	// Click singolo per aprire editor e selezionare
	tree.OnSelected = func(id widget.TreeNodeID) {
		e.selectedID = id
		if activity, ok := e.activityMap[id]; ok {
			e.openActivityEditor(activity)
		}
	}

	return tree
}

// treeChildUIDs restituisce gli ID dei figli di un nodo (solo SubActivities)
func (e *ProjectEditor) treeChildUIDs(id widget.TreeNodeID) []widget.TreeNodeID {
	project := e.appState.GetProject()
	if project == nil || project.Root == nil {
		return []widget.TreeNodeID{}
	}

	// Costruisci la mappa se necessario
	if len(e.activityMap) == 0 {
		e.buildActivityMap(project.Root)
	}

	// Trova l'attività
	activity, ok := e.activityMap[id]
	if !ok {
		return []widget.TreeNodeID{}
	}

	// Restituisci solo le SubActivities (non le risorse)
	children := make([]widget.TreeNodeID, len(activity.SubActivities))
	for i, sub := range activity.SubActivities {
		children[i] = sub.ID
	}

	return children
}

// treeIsBranch verifica se un nodo è un branch (ha SubActivities o risorse)
func (e *ProjectEditor) treeIsBranch(id widget.TreeNodeID) bool {
	activity, ok := e.activityMap[id]
	if !ok {
		return false
	}

	// È un branch se ha sub-attività o risorse
	hasSubActivities := len(activity.SubActivities) > 0
	hasResources := len(activity.Humans) > 0 || len(activity.Materials) > 0 || len(activity.Assets) > 0
	return hasSubActivities || hasResources
}

// treeCreate crea il widget per un nodo
func (e *ProjectEditor) treeCreate(isBranch bool) fyne.CanvasObject {
	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord
	return label
}

// treeUpdate aggiorna il widget di un nodo
func (e *ProjectEditor) treeUpdate(id widget.TreeNodeID, isBranch bool, obj fyne.CanvasObject) {
	label := obj.(*widget.Label)
	activity, ok := e.activityMap[id]
	if !ok {
		label.SetText("Attività non trovata")
		return
	}

	// Crea il testo principale dell'attività
	icon := "📄 "
	if isBranch {
		icon = "📁 "
	}

	tag := "[ ]"
	if activity.Slack == 0 {
		tag = "[CRITICAL]"
	}

	// Conta le risorse
	resourceCount := len(activity.Humans) + len(activity.Materials) + len(activity.Assets)
	resourceInfo := ""
	if resourceCount > 0 {
		resourceInfo = fmt.Sprintf(" [%d risorse]", resourceCount)
	}

	text := fmt.Sprintf("%s%s %s: %s (%d min)%s", icon, tag, activity.ID, activity.Name, activity.Duration, resourceInfo)
	label.SetText(text)
}

// openActivityEditor apre l'editor per un'attività
func (e *ProjectEditor) openActivityEditor(activity *lib.Activity) {
	if activity == nil {
		return
	}

	var canvas fyne.Canvas
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}

	var popup *widget.PopUp
	var createEditorContent func() fyne.CanvasObject

	createEditorContent = func() fyne.CanvasObject {
		editor := NewActivityEditor(e.appState, activity, func() {
			if popup != nil && createEditorContent != nil {
				popup.Content = createEditorContent()
				popup.Refresh()
			}
			e.onChanged()
			e.refreshTree()
		})

		titleLabel := widget.NewRichTextFromMarkdown(fmt.Sprintf("## Modifica Attività: **%s**", activity.Name))
		header := container.NewBorder(nil, widget.NewSeparator(), titleLabel, nil, nil)

		closeBtn := widget.NewButton("Chiudi", func() {
			if popup != nil {
				popup.Hide()
			}
		})

		footer := container.NewBorder(
			widget.NewSeparator(),
			nil,
			nil,
			closeBtn,
			nil,
		)

		return container.NewBorder(
			header,
			footer,
			nil,
			nil,
			editor,
		)
	}

	editorContainer := createEditorContent()

	popup = widget.NewModalPopUp(
		editorContainer,
		canvas,
	)
	popup.Resize(fyne.NewSize(900, 750))
	popup.Show()
}

// addActivityToRoot aggiunge una nuova attività come figlia della root
func (e *ProjectEditor) addActivityToRoot() {
	project := e.appState.GetProject()
	if project == nil {
		project = lib.NewProject()
		e.appState.SetProject(project)
	}

	if project.Root == nil {
		dialog.ShowInformation("Attenzione", "Crea prima la root activity usando il form sopra", e.getWindow())
		return
	}

	e.showAddActivityDialog(project.Root)
}

// showAddActivityDialog mostra il dialog per aggiungere una nuova attività
func (e *ProjectEditor) showAddActivityDialog(parent *lib.Activity) {
	if parent == nil {
		return
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetText("Nuova Attività")
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText("Descrizione attività")
	durationEntry := widget.NewEntry()
	durationEntry.SetText("1")

	form := widget.NewForm(
		widget.NewFormItem("Nome", nameEntry),
		widget.NewFormItem("Descrizione", descEntry),
		widget.NewFormItem("Durata (min)", durationEntry),
	)

	var popup *widget.PopUp
	form.OnSubmit = func() {
		duration := 1
		fmt.Sscanf(durationEntry.Text, "%d", &duration)
		if duration <= 0 {
			duration = 1
		}

		project := e.appState.GetProject()
		if project == nil {
			return
		}

		// Crea la nuova attività
		newActivity := project.Node(nameEntry.Text, descEntry.Text, duration)

		// Aggiungi come sub-activity
		parent.SubActivities = append(parent.SubActivities, newActivity)
		newActivity.Next = append(newActivity.Next, parent.ID)

		e.appState.SetProject(project)
		e.refreshTree()
		if e.onChanged != nil {
			e.onChanged()
		}
		if popup != nil {
			popup.Hide()
		}
	}

	form.OnCancel = func() {
		if popup != nil {
			popup.Hide()
		}
	}

	var canvas fyne.Canvas
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}

	popup = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Aggiungi Attività a: %s", parent.Name)),
			form,
		),
		canvas,
	)
	popup.Resize(fyne.NewSize(500, 400))
	popup.Show()
}

// showContextMenu mostra il menu contestuale per un'attività
func (e *ProjectEditor) showContextMenu(activity *lib.Activity, position fyne.Position) {
	if activity == nil {
		return
	}

	var canvas fyne.Canvas
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}

	// Crea i pulsanti del menu
	addBtn := widget.NewButton("➕ Aggiungi sub-attività", func() {
		e.showAddActivityDialog(activity)
		if e.contextPopup != nil {
			e.contextPopup.Hide()
		}
	})

	editBtn := widget.NewButton("✏️ Modifica", func() {
		e.openActivityEditor(activity)
		if e.contextPopup != nil {
			e.contextPopup.Hide()
		}
	})

	deleteBtn := widget.NewButton("🗑️ Elimina", func() {
		if e.contextPopup != nil {
			e.contextPopup.Hide()
		}
		e.deleteActivity(activity)
	})

	// Verifica se può essere spostato
	parent, index := e.findParentActivity(activity.ID)
	canMoveUp := parent != nil && index > 0
	canMoveDown := parent != nil && index < len(parent.SubActivities)-1

	moveUpBtn := widget.NewButton("↑ Sposta su", func() {
		if e.contextPopup != nil {
			e.contextPopup.Hide()
		}
		e.moveActivityUp(parent, index)
	})
	moveUpBtn.Disable()

	moveDownBtn := widget.NewButton("↓ Sposta giù", func() {
		if e.contextPopup != nil {
			e.contextPopup.Hide()
		}
		e.moveActivityDown(parent, index)
	})
	moveDownBtn.Disable()

	if canMoveUp {
		moveUpBtn.Enable()
	}
	if canMoveDown {
		moveDownBtn.Enable()
	}

	menuContent := container.NewVBox(
		addBtn,
		editBtn,
		widget.NewSeparator(),
		moveUpBtn,
		moveDownBtn,
		widget.NewSeparator(),
		deleteBtn,
	)

	e.contextPopup = widget.NewPopUp(menuContent, canvas)
	e.contextPopup.Move(position)
	e.contextPopup.Show()
}

// deleteActivity elimina un'attività con conferma
func (e *ProjectEditor) deleteActivity(activity *lib.Activity) {
	if activity == nil {
		return
	}

	// Conferma eliminazione
	dialog.ShowConfirm(
		"Conferma Eliminazione",
		fmt.Sprintf("Sei sicuro di voler eliminare l'attività '%s'?\nTutte le sub-attività verranno eliminate.", activity.Name),
		func(confirmed bool) {
			if confirmed {
				parent, _ := e.findParentActivity(activity.ID)
				if parent != nil {
					e.removeActivityFromParent(parent, activity.ID)
					e.appState.SetProject(e.appState.GetProject())
					e.refreshTree()
					if e.onChanged != nil {
						e.onChanged()
					}
				} else if e.appState.GetProject() != nil && e.appState.GetProject().Root != nil && e.appState.GetProject().Root.ID == activity.ID {
					// Tentativo di eliminare la root - non permesso
					dialog.ShowInformation("Errore", "Non è possibile eliminare la root activity", e.getWindow())
				}
			}
		},
		e.getWindow(),
	)
}

// findParentActivity trova il parent di un'attività e il suo indice
func (e *ProjectEditor) findParentActivity(childID string) (*lib.Activity, int) {
	project := e.appState.GetProject()
	if project == nil || project.Root == nil {
		return nil, -1
	}

	return e.findParentRecursive(project.Root, childID)
}

// findParentRecursive cerca ricorsivamente il parent
func (e *ProjectEditor) findParentRecursive(parent *lib.Activity, childID string) (*lib.Activity, int) {
	if parent == nil {
		return nil, -1
	}

	for i, sub := range parent.SubActivities {
		if sub.ID == childID {
			return parent, i
		}
		// Cerca nei figli
		if foundParent, foundIndex := e.findParentRecursive(sub, childID); foundParent != nil {
			return foundParent, foundIndex
		}
	}

	return nil, -1
}

// removeActivityFromParent rimuove un'attività dal parent
func (e *ProjectEditor) removeActivityFromParent(parent *lib.Activity, childID string) {
	if parent == nil {
		return
	}

	for i, sub := range parent.SubActivities {
		if sub.ID == childID {
			// Rimuovi dalla slice
			parent.SubActivities = append(parent.SubActivities[:i], parent.SubActivities[i+1:]...)
			return
		}
	}
}

// moveActivityUp sposta un'attività verso l'alto
func (e *ProjectEditor) moveActivityUp(parent *lib.Activity, index int) {
	if parent == nil || index <= 0 || index >= len(parent.SubActivities) {
		return
	}

	// Scambia con l'elemento precedente
	parent.SubActivities[index-1], parent.SubActivities[index] = parent.SubActivities[index], parent.SubActivities[index-1]

	e.appState.SetProject(e.appState.GetProject())
	e.refreshTree()
	if e.onChanged != nil {
		e.onChanged()
	}
}

// moveActivityDown sposta un'attività verso il basso
func (e *ProjectEditor) moveActivityDown(parent *lib.Activity, index int) {
	if parent == nil || index < 0 || index >= len(parent.SubActivities)-1 {
		return
	}

	// Scambia con l'elemento successivo
	parent.SubActivities[index], parent.SubActivities[index+1] = parent.SubActivities[index+1], parent.SubActivities[index]

	e.appState.SetProject(e.appState.GetProject())
	e.refreshTree()
	if e.onChanged != nil {
		e.onChanged()
	}
}

// getWindow ottiene la window corrente
func (e *ProjectEditor) getWindow() fyne.Window {
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		return windows[0]
	}
	return nil
}
