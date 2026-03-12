// Package viewmodel provides view models for the Explosio GUI.
package viewmodel

import (
	"explosio/core"
	"explosio/core/unit"
	"explosio/gui/app"
	"strconv"
)

// ActivityViewModel holds the editable fields of an activity for form binding.
type ActivityViewModel struct {
	Name          string
	Description   string
	DurationValue float64
	DurationUnit  string
	PriceValue    float64
	Currency      string
}

// FromActivity creates a ViewModel from an Activity.
func FromActivity(a *core.Activity) *ActivityViewModel {
	if a == nil {
		return &ActivityViewModel{
			DurationUnit: "day",
			Currency:     "EUR",
		}
	}
	unitStr := string(a.Duration.Unit)
	if unitStr == "" {
		unitStr = "day"
	}
	curr := a.Price.Currency
	if curr == "" {
		curr = "EUR"
	}
	return &ActivityViewModel{
		Name:          a.Name,
		Description:   a.Description,
		DurationValue: a.Duration.Value,
		DurationUnit:  unitStr,
		PriceValue:    a.Price.Value,
		Currency:      curr,
	}
}

// ToUpdateFields converts the ViewModel to app.UpdateFields for ActivityService.
func (v *ActivityViewModel) ToUpdateFields() app.UpdateFields {
	durUnit := v.DurationUnit
	if durUnit == "" {
		durUnit = "day"
	}
	curr := v.Currency
	if curr == "" {
		curr = "EUR"
	}
	return app.UpdateFields{
		Name:        v.Name,
		Description: v.Description,
		Duration:    unit.Duration{Value: v.DurationValue, Unit: unit.DurationUnit(durUnit)},
		Price:       unit.Price{Value: v.PriceValue, Currency: curr},
	}
}

// ParseFromForm updates the ViewModel from form string values.
// Used when reading from entry widgets.
func (v *ActivityViewModel) ParseFromForm(name, desc, durationStr, durationUnit, priceStr, currency string) {
	v.Name = name
	v.Description = desc
	if f, err := strconv.ParseFloat(durationStr, 64); err == nil {
		v.DurationValue = f
	}
	v.DurationUnit = durationUnit
	if durationUnit == "" {
		v.DurationUnit = "day"
	}
	if f, err := strconv.ParseFloat(priceStr, 64); err == nil {
		v.PriceValue = f
	}
	v.Currency = currency
	if v.Currency == "" {
		v.Currency = "EUR"
	}
}
