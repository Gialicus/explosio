package ui

import (
	"explosio/lib"
	"fmt"
	"strconv"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ActivityEditor gestisce l'editing di un'attività
type ActivityEditor struct {
	appState        *AppState
	activity        *lib.Activity
	onChanged       func()
	humansList      *widget.List
	materialsList   *widget.List
	assetsList      *widget.List
}

// NewActivityEditor crea un nuovo editor di attività
func NewActivityEditor(appState *AppState, activity *lib.Activity, onChanged func()) fyne.CanvasObject {
	editor := &ActivityEditor{
		appState:  appState,
		activity:  activity,
		onChanged: onChanged,
	}

	return editor.createEditor()
}

// createEditor crea l'interfaccia dell'editor
func (e *ActivityEditor) createEditor() fyne.CanvasObject {
	// Form base
	baseForm := e.createBaseForm()

	// Sezioni per risorse
	humansSection := e.createHumansSection()
	materialsSection := e.createMaterialsSection()
	assetsSection := e.createAssetsSection()

	// Sezione crashing
	crashingSection := e.createCrashingSection()

	// Crea tabs per organizzare meglio il contenuto
	tabs := container.NewAppTabs(
		container.NewTabItem("Informazioni", container.NewScroll(container.NewVBox(
			baseForm,
			widget.NewSeparator(),
			crashingSection,
		))),
		container.NewTabItem("Risorse Umane", container.NewScroll(humansSection)),
		container.NewTabItem("Materiali", container.NewScroll(materialsSection)),
		container.NewTabItem("Asset", container.NewScroll(assetsSection)),
	)

	return tabs
}

// createBaseForm crea il form base per nome, descrizione, durata
func (e *ActivityEditor) createBaseForm() fyne.CanvasObject {
	nameEntry := widget.NewEntry()
	descEntry := widget.NewMultiLineEntry()
	durationEntry := widget.NewEntry()

	if e.activity != nil {
		nameEntry.SetText(e.activity.Name)
		descEntry.SetText(e.activity.Description)
		durationEntry.SetText(fmt.Sprintf("%d", e.activity.Duration))
	}

	nameEntry.OnChanged = func(text string) {
		if e.activity != nil {
			e.activity.Name = text
			if e.onChanged != nil {
				e.onChanged()
			}
		}
	}

	descEntry.OnChanged = func(text string) {
		if e.activity != nil {
			e.activity.Description = text
			if e.onChanged != nil {
				e.onChanged()
			}
		}
	}

	durationEntry.OnChanged = func(text string) {
		if e.activity != nil {
			if duration, err := strconv.Atoi(text); err == nil && duration > 0 {
				e.activity.Duration = duration
				if e.activity.MinDuration > duration {
					e.activity.MinDuration = duration
				}
				if e.onChanged != nil {
					e.onChanged()
				}
			}
		}
	}

	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Nome", nameEntry),
			widget.NewFormItem("Descrizione", descEntry),
			widget.NewFormItem("Durata (min)", durationEntry),
		),
	)
}

// createHumansSection crea la sezione per le risorse umane
func (e *ActivityEditor) createHumansSection() fyne.CanvasObject {
	title := widget.NewRichTextFromMarkdown("## Risorse Umane")
	
	e.humansList = widget.NewList(
		func() int {
			if e.activity != nil {
				return len(e.activity.Humans)
			}
			return 0
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewButton("Modifica", nil),
				widget.NewButton("Rimuovi", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if e.activity != nil && id < len(e.activity.Humans) {
				h := e.activity.Humans[id]
				container := obj.(*fyne.Container)
				label := container.Objects[0].(*widget.Label)
				supplierInfo := ""
				if h.Supplier != nil {
					supplierInfo = fmt.Sprintf(" [Fornitore: %s]", h.Supplier.Name)
				}
				label.SetText(fmt.Sprintf("%s (€%.2f/h, qty: %.1f)%s", h.Role, h.CostPerH, h.Quantity, supplierInfo))
				
				modifyBtn := container.Objects[1].(*widget.Button)
				modifyBtn.OnTapped = func() {
					e.editHuman(id)
				}
				
				removeBtn := container.Objects[2].(*widget.Button)
				removeBtn.OnTapped = func() {
					e.removeHuman(id)
				}
			}
		},
	)

	addBtn := widget.NewButton("+ Aggiungi Risorsa Umana", func() {
		e.addHuman()
	})

	return container.NewBorder(
		title,
		addBtn,
		nil,
		nil,
		e.humansList,
	)
}

// createMaterialsSection crea la sezione per i materiali
func (e *ActivityEditor) createMaterialsSection() fyne.CanvasObject {
	title := widget.NewRichTextFromMarkdown("## Materiali")
	
	e.materialsList = widget.NewList(
		func() int {
			if e.activity != nil {
				return len(e.activity.Materials)
			}
			return 0
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewButton("Modifica", nil),
				widget.NewButton("Rimuovi", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if e.activity != nil && id < len(e.activity.Materials) {
				m := e.activity.Materials[id]
				container := obj.(*fyne.Container)
				label := container.Objects[0].(*widget.Label)
				supplierInfo := ""
				if m.Supplier != nil {
					supplierInfo = fmt.Sprintf(" [Fornitore: %s]", m.Supplier.Name)
				}
				label.SetText(fmt.Sprintf("%s (€%.2f/unit, qty: %.1f)%s", m.Name, m.UnitCost, m.Quantity, supplierInfo))
				
				modifyBtn := container.Objects[1].(*widget.Button)
				modifyBtn.OnTapped = func() {
					e.editMaterial(id)
				}
				
				removeBtn := container.Objects[2].(*widget.Button)
				removeBtn.OnTapped = func() {
					e.removeMaterial(id)
				}
			}
		},
	)

	addBtn := widget.NewButton("+ Aggiungi Materiale", func() {
		e.addMaterial()
	})

	return container.NewBorder(
		title,
		addBtn,
		nil,
		nil,
		e.materialsList,
	)
}

// createAssetsSection crea la sezione per gli asset
func (e *ActivityEditor) createAssetsSection() fyne.CanvasObject {
	title := widget.NewRichTextFromMarkdown("## Asset")
	
	e.assetsList = widget.NewList(
		func() int {
			if e.activity != nil {
				return len(e.activity.Assets)
			}
			return 0
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewButton("Modifica", nil),
				widget.NewButton("Rimuovi", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if e.activity != nil && id < len(e.activity.Assets) {
				as := e.activity.Assets[id]
				container := obj.(*fyne.Container)
				label := container.Objects[0].(*widget.Label)
				supplierInfo := ""
				if as.Supplier != nil {
					supplierInfo = fmt.Sprintf(" [Fornitore: %s]", as.Supplier.Name)
				}
				label.SetText(fmt.Sprintf("%s (€%.2f/use, qty: %.1f)%s", as.Name, as.CostPerUse, as.Quantity, supplierInfo))
				
				modifyBtn := container.Objects[1].(*widget.Button)
				modifyBtn.OnTapped = func() {
					e.editAsset(id)
				}
				
				removeBtn := container.Objects[2].(*widget.Button)
				removeBtn.OnTapped = func() {
					e.removeAsset(id)
				}
			}
		},
	)

	addBtn := widget.NewButton("+ Aggiungi Asset", func() {
		e.addAsset()
	})

	return container.NewBorder(
		title,
		addBtn,
		nil,
		nil,
		e.assetsList,
	)
}

// createCrashingSection crea la sezione per i parametri crashing
func (e *ActivityEditor) createCrashingSection() fyne.CanvasObject {
	title := widget.NewRichTextFromMarkdown("## Parametri Crashing")
	
	minDurationEntry := widget.NewEntry()
	crashCostEntry := widget.NewEntry()

	if e.activity != nil {
		minDurationEntry.SetText(fmt.Sprintf("%d", e.activity.MinDuration))
		crashCostEntry.SetText(fmt.Sprintf("%.2f", e.activity.CrashCostStep))
	}

	minDurationEntry.OnChanged = func(text string) {
		if e.activity != nil {
			if minDur, err := strconv.Atoi(text); err == nil && minDur >= 0 {
				if minDur <= e.activity.Duration {
					e.activity.MinDuration = minDur
					if e.onChanged != nil {
						e.onChanged()
					}
				}
			}
		}
	}

	crashCostEntry.OnChanged = func(text string) {
		if e.activity != nil {
			if cost, err := strconv.ParseFloat(text, 64); err == nil && cost >= 0 {
				e.activity.CrashCostStep = cost
				if e.onChanged != nil {
					e.onChanged()
				}
			}
		}
	}

	return container.NewVBox(
		title,
		widget.NewForm(
			widget.NewFormItem("Durata Minima (min)", minDurationEntry),
			widget.NewFormItem("Costo Extra Crashing (€)", crashCostEntry),
		),
	)
}

// Metodi helper per gestire risorse

func (e *ActivityEditor) addHuman() {
	if e.activity == nil {
		return
	}
	
	// Crea una nuova risorsa umana di default
	newHuman := lib.HumanResource{
		Role:        "Nuovo Ruolo",
		Description:  "",
		CostPerH:    15.0,
		Quantity:    1.0,
		Supplier:    nil,
	}
	
	index := len(e.activity.Humans)
	e.activity.Humans = append(e.activity.Humans, newHuman)
	e.showEditHumanDialog(&e.activity.Humans[index], index)
}

func (e *ActivityEditor) editHuman(index int) {
	if e.activity == nil || index >= len(e.activity.Humans) {
		return
	}
	
	human := &e.activity.Humans[index]
	e.showEditHumanDialog(human, index)
}

func (e *ActivityEditor) removeHuman(index int) {
	if e.activity == nil || index >= len(e.activity.Humans) {
		return
	}
	
	e.activity.Humans = append(e.activity.Humans[:index], e.activity.Humans[index+1:]...)
	if e.humansList != nil {
		e.humansList.Refresh()
	}
	if e.onChanged != nil {
		e.onChanged()
	}
}

func (e *ActivityEditor) addMaterial() {
	if e.activity == nil {
		return
	}
	
	// Crea una nuova risorsa materiale di default
	newMaterial := lib.MaterialResource{
		Name:        "Nuovo Materiale",
		Description: "",
		UnitCost:    1.0,
		Quantity:    1.0,
		Supplier:    nil,
	}
	
	index := len(e.activity.Materials)
	e.activity.Materials = append(e.activity.Materials, newMaterial)
	e.showEditMaterialDialog(&e.activity.Materials[index], index)
}

func (e *ActivityEditor) editMaterial(index int) {
	if e.activity == nil || index >= len(e.activity.Materials) {
		return
	}
	
	material := &e.activity.Materials[index]
	e.showEditMaterialDialog(material, index)
}

func (e *ActivityEditor) removeMaterial(index int) {
	if e.activity == nil || index >= len(e.activity.Materials) {
		return
	}
	
	e.activity.Materials = append(e.activity.Materials[:index], e.activity.Materials[index+1:]...)
	if e.materialsList != nil {
		e.materialsList.Refresh()
	}
	if e.onChanged != nil {
		e.onChanged()
	}
}

func (e *ActivityEditor) addAsset() {
	if e.activity == nil {
		return
	}
	
	// Crea un nuovo asset di default
	newAsset := lib.Asset{
		Name:        "Nuovo Asset",
		Description: "",
		CostPerUse:  1.0,
		Quantity:    1.0,
		Supplier:    nil,
	}
	
	index := len(e.activity.Assets)
	e.activity.Assets = append(e.activity.Assets, newAsset)
	e.showEditAssetDialog(&e.activity.Assets[index], index)
}

func (e *ActivityEditor) editAsset(index int) {
	if e.activity == nil || index >= len(e.activity.Assets) {
		return
	}
	
	asset := &e.activity.Assets[index]
	e.showEditAssetDialog(asset, index)
}

func (e *ActivityEditor) removeAsset(index int) {
	if e.activity == nil || index >= len(e.activity.Assets) {
		return
	}
	
	e.activity.Assets = append(e.activity.Assets[:index], e.activity.Assets[index+1:]...)
	if e.assetsList != nil {
		e.assetsList.Refresh()
	}
	if e.onChanged != nil {
		e.onChanged()
	}
}

// showEditHumanDialog mostra il dialog per modificare una risorsa umana
func (e *ActivityEditor) showEditHumanDialog(human *lib.HumanResource, index int) {
	roleEntry := widget.NewEntry()
	roleEntry.SetText(human.Role)
	
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(human.Description)
	
	costEntry := widget.NewEntry()
	costEntry.SetText(fmt.Sprintf("%.2f", human.CostPerH))
	
	quantityEntry := widget.NewEntry()
	quantityEntry.SetText(fmt.Sprintf("%.2f", human.Quantity))
	
	// Dropdown fornitori
	suppliers := e.appState.GetSuppliers()
	supplierOptions := []string{"Nessuno"}
	supplierMap := make(map[string]*lib.Supplier)
	supplierMap["Nessuno"] = nil
	
	for _, s := range suppliers {
		option := s.Name
		supplierOptions = append(supplierOptions, option)
		supplierMap[option] = s
	}
	
	supplierSelect := widget.NewSelect(supplierOptions, nil)
	if human.Supplier != nil {
		supplierSelect.SetSelected(human.Supplier.Name)
	} else {
		supplierSelect.SetSelected("Nessuno")
	}
	
	form := widget.NewForm(
		widget.NewFormItem("Ruolo", roleEntry),
		widget.NewFormItem("Descrizione", descEntry),
		widget.NewFormItem("Costo per Ora (€)", costEntry),
		widget.NewFormItem("Quantità", quantityEntry),
		widget.NewFormItem("Fornitore", supplierSelect),
	)
	
	var popup *widget.PopUp
	form.OnSubmit = func() {
		cost, err1 := strconv.ParseFloat(costEntry.Text, 64)
		quantity, err2 := strconv.ParseFloat(quantityEntry.Text, 64)
		
		if err1 != nil || err2 != nil || cost < 0 || quantity < 0 {
			// Errore di validazione - potremmo mostrare un messaggio
			return
		}
		
		// Aggiorna la risorsa
		human.Role = roleEntry.Text
		human.Description = descEntry.Text
		human.CostPerH = cost
		human.Quantity = quantity
		human.Supplier = supplierMap[supplierSelect.Selected]
		
		// Aggiorna l'UI
		if e.humansList != nil {
			e.humansList.Refresh()
		}
		if e.onChanged != nil {
			e.onChanged()
		}
		
		popup.Hide()
	}
	
	form.OnCancel = func() {
		popup.Hide()
	}
	
	// Ottieni canvas per il popup
	var canvas fyne.Canvas
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}
	
	popup = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Modifica Risorsa Umana"),
			form,
		),
		canvas,
	)
	popup.Resize(fyne.NewSize(500, 400))
	popup.Show()
}

// showEditMaterialDialog mostra il dialog per modificare un materiale
func (e *ActivityEditor) showEditMaterialDialog(material *lib.MaterialResource, index int) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(material.Name)
	
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(material.Description)
	
	costEntry := widget.NewEntry()
	costEntry.SetText(fmt.Sprintf("%.2f", material.UnitCost))
	
	quantityEntry := widget.NewEntry()
	quantityEntry.SetText(fmt.Sprintf("%.2f", material.Quantity))
	
	// Dropdown fornitori
	suppliers := e.appState.GetSuppliers()
	supplierOptions := []string{"Nessuno"}
	supplierMap := make(map[string]*lib.Supplier)
	supplierMap["Nessuno"] = nil
	
	for _, s := range suppliers {
		option := s.Name
		supplierOptions = append(supplierOptions, option)
		supplierMap[option] = s
	}
	
	supplierSelect := widget.NewSelect(supplierOptions, nil)
	if material.Supplier != nil {
		supplierSelect.SetSelected(material.Supplier.Name)
	} else {
		supplierSelect.SetSelected("Nessuno")
	}
	
	form := widget.NewForm(
		widget.NewFormItem("Nome", nameEntry),
		widget.NewFormItem("Descrizione", descEntry),
		widget.NewFormItem("Costo Unitario (€)", costEntry),
		widget.NewFormItem("Quantità", quantityEntry),
		widget.NewFormItem("Fornitore", supplierSelect),
	)
	
	var popup *widget.PopUp
	form.OnSubmit = func() {
		cost, err1 := strconv.ParseFloat(costEntry.Text, 64)
		quantity, err2 := strconv.ParseFloat(quantityEntry.Text, 64)
		
		if err1 != nil || err2 != nil || cost < 0 || quantity < 0 {
			// Errore di validazione
			return
		}
		
		// Aggiorna il materiale
		material.Name = nameEntry.Text
		material.Description = descEntry.Text
		material.UnitCost = cost
		material.Quantity = quantity
		material.Supplier = supplierMap[supplierSelect.Selected]
		
		// Aggiorna l'UI
		if e.materialsList != nil {
			e.materialsList.Refresh()
		}
		if e.onChanged != nil {
			e.onChanged()
		}
		
		popup.Hide()
	}
	
	form.OnCancel = func() {
		popup.Hide()
	}
	
	// Ottieni canvas per il popup
	var canvas fyne.Canvas
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}
	
	popup = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Modifica Materiale"),
			form,
		),
		canvas,
	)
	popup.Resize(fyne.NewSize(500, 400))
	popup.Show()
}

// showEditAssetDialog mostra il dialog per modificare un asset
func (e *ActivityEditor) showEditAssetDialog(asset *lib.Asset, index int) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(asset.Name)
	
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(asset.Description)
	
	costEntry := widget.NewEntry()
	costEntry.SetText(fmt.Sprintf("%.2f", asset.CostPerUse))
	
	quantityEntry := widget.NewEntry()
	quantityEntry.SetText(fmt.Sprintf("%.2f", asset.Quantity))
	
	// Dropdown fornitori
	suppliers := e.appState.GetSuppliers()
	supplierOptions := []string{"Nessuno"}
	supplierMap := make(map[string]*lib.Supplier)
	supplierMap["Nessuno"] = nil
	
	for _, s := range suppliers {
		option := s.Name
		supplierOptions = append(supplierOptions, option)
		supplierMap[option] = s
	}
	
	supplierSelect := widget.NewSelect(supplierOptions, nil)
	if asset.Supplier != nil {
		supplierSelect.SetSelected(asset.Supplier.Name)
	} else {
		supplierSelect.SetSelected("Nessuno")
	}
	
	form := widget.NewForm(
		widget.NewFormItem("Nome", nameEntry),
		widget.NewFormItem("Descrizione", descEntry),
		widget.NewFormItem("Costo per Uso (€)", costEntry),
		widget.NewFormItem("Quantità", quantityEntry),
		widget.NewFormItem("Fornitore", supplierSelect),
	)
	
	var popup *widget.PopUp
	form.OnSubmit = func() {
		cost, err1 := strconv.ParseFloat(costEntry.Text, 64)
		quantity, err2 := strconv.ParseFloat(quantityEntry.Text, 64)
		
		if err1 != nil || err2 != nil || cost < 0 || quantity < 0 {
			// Errore di validazione
			return
		}
		
		// Aggiorna l'asset
		asset.Name = nameEntry.Text
		asset.Description = descEntry.Text
		asset.CostPerUse = cost
		asset.Quantity = quantity
		asset.Supplier = supplierMap[supplierSelect.Selected]
		
		// Aggiorna l'UI
		if e.assetsList != nil {
			e.assetsList.Refresh()
		}
		if e.onChanged != nil {
			e.onChanged()
		}
		
		popup.Hide()
	}
	
	form.OnCancel = func() {
		popup.Hide()
	}
	
	// Ottieni canvas per il popup
	var canvas fyne.Canvas
	if windows := fyne.CurrentApp().Driver().AllWindows(); len(windows) > 0 {
		canvas = windows[0].Canvas()
	} else {
		return
	}
	
	popup = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Modifica Asset"),
			form,
		),
		canvas,
	)
	popup.Resize(fyne.NewSize(500, 400))
	popup.Show()
}
