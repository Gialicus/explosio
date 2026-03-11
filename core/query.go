// Package core provides the Activity model and operations to build the tree of activities, materials, and durations.
package core

import (
	"regexp"
	"sort"
	"strings"

	"explosio/core/material"
	"explosio/core/resource/asset"
	"explosio/core/resource/human"
)

// FilterOptions holds criteria for filtering activities.
type FilterOptions struct {
	PriceMin    float64 // Minimum price (0 = no filter)
	PriceMax    float64 // Maximum price (0 = no filter)
	DurationMin float64 // Minimum duration in hours (0 = no filter)
	DurationMax float64 // Maximum duration in hours (0 = no filter)
	Name        string  // Substring match on activity name (empty = no filter)
	NameRegex   string  // Regex match on activity name (empty = no filter)
	MaterialName string // Filter activities that use a material with this name (substring)
	ResourceName string // Filter activities that use a human resource with this name (substring)
}

// FilterActivities returns activities that match the given filter options.
func FilterActivities(activities []*Activity, opts FilterOptions) []*Activity {
	var result []*Activity
	var nameRe *regexp.Regexp
	if opts.NameRegex != "" {
		nameRe, _ = regexp.Compile(opts.NameRegex)
	}
	for _, a := range activities {
		if opts.PriceMin > 0 && a.CalculatePrice() < opts.PriceMin {
			continue
		}
		if opts.PriceMax > 0 && a.CalculatePrice() > opts.PriceMax {
			continue
		}
		durHours := a.calculateDurationHours()
		if opts.DurationMin > 0 && durHours < opts.DurationMin {
			continue
		}
		if opts.DurationMax > 0 && durHours > opts.DurationMax {
			continue
		}
		if opts.Name != "" && !strings.Contains(strings.ToLower(a.Name), strings.ToLower(opts.Name)) {
			continue
		}
		if nameRe != nil && !nameRe.MatchString(a.Name) {
			continue
		}
		if opts.MaterialName != "" {
			hasMaterial := false
			for _, m := range a.GetMeasurableMaterials() {
				if strings.Contains(strings.ToLower(m.Name), strings.ToLower(opts.MaterialName)) {
					hasMaterial = true
					break
				}
			}
			if !hasMaterial {
				for _, m := range a.GetComplexMaterials() {
					if strings.Contains(strings.ToLower(m.Name), strings.ToLower(opts.MaterialName)) {
						hasMaterial = true
						break
					}
				}
			}
			if !hasMaterial {
				for _, m := range a.GetCountableMaterials() {
					if strings.Contains(strings.ToLower(m.Name), strings.ToLower(opts.MaterialName)) {
						hasMaterial = true
						break
					}
				}
			}
			if !hasMaterial {
				continue
			}
		}
		if opts.ResourceName != "" {
			hasResource := false
			for _, h := range a.GetHumanResources() {
				if strings.Contains(strings.ToLower(h.Name), strings.ToLower(opts.ResourceName)) {
					hasResource = true
					break
				}
			}
			if !hasResource {
				continue
			}
		}
		result = append(result, a)
	}
	return result
}

// SortOrder specifies how to sort activities.
type SortOrder int

const (
	SortByName     SortOrder = iota
	SortByPrice    SortOrder = iota
	SortByDuration SortOrder = iota
)

// SortActivities sorts activities by the given order.
func SortActivities(activities []*Activity, order SortOrder) {
	sort.Slice(activities, func(i, j int) bool {
		switch order {
		case SortByPrice:
			return activities[i].CalculatePrice() < activities[j].CalculatePrice()
		case SortByDuration:
			return activities[i].CalculateDuration() < activities[j].CalculateDuration()
		default:
			return strings.Compare(strings.ToLower(activities[i].Name), strings.ToLower(activities[j].Name)) < 0
		}
	})
}

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
