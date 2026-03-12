package app

import (
	"explosio/core"
	"explosio/core/material"
	"explosio/core/resource/asset"
	"explosio/core/resource/human"
)

// MaterialsService handles material and resource operations on activities.
type MaterialsService struct{}

// NewMaterialsService creates a new MaterialsService.
func NewMaterialsService() *MaterialsService {
	return &MaterialsService{}
}

// AddComplexMaterial adds a complex material to the activity.
func (s *MaterialsService) AddComplexMaterial(a *core.Activity, m *material.ComplexMaterial) {
	if a == nil || m == nil {
		return
	}
	a.AddComplexMaterial(m)
}

// ReplaceComplexMaterial replaces an existing complex material with a new one.
func (s *MaterialsService) ReplaceComplexMaterial(a *core.Activity, old *material.ComplexMaterial, new *material.ComplexMaterial) {
	if a == nil || new == nil {
		return
	}
	for i, c := range a.ComplexMaterials {
		if c == old {
			a.ComplexMaterials[i] = new
			return
		}
	}
}

// RemoveComplexMaterial removes the complex material at the given index.
func (s *MaterialsService) RemoveComplexMaterial(a *core.Activity, index int) {
	if a == nil || index < 0 || index >= len(a.ComplexMaterials) {
		return
	}
	a.ComplexMaterials = append(a.ComplexMaterials[:index], a.ComplexMaterials[index+1:]...)
}

// AddCountableMaterial adds a countable material to the activity.
func (s *MaterialsService) AddCountableMaterial(a *core.Activity, m *material.CountableMaterial) {
	if a == nil || m == nil {
		return
	}
	a.AddCountableMaterial(m)
}

// ReplaceCountableMaterial replaces an existing countable material with a new one.
func (s *MaterialsService) ReplaceCountableMaterial(a *core.Activity, old *material.CountableMaterial, new *material.CountableMaterial) {
	if a == nil || new == nil {
		return
	}
	for i, c := range a.CountableMaterials {
		if c == old {
			a.CountableMaterials[i] = new
			return
		}
	}
}

// RemoveCountableMaterial removes the countable material at the given index.
func (s *MaterialsService) RemoveCountableMaterial(a *core.Activity, index int) {
	if a == nil || index < 0 || index >= len(a.CountableMaterials) {
		return
	}
	a.CountableMaterials = append(a.CountableMaterials[:index], a.CountableMaterials[index+1:]...)
}

// AddMeasurableMaterial adds a measurable material to the activity.
func (s *MaterialsService) AddMeasurableMaterial(a *core.Activity, m *material.MeasurableMaterial) {
	if a == nil || m == nil {
		return
	}
	a.AddMeasurableMaterial(m)
}

// ReplaceMeasurableMaterial replaces an existing measurable material with a new one.
func (s *MaterialsService) ReplaceMeasurableMaterial(a *core.Activity, old *material.MeasurableMaterial, new *material.MeasurableMaterial) {
	if a == nil || new == nil {
		return
	}
	for i, c := range a.MeasurableMaterials {
		if c == old {
			a.MeasurableMaterials[i] = new
			return
		}
	}
}

// RemoveMeasurableMaterial removes the measurable material at the given index.
func (s *MaterialsService) RemoveMeasurableMaterial(a *core.Activity, index int) {
	if a == nil || index < 0 || index >= len(a.MeasurableMaterials) {
		return
	}
	a.MeasurableMaterials = append(a.MeasurableMaterials[:index], a.MeasurableMaterials[index+1:]...)
}

// AddHumanResource adds a human resource to the activity.
func (s *MaterialsService) AddHumanResource(a *core.Activity, h *human.HumanResource) {
	if a == nil || h == nil {
		return
	}
	a.AddHumanResource(h)
}

// ReplaceHumanResource replaces an existing human resource with a new one.
func (s *MaterialsService) ReplaceHumanResource(a *core.Activity, old *human.HumanResource, new *human.HumanResource) {
	if a == nil || new == nil {
		return
	}
	for i, c := range a.HumanResources {
		if c == old {
			a.HumanResources[i] = new
			return
		}
	}
}

// RemoveHumanResource removes the human resource at the given index.
func (s *MaterialsService) RemoveHumanResource(a *core.Activity, index int) {
	if a == nil || index < 0 || index >= len(a.HumanResources) {
		return
	}
	a.HumanResources = append(a.HumanResources[:index], a.HumanResources[index+1:]...)
}

// AddAsset adds an asset to the activity.
func (s *MaterialsService) AddAsset(a *core.Activity, as *asset.Asset) {
	if a == nil || as == nil {
		return
	}
	a.AddAsset(as)
}

// ReplaceAsset replaces an existing asset with a new one.
func (s *MaterialsService) ReplaceAsset(a *core.Activity, old *asset.Asset, new *asset.Asset) {
	if a == nil || new == nil {
		return
	}
	for i, c := range a.Assets {
		if c == old {
			a.Assets[i] = new
			return
		}
	}
}

// RemoveAsset removes the asset at the given index.
func (s *MaterialsService) RemoveAsset(a *core.Activity, index int) {
	if a == nil || index < 0 || index >= len(a.Assets) {
		return
	}
	a.Assets = append(a.Assets[:index], a.Assets[index+1:]...)
}
