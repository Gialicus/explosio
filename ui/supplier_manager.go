package ui

import (
	"explosio/lib"
	"fmt"
	"strconv"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SupplierManager gestisce la lista dei fornitori
type SupplierManager struct {
	appState  *AppState
	onChanged func()
	list      *widget.List
}

// NewSupplierManager crea un nuovo gestore di fornitori
func NewSupplierManager(appState *AppState, onChanged func()) fyne.CanvasObject {
	manager := &SupplierManager{
		appState:  appState,
		onChanged: onChanged,
	}

	// Lista fornitori
	manager.list = widget.NewList(
		func() int {
			return len(appState.GetSuppliers())
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewButton("Modifica", nil),
				widget.NewButton("Rimuovi", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			suppliers := appState.GetSuppliers()
			if id < len(suppliers) {
				s := suppliers[id]
				container := obj.(*fyne.Container)
				label := container.Objects[0].(*widget.Label)
				label.SetText(fmt.Sprintf("%s - %s (€%.2f, %.1f/%s)", 
					s.Name, s.Description, s.UnitCost, s.AvailableQuantity, s.Period.String()))
				
				modifyBtn := container.Objects[1].(*widget.Button)
				modifyBtn.OnTapped = func() {
					manager.editSupplier(s)
				}
				
				removeBtn := container.Objects[2].(*widget.Button)
				removeBtn.OnTapped = func() {
					manager.removeSupplier(s.Name)
				}
			}
		},
	)

	// Pulsante per aggiungere fornitore
	addBtn := widget.NewButton("Aggiungi Fornitore", func() {
		manager.addSupplier()
	})

	// Layout
	content := container.NewBorder(
		nil,
		addBtn,
		nil,
		nil,
		manager.list,
	)

	// Aggiorna la lista quando cambiano i fornitori
	appState.OnSuppliersChanged(func([]*lib.Supplier) {
		manager.list.Refresh()
	})

	return content
}

// addSupplier aggiunge un nuovo fornitore
func (m *SupplierManager) addSupplier() {
	m.showSupplierDialog(nil)
}

// editSupplier modifica un fornitore esistente
func (m *SupplierManager) editSupplier(supplier *lib.Supplier) {
	m.showSupplierDialog(supplier)
}

// removeSupplier rimuove un fornitore
func (m *SupplierManager) removeSupplier(name string) {
	if m.appState.RemoveSupplier(name) {
		if m.onChanged != nil {
			m.onChanged()
		}
	}
}

// showSupplierDialog mostra il dialog per creare/modificare un fornitore
func (m *SupplierManager) showSupplierDialog(supplier *lib.Supplier) {
	nameEntry := widget.NewEntry()
	descEntry := widget.NewMultiLineEntry()
	unitCostEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()
	periodSelect := widget.NewSelect([]string{
		string(lib.PeriodMinute),
		string(lib.PeriodHour),
		string(lib.PeriodDay),
		string(lib.PeriodWeek),
		string(lib.PeriodMonth),
		string(lib.PeriodYear),
	}, nil)

	if supplier != nil {
		nameEntry.SetText(supplier.Name)
		descEntry.SetText(supplier.Description)
		unitCostEntry.SetText(fmt.Sprintf("%.2f", supplier.UnitCost))
		quantityEntry.SetText(fmt.Sprintf("%.2f", supplier.AvailableQuantity))
		periodSelect.SetSelected(string(supplier.Period))
	} else {
		periodSelect.SetSelected(string(lib.PeriodMonth))
	}

	form := widget.NewForm(
		widget.NewFormItem("Nome", nameEntry),
		widget.NewFormItem("Descrizione", descEntry),
		widget.NewFormItem("Costo Unitario (€)", unitCostEntry),
		widget.NewFormItem("Quantità Disponibile", quantityEntry),
		widget.NewFormItem("Periodo", periodSelect),
	)

	var dialog *widget.PopUp
	form.OnSubmit = func() {
		unitCost, err1 := strconv.ParseFloat(unitCostEntry.Text, 64)
		quantity, err2 := strconv.ParseFloat(quantityEntry.Text, 64)
		
		if err1 != nil || err2 != nil || unitCost < 0 || quantity < 0 {
			// Mostra errore
			return
		}

		period := lib.PeriodType(periodSelect.Selected)
		if !period.IsValid() {
			return
		}

		if supplier != nil {
			// Modifica fornitore esistente
			supplier.Name = nameEntry.Text
			supplier.Description = descEntry.Text
			supplier.UnitCost = unitCost
			supplier.AvailableQuantity = quantity
			supplier.Period = period
		} else {
			// Crea nuovo fornitore
			newSupplier := lib.NewSupplier(
				nameEntry.Text,
				descEntry.Text,
				unitCost,
				quantity,
				period,
			)
			m.appState.AddSupplier(newSupplier)
		}

		if m.onChanged != nil {
			m.onChanged()
		}

		dialog.Hide()
	}

	form.OnCancel = func() {
		dialog.Hide()
	}

	// Ottieni la window corrente per il dialog
	var canvas fyne.Canvas
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}

	dialog = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Fornitore"),
			form,
		),
		canvas,
	)
	dialog.Resize(fyne.NewSize(400, 300))
	dialog.Show()
}
