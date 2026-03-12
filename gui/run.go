package gui

import (
	"explosio/core"
	"explosio/gui/app"
	"explosio/gui/components"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Run starts the GUI with the given root activity. If root is nil, creates a minimal default project.
func Run(root *core.Activity) {
	projSvc := app.NewProjectService()
	actSvc := app.NewActivityService()
	matSvc := app.NewMaterialsService()

	if root == nil {
		root = projSvc.NewProject().Root
	}

	w := fyne.CurrentApp().NewWindow("Explosio - Activity tree")
	w.Resize(fyne.NewSize(900, 600))

	var tree *widget.Tree
	var form *ActivityForm
	proj := core.NewProject(root)

	refreshTree := func() {
		if tree != nil {
			tree.Refresh()
		}
	}

	statusBar := components.NewStatusBar()
	updateStatus := func() {
		statusBar.Update(projSvc.GetSummary(root))
	}

	tree = components.NewActivityTree(root, actSvc, func(selected *core.Activity) {
		if form != nil {
			form.SelectActivity(selected)
		}
	})
	tree.OpenAllBranches()

	form = NewActivityForm(root, actSvc, matSvc, func() {
		refreshTree()
		updateStatus()
	}, w)
	form.SetWindow(w)

	if installPipes := actSvc.PathToActivity(root, "1-0"); installPipes != nil {
		tree.Select(components.TreeSelectActivityPath("1-0"))
		form.SelectActivity(installPipes)
	} else {
		form.SelectActivity(root)
	}

	split := container.NewHSplit(tree, form.Content())
	split.SetOffset(0.3)

	toolbar := components.NewToolbar(components.ToolbarCallbacks{
		OnOpen: func() {
			dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
				if err != nil || uc == nil {
					return
				}
				defer uc.Close()
				path := uc.URI().Path()
				format := "json"
				if strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml") {
					format = "yaml"
				}
				loaded, err := projSvc.LoadProject(uc, format)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				proj = loaded
				root = loaded.Root
				tree = components.NewActivityTree(root, actSvc, func(selected *core.Activity) {
					if form != nil {
						form.SelectActivity(selected)
					}
				})
				tree.OpenAllBranches()
				form = NewActivityForm(root, actSvc, matSvc, func() {
					refreshTree()
					updateStatus()
				}, w)
				form.SetWindow(w)
				form.SelectActivity(root)
				split.Leading = tree
				split.Trailing = form.Content()
				split.Refresh()
				updateStatus()
			}, w)
		},
		OnSave: func() {
			r := projSvc.Validate(proj)
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
				if err := projSvc.SaveProject(proj, uc); err != nil {
					dialog.ShowError(err, w)
					return
				}
			}, w)
		},
		OnAddActivity: func() {
			curr := form.Current()
			if curr == nil {
				curr = root
			}
			child := actSvc.AddChild(curr, "Nuova attività")
			tree.Refresh()
			form.SelectActivity(child)
			updateStatus()
		},
		OnDeleteActivity: func() {
			curr := form.Current()
			if curr == nil || curr == root {
				dialog.ShowInformation("Elimina", "Non è possibile eliminare la radice del progetto.", w)
				return
			}
			dialog.ShowConfirm("Elimina attività", "Eliminare \""+curr.Name+"\" e tutte le sottostrutture?", func(ok bool) {
				if !ok {
					return
				}
				parent, _ := actSvc.GetParent(root, curr)
				if actSvc.DeleteActivity(root, curr) {
					tree.Refresh()
					form.SelectActivity(parent)
					updateStatus()
				}
			}, w)
		},
	})

	updateStatus()

	content := container.NewBorder(
		container.NewVBox(toolbar, widget.NewSeparator()),
		statusBar.Content(),
		nil, nil,
		split,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
