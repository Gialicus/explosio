package core

import (
	"explosio/core/asset"
	"explosio/core/human"
	"explosio/core/material"
)

func (a *Activity) GetActivities() []*Activity {
	activities := []*Activity{a}
	for _, activity := range a.Activities {
		activities = append(activities, activity.GetActivities()...)
	}
	return activities
}

func (a *Activity) GetComplexMaterials() []*material.ComplexMaterial {
	complexMaterials := a.ComplexMaterials
	for _, activity := range a.Activities {
		complexMaterials = append(complexMaterials, activity.GetComplexMaterials()...)
	}
	return complexMaterials
}

func (a *Activity) GetCountableMaterials() []*material.CountableMaterial {
	countableMaterials := a.CountableMaterials
	for _, activity := range a.Activities {
		countableMaterials = append(countableMaterials, activity.GetCountableMaterials()...)
	}
	return countableMaterials
}

func (a *Activity) GetMeasurableMaterials() []*material.MeasurableMaterial {
	measurableMaterials := a.MeasurableMaterials
	for _, activity := range a.Activities {
		measurableMaterials = append(measurableMaterials, activity.GetMeasurableMaterials()...)
	}
	return measurableMaterials
}

func (a *Activity) GetHumanResources() []*human.HumanResource {
	humanResources := a.HumanResources
	for _, activity := range a.Activities {
		humanResources = append(humanResources, activity.GetHumanResources()...)
	}
	return humanResources
}

func (a *Activity) GetAssets() []*asset.Asset {
	assets := a.Assets
	for _, activity := range a.Activities {
		assets = append(assets, activity.GetAssets()...)
	}
	return assets
}
