package gui

import (
	"explosio/core"
	"explosio/core/unit"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Run starts the GUI with the given root activity. If root is nil, creates a minimal default project.
func Run(root *core.Activity) {
	if root == nil {
		root = core.NewActivity("Nuovo progetto", "", unit.Duration{Value: 0, Unit: unit.DurationUnitDay}, unit.Price{Value: 0, Currency: "EUR"})
	}

	w := fyne.CurrentApp().NewWindow("Explosio - Activity tree")
	w.Resize(fyne.NewSize(900, 600))

	var tree *widget.Tree
	var form *ActivityForm
	var proj *core.Project = core.NewProject(root)

	refreshTree := func() {
		if tree != nil {
			tree.Refresh()
		}
	}

	tree = NewActivityTree(root, func(selected *core.Activity) {
		if form != nil {
			form.SelectActivity(selected)
		}
	})
	tree.OpenAllBranches()

	form = NewActivityForm(root, refreshTree, w)
	form.SetWindow(w)
	// Seleziona "Install pipes" (path 1-0) se esiste, così si vedono subito materiali e risorse
	if installPipes := pathToActivity(root, "1-0"); installPipes != nil {
		tree.Select("1-0")
		form.SelectActivity(installPipes)
	} else {
		form.SelectActivity(root)
	}

	split := container.NewHSplit(
		container.NewBorder(nil, nil, nil, nil, tree),
		form.Content(),
	)
	split.SetOffset(0.3)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentIcon(), func() {
			dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
				if err != nil || uc == nil {
					return
				}
				defer uc.Close()
				path := uc.URI().Path()
				var loaded *core.Project
				if strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml") {
					loaded, err = core.ReadYAML(uc)
				} else {
					loaded, err = core.ReadJSON(uc)
				}
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				proj = loaded
				root = loaded.Root
				tree = NewActivityTree(root, func(selected *core.Activity) {
					if form != nil {
						form.SelectActivity(selected)
					}
				})
				tree.OpenAllBranches()
				form = NewActivityForm(root, refreshTree, w)
				form.SetWindow(w)
				form.SelectActivity(root)
				split.Leading = tree
				split.Trailing = form.Content()
				split.Refresh()
			}, w)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			r := root.Validate()
			if !r.Valid() {
				var msg strings.Builder
				for _, e := range r.Errors {
					msg.WriteString("• ")
					msg.WriteString(e.Error())
					msg.WriteString("\n")
				}
				for _, e := range r.Warnings {
					msg.WriteString("⚠ ")
					msg.WriteString(e.Error())
					msg.WriteString("\n")
				}
				dialog.ShowInformation("Validazione", "Errori prima del salvataggio:\n\n"+msg.String(), w)
				return
			}
			dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
				if err != nil || uc == nil {
					return
				}
				defer uc.Close()
				proj = core.NewProject(root)
				if err := proj.WriteJSON(uc); err != nil {
					dialog.ShowError(err, w)
					return
				}
			}, w)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			curr := form.Current()
			if curr == nil {
				curr = root
			}
			child := core.NewActivity("Nuova attività", "", unit.Duration{Value: 0, Unit: unit.DurationUnitDay}, unit.Price{Value: 0, Currency: "EUR"})
			curr.AddActivity(child)
			tree.Refresh()
			form.SelectActivity(child)
		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			curr := form.Current()
			if curr == nil || curr == root {
				dialog.ShowInformation("Elimina", "Non è possibile eliminare la radice del progetto.", w)
				return
			}
			parent, idx := activityParent(root, curr)
			if parent == nil {
				return
			}
			parent.Activities = append(parent.Activities[:idx], parent.Activities[idx+1:]...)
			tree.Refresh()
			form.SelectActivity(parent)
		}),
	)

	statusLabel := widget.NewLabel("")
	updateStatus := func() {
		totalPrice := root.CalculatePrice()
		totalDur := root.CalculateDuration()
		statusLabel.SetText(fmt.Sprintf("Totale: %.0f %s | Durata: %.0f %s", totalPrice, root.Price.Currency, totalDur, root.Duration.Unit))
	}
	updateStatus()
	form.onRefresh = func() {
		refreshTree()
		updateStatus()
	}

	content := container.NewBorder(
		container.NewVBox(toolbar, widget.NewSeparator()),
		container.NewVBox(widget.NewSeparator(), statusLabel),
		nil, nil,
		split,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
