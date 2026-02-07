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

func NewActivity(name string, description string, duration unit.Duration, price unit.Price) *Activity {
	return &Activity{
		Name:        name,
		Description: description,
		Duration:    duration,
		Price:       price,
	}
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
