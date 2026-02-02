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
	appState  *AppState
	activity  *lib.Activity
	onChanged func()
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

	// Layout verticale con scroll
	content := container.NewVBox(
		baseForm,
		widget.NewSeparator(),
		humansSection,
		widget.NewSeparator(),
		materialsSection,
		widget.NewSeparator(),
		assetsSection,
		widget.NewSeparator(),
		crashingSection,
	)

	return container.NewScroll(content)
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
	
	list := widget.NewList(
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
				label.SetText(fmt.Sprintf("%s (€%.2f/h, qty: %.1f)", h.Role, h.CostPerH, h.Quantity))
				
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

	addBtn := widget.NewButton("Aggiungi Risorsa Umana", func() {
		e.addHuman()
	})

	return container.NewVBox(
		title,
		list,
		addBtn,
	)
}

// createMaterialsSection crea la sezione per i materiali
func (e *ActivityEditor) createMaterialsSection() fyne.CanvasObject {
	title := widget.NewRichTextFromMarkdown("## Materiali")
	
	list := widget.NewList(
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
				label.SetText(fmt.Sprintf("%s (€%.2f/unit, qty: %.1f)", m.Name, m.UnitCost, m.Quantity))
				
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

	addBtn := widget.NewButton("Aggiungi Materiale", func() {
		e.addMaterial()
	})

	return container.NewVBox(
		title,
		list,
		addBtn,
	)
}

// createAssetsSection crea la sezione per gli asset
func (e *ActivityEditor) createAssetsSection() fyne.CanvasObject {
	title := widget.NewRichTextFromMarkdown("## Asset")
	
	list := widget.NewList(
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
				label.SetText(fmt.Sprintf("%s (€%.2f/use, qty: %.1f)", as.Name, as.CostPerUse, as.Quantity))
				
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

	addBtn := widget.NewButton("Aggiungi Asset", func() {
		e.addAsset()
	})

	return container.NewVBox(
		title,
		list,
		addBtn,
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
	
	e.activity.Humans = append(e.activity.Humans, newHuman)
	if e.onChanged != nil {
		e.onChanged()
	}
}

func (e *ActivityEditor) editHuman(index int) {
	if e.activity == nil || index >= len(e.activity.Humans) {
		return
	}
	
	// Per semplicità, per ora non implementiamo un dialog di editing
	// In una versione completa, apriremmo un dialog con form
}

func (e *ActivityEditor) removeHuman(index int) {
	if e.activity == nil || index >= len(e.activity.Humans) {
		return
	}
	
	e.activity.Humans = append(e.activity.Humans[:index], e.activity.Humans[index+1:]...)
	if e.onChanged != nil {
		e.onChanged()
	}
}

func (e *ActivityEditor) addMaterial() {
	if e.activity == nil {
		return
	}
	
	newMaterial := lib.MaterialResource{
		Name:        "Nuovo Materiale",
		Description: "",
		UnitCost:    1.0,
		Quantity:    1.0,
		Supplier:    nil,
	}
	
	e.activity.Materials = append(e.activity.Materials, newMaterial)
	if e.onChanged != nil {
		e.onChanged()
	}
}

func (e *ActivityEditor) editMaterial(index int) {
	// Placeholder
}

func (e *ActivityEditor) removeMaterial(index int) {
	if e.activity == nil || index >= len(e.activity.Materials) {
		return
	}
	
	e.activity.Materials = append(e.activity.Materials[:index], e.activity.Materials[index+1:]...)
	if e.onChanged != nil {
		e.onChanged()
	}
}

func (e *ActivityEditor) addAsset() {
	if e.activity == nil {
		return
	}
	
	newAsset := lib.Asset{
		Name:        "Nuovo Asset",
		Description: "",
		CostPerUse:  1.0,
		Quantity:    1.0,
		Supplier:    nil,
	}
	
	e.activity.Assets = append(e.activity.Assets, newAsset)
	if e.onChanged != nil {
		e.onChanged()
	}
}

func (e *ActivityEditor) editAsset(index int) {
	// Placeholder
}

func (e *ActivityEditor) removeAsset(index int) {
	if e.activity == nil || index >= len(e.activity.Assets) {
		return
	}
	
	e.activity.Assets = append(e.activity.Assets[:index], e.activity.Assets[index+1:]...)
	if e.onChanged != nil {
		e.onChanged()
	}
}
