package gui

import (
	"explosio/core"
	"explosio/gui/app"
	"explosio/gui/viewmodel"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)


var durationUnits = []string{"minute", "hour", "day", "week", "month", "year"}

// ActivityForm is a form panel for editing an Activity's base fields.
type ActivityForm struct {
	root        *core.Activity
	current     *core.Activity
	actSvc      *app.ActivityService
	matSvc      *app.MaterialsService
	onRefresh   func()
	win         fyne.Window
	loading     bool // true durante loadFromActivity, evita che SetText attivi saveToActivity

	// Base fields
	nameEntry       *widget.Entry
	descEntry       *widget.Entry
	durationEntry   *widget.Entry
	durationSelect  *widget.Select
	priceEntry      *widget.Entry
	currencyEntry   *widget.Entry

	// Materials/resources accordion (from form_materials.go)
	materialsAccordion *materialsAccordion

	content         *fyne.Container
	centerContent   *fyne.Container // Stack: placeholder o form scroll
	emptyPlaceholder fyne.CanvasObject
	formScroll      fyne.CanvasObject
}

// NewActivityForm creates a new activity form. actSvc and matSvc are used for updating activities and materials; onRefresh is called when the tree should refresh.
// win is used for dialogs (can be set later via SetWindow).
func NewActivityForm(root *core.Activity, actSvc *app.ActivityService, matSvc *app.MaterialsService, onRefresh func(), win fyne.Window) *ActivityForm {
	f := &ActivityForm{
		root:      root,
		actSvc:    actSvc,
		matSvc:    matSvc,
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

	baseInfoForm := widget.NewForm(
		widget.NewFormItem("Nome", f.nameEntry),
		widget.NewFormItem("Descrizione", f.descEntry),
	)
	tempiCostiForm := widget.NewForm(
		widget.NewFormItem("Durata", container.NewHBox(f.durationEntry, f.durationSelect)),
		widget.NewFormItem("Prezzo", container.NewHBox(f.priceEntry, f.currencyEntry)),
	)

	f.materialsAccordion = newMaterialsAccordion(f, matSvc)

	baseLabel := widget.NewLabelWithStyle("Informazioni base", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	tempiLabel := widget.NewLabelWithStyle("Tempi e costi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	scroll := container.NewScroll(container.NewVBox(
		baseLabel,
		baseInfoForm,
		widget.NewSeparator(),
		tempiLabel,
		tempiCostiForm,
		widget.NewSeparator(),
		f.materialsAccordion.accordion,
	))
	scroll.SetMinSize(fyne.NewSize(300, 550))

	f.emptyPlaceholder = widget.NewLabelWithStyle("Seleziona un'attività dal tree", fyne.TextAlignCenter, fyne.TextStyle{})
	f.formScroll = scroll
	f.centerContent = container.NewStack(f.emptyPlaceholder)
	f.content = container.NewBorder(
		widget.NewLabelWithStyle("Dettagli attività", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		f.centerContent,
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
	if f.current == nil || f.actSvc == nil {
		return
	}
	vm := &viewmodel.ActivityViewModel{}
	vm.ParseFromForm(
		f.nameEntry.Text,
		f.descEntry.Text,
		f.durationEntry.Text,
		f.durationSelect.Selected,
		f.priceEntry.Text,
		f.currencyEntry.Text,
	)
	f.actSvc.UpdateActivity(f.current, vm.ToUpdateFields())
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
		f.centerContent.Objects = []fyne.CanvasObject{f.emptyPlaceholder}
		f.centerContent.Refresh()
		return
	}
	f.centerContent.Objects = []fyne.CanvasObject{f.formScroll}
	f.centerContent.Refresh()

	vm := viewmodel.FromActivity(a)
	f.nameEntry.SetText(vm.Name)
	f.descEntry.SetText(vm.Description)
	f.durationEntry.SetText(strconv.FormatFloat(vm.DurationValue, 'f', -1, 64))
	f.durationSelect.SetSelected(vm.DurationUnit)
	f.priceEntry.SetText(strconv.FormatFloat(vm.PriceValue, 'f', -1, 64))
	f.currencyEntry.SetText(vm.Currency)
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
