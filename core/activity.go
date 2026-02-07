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

/*
Quantity calculation:
*/
func (a *Activity) CalculateQuantity() int {
	quantity := 0
	for _, complexMaterial := range a.ComplexMaterials {
		quantity += complexMaterial.UnitQuantity
	}
	for _, countableMaterial := range a.CountableMaterials {
		quantity += countableMaterial.Quantity
	}
	for range a.MeasurableMaterials {
		quantity += 1
	}
	for _, activity := range a.Activities {
		quantity += activity.CalculateQuantity()
	}
	return quantity
}

/*
Critical Path (classic CPM interpretation).

When sub-activities are considered in parallel, the critical path is the longest
path from the root activity to any leaf. It determines the minimum project
duration: any delay on this path delays the whole project.
*/
func (a *Activity) CalculateCriticalPath() []*Activity {
	path, _ := a.criticalPathAndDuration()
	return path
}

// criticalPathAndDuration returns the critical path and its total duration in hours.
// For a leaf (no children), the path is just this activity.
// For a node with children, it picks the child whose subtree has the longest
// total duration and appends that child's critical path to this activity.
func (a *Activity) criticalPathAndDuration() ([]*Activity, float64) {
	myHours := a.Duration.ToHours()
	if len(a.Activities) == 0 {
		return []*Activity{a}, myHours
	}
	var bestPath []*Activity
	bestHours := -1.0
	for _, child := range a.Activities {
		childPath, childHours := child.criticalPathAndDuration()
		total := myHours + childHours
		if total > bestHours {
			bestHours = total
			bestPath = childPath
		}
	}
	return append([]*Activity{a}, bestPath...), bestHours
}
