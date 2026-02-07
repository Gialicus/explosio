package core

import (
	"explosio/core/material"
	"explosio/core/unit"
)

type Activity struct {
	Name                string
	Description         string
	Duration            unit.Duration
	Price               unit.Price
	Activities          []*Activity
	ComplexMaterials    []*material.ComplexMaterial
	CountableMaterials  []*material.CountableMaterial
	MeasurableMaterials []*material.MeasurableMaterial
}

func NewActivity(name string, description string) *Activity {
	return &Activity{
		Name:        name,
		Description: description,
		Duration:    *unit.NewDuration(0, unit.DurationUnitHour),
		Price:       *unit.NewPrice(0, "EUR"),
	}
}

func (a *Activity) SetDuration(duration unit.Duration) *Activity {
	a.Duration = duration
	return a
}

func (a *Activity) SetPrice(price unit.Price) *Activity {
	a.Price = price
	return a
}

func (a *Activity) AddActivity(activity *Activity) {
	a.Activities = append(a.Activities, activity)
}

func (a *Activity) AddComplexMaterial(complexMaterial *material.ComplexMaterial) {
	a.ComplexMaterials = append(a.ComplexMaterials, complexMaterial)
}

func (a *Activity) AddCountableMaterial(countableMaterial *material.CountableMaterial) {
	a.CountableMaterials = append(a.CountableMaterials, countableMaterial)
}

func (a *Activity) AddMeasurableMaterial(measurableMaterial *material.MeasurableMaterial) {
	a.MeasurableMaterials = append(a.MeasurableMaterials, measurableMaterial)
}

/*
Price calculation:
*/
func (a *Activity) CalculatePrice() float64 {
	price := a.Price.Value
	for _, complexMaterial := range a.ComplexMaterials {
		price += complexMaterial.CalculatePrice()
	}
	for _, countableMaterial := range a.CountableMaterials {
		price += countableMaterial.CalculatePrice()
	}
	for _, measurableMaterial := range a.MeasurableMaterials {
		price += measurableMaterial.CalculatePrice()
	}
	for _, activity := range a.Activities {
		price += activity.CalculatePrice()
	}
	return price
}

/*
Duration calculation:
*/
func (a *Activity) CalculateDuration() float64 {
	duration := a.Duration.Value
	for _, activity := range a.Activities {
		duration += activity.CalculateDuration()
	}
	return duration
}
