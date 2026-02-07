// Package core provides the Activity model and operations to build the tree of activities, materials, and durations.
package core

import (
	"explosio/core/material"
	"explosio/core/unit"
)

// Activity represents a task with name, description, duration, price, sub-activities, and materials (complex, countable, measurable).
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

// NewActivity creates an activity with name and description, zero duration and zero EUR price.
func NewActivity(name string, description string) *Activity {
	return &Activity{
		Name:        name,
		Description: description,
		Duration:    *unit.NewDuration(0, unit.DurationUnitHour),
		Price:       *unit.NewPrice(0, "EUR"),
	}
}

// SetDuration sets the activity duration and returns the activity for chaining.
func (a *Activity) SetDuration(duration unit.Duration) *Activity {
	a.Duration = duration
	return a
}

// SetPrice sets the activity price and returns the activity for chaining.
func (a *Activity) SetPrice(price unit.Price) *Activity {
	a.Price = price
	return a
}

// AddActivity adds a sub-activity.
func (a *Activity) AddActivity(activity *Activity) {
	a.Activities = append(a.Activities, activity)
}

// AddComplexMaterial adds a complex material.
func (a *Activity) AddComplexMaterial(complexMaterial *material.ComplexMaterial) {
	a.ComplexMaterials = append(a.ComplexMaterials, complexMaterial)
}

// AddCountableMaterial adds a countable material.
func (a *Activity) AddCountableMaterial(countableMaterial *material.CountableMaterial) {
	a.CountableMaterials = append(a.CountableMaterials, countableMaterial)
}

// AddMeasurableMaterial adds a measurable material.
func (a *Activity) AddMeasurableMaterial(measurableMaterial *material.MeasurableMaterial) {
	a.MeasurableMaterials = append(a.MeasurableMaterials, measurableMaterial)
}

