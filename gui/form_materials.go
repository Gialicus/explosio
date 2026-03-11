package gui

import (
	"explosio/core"
	"explosio/core/material"
	"explosio/core/resource/asset"
	"explosio/core/resource/human"
	"explosio/core/unit"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// materialsAccordion wraps an Accordion for materials/resources with setActivity.
type materialsAccordion struct {
	accordion       *widget.Accordion
	activity        *core.Activity
	form            *ActivityForm
	complexList     *widget.List
	countableList   *widget.List
	measurableList  *widget.List
	humanList       *widget.List
	assetList       *widget.List
}

func newMaterialsAccordion(form *ActivityForm) *materialsAccordion {
	m := &materialsAccordion{form: form}
	m.accordion = widget.NewAccordion(
		widget.NewAccordionItem("Materiali complessi", m.buildComplexPanel()),
		widget.NewAccordionItem("Materiali numerabili", m.buildCountablePanel()),
		widget.NewAccordionItem("Materiali misurabili", m.buildMeasurablePanel()),
		widget.NewAccordionItem("Risorse umane", m.buildHumanPanel()),
		widget.NewAccordionItem("Asset", m.buildAssetPanel()),
	)
	m.accordion.MultiOpen = true
	m.accordion.OpenAll()
	return m
}

func (m *materialsAccordion) setActivity(a *core.Activity) {
	m.activity = a
	m.refreshPanels()
}

func (m *materialsAccordion) refreshPanels() {
	// Rebuild panel content - the accordion items hold static content, we need to refresh the lists
	// For simplicity, we'll rebuild the accordion items when activity changes
	if m.activity == nil {
		m.accordion.Items[0].Detail = m.buildComplexPanel()
		m.accordion.Items[1].Detail = m.buildCountablePanel()
		m.accordion.Items[2].Detail = m.buildMeasurablePanel()
		m.accordion.Items[3].Detail = m.buildHumanPanel()
		m.accordion.Items[4].Detail = m.buildAssetPanel()
	} else {
		m.accordion.Items[0].Detail = m.buildComplexPanel()
		m.accordion.Items[1].Detail = m.buildCountablePanel()
		m.accordion.Items[2].Detail = m.buildMeasurablePanel()
		m.accordion.Items[3].Detail = m.buildHumanPanel()
		m.accordion.Items[4].Detail = m.buildAssetPanel()
	}
	m.accordion.Refresh()
}

func (m *materialsAccordion) buildComplexPanel() fyne.CanvasObject {
	m.complexList = widget.NewList(
		func() int {
			if m.activity == nil {
				return 0
			}
			return len(m.activity.ComplexMaterials)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), widget.NewButton("Elimina", nil))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			btn := box.Objects[1].(*widget.Button)
			if m.activity != nil && id < len(m.activity.ComplexMaterials) {
				mat := m.activity.ComplexMaterials[id]
				label.SetText(fmt.Sprintf("%s (x%d) - %.2f %s", mat.Name, mat.UnitQuantity, mat.CalculatePrice(), mat.Price.Currency))
				btn.OnTapped = func() {
					m.activity.ComplexMaterials = append(m.activity.ComplexMaterials[:id], m.activity.ComplexMaterials[id+1:]...)
					m.refreshLists()
					if m.form.onRefresh != nil {
						m.form.onRefresh()
					}
				}
			}
		},
	)
	addBtn := widget.NewButton("+ Aggiungi", func() {
		m.showComplexMaterialDialog(nil)
	})
	return container.NewBorder(nil, addBtn, nil, nil, m.complexList)
}

func (m *materialsAccordion) buildCountablePanel() fyne.CanvasObject {
	m.countableList = widget.NewList(
		func() int {
			if m.activity == nil {
				return 0
			}
			return len(m.activity.CountableMaterials)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), widget.NewButton("Elimina", nil))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			btn := box.Objects[1].(*widget.Button)
			if m.activity != nil && id < len(m.activity.CountableMaterials) {
				mat := m.activity.CountableMaterials[id]
				label.SetText(fmt.Sprintf("%s (x%d) - %.2f %s", mat.Name, mat.Quantity, mat.CalculatePrice(), mat.Price.Currency))
				btn.OnTapped = func() {
					m.activity.CountableMaterials = append(m.activity.CountableMaterials[:id], m.activity.CountableMaterials[id+1:]...)
					m.refreshLists()
					if m.form.onRefresh != nil {
						m.form.onRefresh()
					}
				}
			}
		},
	)
	addBtn := widget.NewButton("+ Aggiungi", func() {
		m.showCountableMaterialDialog(nil)
	})
	return container.NewBorder(nil, addBtn, nil, nil, m.countableList)
}

func (m *materialsAccordion) buildMeasurablePanel() fyne.CanvasObject {
	m.measurableList = widget.NewList(
		func() int {
			if m.activity == nil {
				return 0
			}
			return len(m.activity.MeasurableMaterials)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), widget.NewButton("Elimina", nil))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			btn := box.Objects[1].(*widget.Button)
			if m.activity != nil && id < len(m.activity.MeasurableMaterials) {
				mat := m.activity.MeasurableMaterials[id]
				label.SetText(fmt.Sprintf("%s (%.0f %s) - %.2f %s", mat.Name, mat.Quantity.Value, mat.Quantity.Unit, mat.CalculatePrice(), mat.Price.Currency))
				btn.OnTapped = func() {
					m.activity.MeasurableMaterials = append(m.activity.MeasurableMaterials[:id], m.activity.MeasurableMaterials[id+1:]...)
					m.refreshLists()
					if m.form.onRefresh != nil {
						m.form.onRefresh()
					}
				}
			}
		},
	)
	addBtn := widget.NewButton("+ Aggiungi", func() {
		m.showMeasurableMaterialDialog(nil)
	})
	return container.NewBorder(nil, addBtn, nil, nil, m.measurableList)
}

func (m *materialsAccordion) buildHumanPanel() fyne.CanvasObject {
	m.humanList = widget.NewList(
		func() int {
			if m.activity == nil {
				return 0
			}
			return len(m.activity.HumanResources)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), widget.NewButton("Elimina", nil))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			btn := box.Objects[1].(*widget.Button)
			if m.activity != nil && id < len(m.activity.HumanResources) {
				hr := m.activity.HumanResources[id]
				label.SetText(fmt.Sprintf("%s - %.0f %s - %.2f %s", hr.Name, hr.Duration.Value, hr.Duration.Unit, hr.Price.Value, hr.Price.Currency))
				btn.OnTapped = func() {
					m.activity.HumanResources = append(m.activity.HumanResources[:id], m.activity.HumanResources[id+1:]...)
					m.refreshLists()
					if m.form.onRefresh != nil {
						m.form.onRefresh()
					}
				}
			}
		},
	)
	addBtn := widget.NewButton("+ Aggiungi", func() {
		m.showHumanResourceDialog(nil)
	})
	return container.NewBorder(nil, addBtn, nil, nil, m.humanList)
}

func (m *materialsAccordion) buildAssetPanel() fyne.CanvasObject {
	m.assetList = widget.NewList(
		func() int {
			if m.activity == nil {
				return 0
			}
			return len(m.activity.Assets)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), widget.NewButton("Elimina", nil))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			btn := box.Objects[1].(*widget.Button)
			if m.activity != nil && id < len(m.activity.Assets) {
				as := m.activity.Assets[id]
				label.SetText(fmt.Sprintf("%s - %.2f %s", as.Name, as.Price.Value, as.Price.Currency))
				btn.OnTapped = func() {
					m.activity.Assets = append(m.activity.Assets[:id], m.activity.Assets[id+1:]...)
					m.refreshLists()
					if m.form.onRefresh != nil {
						m.form.onRefresh()
					}
				}
			}
		},
	)
	addBtn := widget.NewButton("+ Aggiungi", func() {
		m.showAssetDialog(nil)
	})
	return container.NewBorder(nil, addBtn, nil, nil, m.assetList)
}

func (m *materialsAccordion) refreshLists() {
	if m.complexList != nil {
		m.complexList.Refresh()
	}
	if m.countableList != nil {
		m.countableList.Refresh()
	}
	if m.measurableList != nil {
		m.measurableList.Refresh()
	}
	if m.humanList != nil {
		m.humanList.Refresh()
	}
	if m.assetList != nil {
		m.assetList.Refresh()
	}
}

func (m *materialsAccordion) showComplexMaterialDialog(existing *material.ComplexMaterial) {
	if m.activity == nil {
		return
	}
	nameE := widget.NewEntry()
	descE := widget.NewEntry()
	priceE := widget.NewEntry()
	priceE.SetPlaceHolder("0")
	currE := widget.NewEntry()
	currE.SetText("EUR")
	qtyE := widget.NewEntry()
	qtyE.SetPlaceHolder("1")
	if existing != nil {
		nameE.SetText(existing.Name)
		descE.SetText(existing.Description)
		priceE.SetText(strconv.FormatFloat(existing.Price.Value, 'f', -1, 64))
		currE.SetText(existing.Price.Currency)
		qtyE.SetText(strconv.Itoa(existing.UnitQuantity))
	}
	items := []*widget.FormItem{
		widget.NewFormItem("Nome", nameE),
		widget.NewFormItem("Descrizione", descE),
		widget.NewFormItem("Prezzo", container.NewHBox(priceE, currE)),
		widget.NewFormItem("Quantità unità", qtyE),
	}
	callback := func(ok bool) {
		if !ok {
			return
		}
		priceVal, _ := strconv.ParseFloat(priceE.Text, 64)
		qty, _ := strconv.Atoi(qtyE.Text)
		if qty <= 0 {
			qty = 1
		}
		curr := currE.Text
		if curr == "" {
			curr = "EUR"
		}
		p := unit.Price{Value: priceVal, Currency: curr}
		meas := material.NewMeasurableMaterial("", "", unit.Price{Value: 0, Currency: curr}, unit.MeasurableQuantity{Value: 1, Unit: unit.UnitMeter})
		cm := material.NewComplexMaterial(nameE.Text, descE.Text, p, qty, meas)
		if existing != nil {
			for i, c := range m.activity.ComplexMaterials {
				if c == existing {
					m.activity.ComplexMaterials[i] = cm
					break
				}
			}
		} else {
			m.activity.AddComplexMaterial(cm)
		}
		m.refreshLists()
		if m.form.onRefresh != nil {
			m.form.onRefresh()
		}
	}
	showFormDialog("Materiale complesso", items, callback, m.form.win)
}

func (m *materialsAccordion) showCountableMaterialDialog(existing *material.CountableMaterial) {
	if m.activity == nil {
		return
	}
	nameE := widget.NewEntry()
	descE := widget.NewEntry()
	priceE := widget.NewEntry()
	priceE.SetPlaceHolder("0")
	currE := widget.NewEntry()
	currE.SetText("EUR")
	qtyE := widget.NewEntry()
	qtyE.SetPlaceHolder("1")
	if existing != nil {
		nameE.SetText(existing.Name)
		descE.SetText(existing.Description)
		priceE.SetText(strconv.FormatFloat(existing.Price.Value, 'f', -1, 64))
		currE.SetText(existing.Price.Currency)
		qtyE.SetText(strconv.Itoa(existing.Quantity))
	}
	items := []*widget.FormItem{
		widget.NewFormItem("Nome", nameE),
		widget.NewFormItem("Descrizione", descE),
		widget.NewFormItem("Prezzo", container.NewHBox(priceE, currE)),
		widget.NewFormItem("Quantità", qtyE),
	}
	callback := func(ok bool) {
		if !ok {
			return
		}
		priceVal, _ := strconv.ParseFloat(priceE.Text, 64)
		qty, _ := strconv.Atoi(qtyE.Text)
		if qty < 0 {
			qty = 0
		}
		curr := currE.Text
		if curr == "" {
			curr = "EUR"
		}
		p := unit.Price{Value: priceVal, Currency: curr}
		cm := material.NewCountableMaterial(nameE.Text, descE.Text, p, qty)
		if existing != nil {
			for i, c := range m.activity.CountableMaterials {
				if c == existing {
					m.activity.CountableMaterials[i] = cm
					break
				}
			}
		} else {
			m.activity.AddCountableMaterial(cm)
		}
		m.refreshLists()
		if m.form.onRefresh != nil {
			m.form.onRefresh()
		}
	}
	showFormDialog("Materiale numerabile", items, callback, m.form.win)
}

func (m *materialsAccordion) showMeasurableMaterialDialog(existing *material.MeasurableMaterial) {
	if m.activity == nil {
		return
	}
	nameE := widget.NewEntry()
	descE := widget.NewEntry()
	priceE := widget.NewEntry()
	priceE.SetPlaceHolder("0")
	currE := widget.NewEntry()
	currE.SetText("EUR")
	valE := widget.NewEntry()
	valE.SetPlaceHolder("0")
	unitSelect := widget.NewSelect([]string{"m", "m²", "kg", "g", "day"}, nil)
	unitSelect.SetSelected("m")
	if existing != nil {
		nameE.SetText(existing.Name)
		descE.SetText(existing.Description)
		priceE.SetText(strconv.FormatFloat(existing.Price.Value, 'f', -1, 64))
		currE.SetText(existing.Price.Currency)
		valE.SetText(strconv.FormatFloat(existing.Quantity.Value, 'f', -1, 64))
		unitSelect.SetSelected(string(existing.Quantity.Unit))
	}
	items := []*widget.FormItem{
		widget.NewFormItem("Nome", nameE),
		widget.NewFormItem("Descrizione", descE),
		widget.NewFormItem("Prezzo unitario", container.NewHBox(priceE, currE)),
		widget.NewFormItem("Quantità", container.NewHBox(valE, unitSelect)),
	}
	callback := func(ok bool) {
		if !ok {
			return
		}
		priceVal, _ := strconv.ParseFloat(priceE.Text, 64)
		qtyVal, _ := strconv.ParseFloat(valE.Text, 64)
		if qtyVal < 0 {
			qtyVal = 0
		}
		curr := currE.Text
		if curr == "" {
			curr = "EUR"
		}
		p := unit.Price{Value: priceVal, Currency: curr}
		q := unit.MeasurableQuantity{Value: qtyVal, Unit: unit.MeasurableUnit(unitSelect.Selected)}
		mm := material.NewMeasurableMaterial(nameE.Text, descE.Text, p, q)
		if existing != nil {
			for i, c := range m.activity.MeasurableMaterials {
				if c == existing {
					m.activity.MeasurableMaterials[i] = mm
					break
				}
			}
		} else {
			m.activity.AddMeasurableMaterial(mm)
		}
		m.refreshLists()
		if m.form.onRefresh != nil {
			m.form.onRefresh()
		}
	}
	showFormDialog("Materiale misurabile", items, callback, m.form.win)
}

func (m *materialsAccordion) showHumanResourceDialog(existing *human.HumanResource) {
	if m.activity == nil {
		return
	}
	nameE := widget.NewEntry()
	descE := widget.NewEntry()
	priceE := widget.NewEntry()
	priceE.SetPlaceHolder("0")
	currE := widget.NewEntry()
	currE.SetText("EUR")
	durE := widget.NewEntry()
	durE.SetPlaceHolder("0")
	durSelect := widget.NewSelect(durationUnits, nil)
	durSelect.SetSelected("day")
	if existing != nil {
		nameE.SetText(existing.Name)
		descE.SetText(existing.Description)
		priceE.SetText(strconv.FormatFloat(existing.Price.Value, 'f', -1, 64))
		currE.SetText(existing.Price.Currency)
		durE.SetText(strconv.FormatFloat(existing.Duration.Value, 'f', -1, 64))
		durSelect.SetSelected(string(existing.Duration.Unit))
	}
	items := []*widget.FormItem{
		widget.NewFormItem("Nome", nameE),
		widget.NewFormItem("Descrizione", descE),
		widget.NewFormItem("Prezzo", container.NewHBox(priceE, currE)),
		widget.NewFormItem("Durata", container.NewHBox(durE, durSelect)),
	}
	callback := func(ok bool) {
		if !ok {
			return
		}
		priceVal, _ := strconv.ParseFloat(priceE.Text, 64)
		durVal, _ := strconv.ParseFloat(durE.Text, 64)
		curr := currE.Text
		if curr == "" {
			curr = "EUR"
		}
		p := unit.Price{Value: priceVal, Currency: curr}
		d := unit.Duration{Value: durVal, Unit: unit.DurationUnit(durSelect.Selected)}
		hr := human.NewHumanResource(nameE.Text, descE.Text, d, p)
		if existing != nil {
			for i, c := range m.activity.HumanResources {
				if c == existing {
					m.activity.HumanResources[i] = hr
					break
				}
			}
		} else {
			m.activity.AddHumanResource(hr)
		}
		m.refreshLists()
		if m.form.onRefresh != nil {
			m.form.onRefresh()
		}
	}
	showFormDialog("Risorsa umana", items, callback, m.form.win)
}

func (m *materialsAccordion) showAssetDialog(existing *asset.Asset) {
	if m.activity == nil {
		return
	}
	nameE := widget.NewEntry()
	descE := widget.NewEntry()
	priceE := widget.NewEntry()
	priceE.SetPlaceHolder("0")
	currE := widget.NewEntry()
	currE.SetText("EUR")
	durE := widget.NewEntry()
	durE.SetPlaceHolder("0")
	durSelect := widget.NewSelect(durationUnits, nil)
	durSelect.SetSelected("day")
	if existing != nil {
		nameE.SetText(existing.Name)
		descE.SetText(existing.Description)
		priceE.SetText(strconv.FormatFloat(existing.Price.Value, 'f', -1, 64))
		currE.SetText(existing.Price.Currency)
		durE.SetText(strconv.FormatFloat(existing.Duration.Value, 'f', -1, 64))
		durSelect.SetSelected(string(existing.Duration.Unit))
	}
	items := []*widget.FormItem{
		widget.NewFormItem("Nome", nameE),
		widget.NewFormItem("Descrizione", descE),
		widget.NewFormItem("Prezzo", container.NewHBox(priceE, currE)),
		widget.NewFormItem("Durata", container.NewHBox(durE, durSelect)),
	}
	callback := func(ok bool) {
		if !ok {
			return
		}
		priceVal, _ := strconv.ParseFloat(priceE.Text, 64)
		durVal, _ := strconv.ParseFloat(durE.Text, 64)
		curr := currE.Text
		if curr == "" {
			curr = "EUR"
		}
		p := unit.Price{Value: priceVal, Currency: curr}
		d := unit.Duration{Value: durVal, Unit: unit.DurationUnit(durSelect.Selected)}
		as := asset.NewAsset(nameE.Text, descE.Text, p, d)
		if existing != nil {
			for i, c := range m.activity.Assets {
				if c == existing {
					m.activity.Assets[i] = as
					break
				}
			}
		} else {
			m.activity.AddAsset(as)
		}
		m.refreshLists()
		if m.form.onRefresh != nil {
			m.form.onRefresh()
		}
	}
	showFormDialog("Asset", items, callback, m.form.win)
}

func showFormDialog(title string, items []*widget.FormItem, callback func(bool), win fyne.Window) {
	if win == nil {
		return
	}
	d := dialog.NewForm(title, "OK", "Annulla", items, callback, win)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}
