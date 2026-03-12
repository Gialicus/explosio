// Package components provides reusable GUI components for Explosio.
package components

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ToolbarCallbacks holds the callbacks for toolbar actions.
type ToolbarCallbacks struct {
	OnOpen           func()
	OnSave           func()
	OnAddActivity    func()
	OnDeleteActivity func()
}

// NewToolbar creates a toolbar with the given callbacks.
func NewToolbar(callbacks ToolbarCallbacks) *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentIcon(), func() {
			if callbacks.OnOpen != nil {
				callbacks.OnOpen()
			}
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			if callbacks.OnSave != nil {
				callbacks.OnSave()
			}
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			if callbacks.OnAddActivity != nil {
				callbacks.OnAddActivity()
			}
		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			if callbacks.OnDeleteActivity != nil {
				callbacks.OnDeleteActivity()
			}
		}),
	)
}
