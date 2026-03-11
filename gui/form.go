package gui

import (
	"explosio/core"
	"explosio/core/unit"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)


var durationUnits = []string{"minute", "hour", "day", "week", "month", "year"}

// ActivityForm is a form panel for editing an Activity's base fields.
type ActivityForm struct {
	root       *core.Activity
	current    *core.Activity
	onRefresh  func()
	win        fyne.Window
	loading    bool // true durante loadFromActivity, evita che SetText attivi saveToActivity

	// Base fields
	nameEntry       *widget.Entry
	descEntry       *widget.Entry
	durationEntry   *widget.Entry
	durationSelect  *widget.Select
	priceEntry      *widget.Entry
	currencyEntry   *widget.Entry

	// Materials/resources accordion (from form_materials.go)
	materialsAccordion *materialsAccordion

	content *fyne.Container
}

// NewActivityForm creates a new activity form. onRefresh is called when the tree should refresh (e.g. after name change).
// win is used for dialogs (can be set later via SetWindow).
func NewActivityForm(root *core.Activity, onRefresh func(), win fyne.Window) *ActivityForm {
	f := &ActivityForm{
		root:      root,
		onRefresh: onRefresh,
		win:       win,
	}

	f.nameEntry = widget.NewEntry()
	f.nameEntry.SetPlaceHolder("Nome attività")
	f.nameEntry.OnChanged = f.onFieldChanged

	f.descEntry = widget.NewMultiLineEntry()
	f.descEntry.SetPlaceHolder("Descrizione")
	f.descEntry.OnChanged = f.onFieldChanged

	f.durationEntry = widget.NewEntry()
	f.durationEntry.SetPlaceHolder("0")
	f.durationEntry.OnChanged = f.onFieldChanged

	f.durationSelect = widget.NewSelect(durationUnits, func(string) { f.onFieldChanged("") })
	f.durationSelect.SetSelected("day")

	f.priceEntry = widget.NewEntry()
	f.priceEntry.SetPlaceHolder("0")
	f.priceEntry.OnChanged = f.onFieldChanged

	f.currencyEntry = widget.NewEntry()
	f.currencyEntry.SetPlaceHolder("EUR")
	f.currencyEntry.OnChanged = f.onFieldChanged

	baseForm := widget.NewForm(
		widget.NewFormItem("Nome", f.nameEntry),
		widget.NewFormItem("Descrizione", f.descEntry),
		widget.NewFormItem("Durata", container.NewHBox(f.durationEntry, f.durationSelect)),
		widget.NewFormItem("Prezzo", f.priceEntry),
		widget.NewFormItem("Valuta", f.currencyEntry),
	)

	f.materialsAccordion = newMaterialsAccordion(f)

	scroll := container.NewScroll(container.NewVBox(
		baseForm,
		widget.NewSeparator(),
		f.materialsAccordion.accordion,
	))
	scroll.SetMinSize(fyne.NewSize(300, 400))

	f.content = container.NewBorder(
		widget.NewLabelWithStyle("Dettagli attività", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		scroll,
	)

	return f
}

func (f *ActivityForm) onFieldChanged(string) {
	if f.loading {
		return
	}
	f.saveToActivity()
	if f.onRefresh != nil {
		f.onRefresh()
	}
}

func (f *ActivityForm) saveToActivity() {
	if f.current == nil {
		return
	}
	f.current.Name = f.nameEntry.Text
	f.current.Description = f.descEntry.Text

	if v, err := strconv.ParseFloat(f.durationEntry.Text, 64); err == nil {
		f.current.Duration.Value = v
	}
	f.current.Duration.Unit = unit.DurationUnit(f.durationSelect.Selected)

	if v, err := strconv.ParseFloat(f.priceEntry.Text, 64); err == nil {
		f.current.Price.Value = v
	}
	f.current.Price.Currency = f.currencyEntry.Text
	if f.current.Price.Currency == "" {
		f.current.Price.Currency = "EUR"
	}
}

func (f *ActivityForm) loadFromActivity(a *core.Activity) {
	f.loading = true
	defer func() { f.loading = false }()
	f.current = a
	if a == nil {
		f.nameEntry.SetText("")
		f.descEntry.SetText("")
		f.durationEntry.SetText("0")
		f.durationSelect.SetSelected("day")
		f.priceEntry.SetText("0")
		f.currencyEntry.SetText("EUR")
		f.materialsAccordion.setActivity(nil)
		return
	}

	f.nameEntry.SetText(a.Name)
	f.descEntry.SetText(a.Description)
	f.durationEntry.SetText(strconv.FormatFloat(a.Duration.Value, 'f', -1, 64))
	unitStr := string(a.Duration.Unit)
	if unitStr == "" {
		unitStr = "day"
	}
	f.durationSelect.SetSelected(unitStr)
	f.priceEntry.SetText(strconv.FormatFloat(a.Price.Value, 'f', -1, 64))
	f.currencyEntry.SetText(a.Price.Currency)
	if f.currencyEntry.Text == "" {
		f.currencyEntry.SetText("EUR")
	}

	f.materialsAccordion.setActivity(a)
}

// SelectActivity loads the given activity into the form. Call when tree selection changes.
func (f *ActivityForm) SelectActivity(a *core.Activity) {
	f.loadFromActivity(a)
}

// Content returns the form's container for embedding in the layout.
func (f *ActivityForm) Content() fyne.CanvasObject {
	return f.content
}

// Current returns the currently selected activity.
func (f *ActivityForm) Current() *core.Activity {
	return f.current
}

// SetWindow sets the window for dialogs (e.g. when adding materials).
func (f *ActivityForm) SetWindow(win fyne.Window) {
	f.win = win
}
