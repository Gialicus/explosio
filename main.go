// Package main is the entry point of the explosio application.
package main

import (
	"explosio/core"
	"explosio/core/material"
	"explosio/core/unit"
	"fmt"
)

// main builds a sample activity tree (home renovation), prints it with PrettyPrint, and shows totals and critical path.
func main() {
	// Level 1 (root)
	homeRenovation := core.NewActivity("Home Renovation", "Complete home renovation project")
	homeRenovation.SetDuration(unit.Duration{Value: 40, Unit: unit.DurationUnitDay})
	homeRenovation.SetPrice(unit.Price{Value: 50000, Currency: "EUR"})

	// Level 2
	kitchen := core.NewActivity("Kitchen Renovation", "Kitchen remodeling")
	kitchen.SetDuration(unit.Duration{Value: 15, Unit: unit.DurationUnitDay})
	kitchen.SetPrice(unit.Price{Value: 20000, Currency: "EUR"})

	bathroom := core.NewActivity("Bathroom Renovation", "Bathroom remodeling")
	bathroom.SetDuration(unit.Duration{Value: 10, Unit: unit.DurationUnitDay})
	bathroom.SetPrice(unit.Price{Value: 12000, Currency: "EUR"})

	// Level 3
	installPipes := core.NewActivity("Install pipes", "Plumbing installation in kitchen")
	installPipes.SetDuration(unit.Duration{Value: 3, Unit: unit.DurationUnitDay})
	installPipes.SetPrice(unit.Price{Value: 2500, Currency: "EUR"})

	installElectrical := core.NewActivity("Install electrical", "Electrical wiring in kitchen")
	installElectrical.SetDuration(unit.Duration{Value: 2, Unit: unit.DurationUnitDay})
	installElectrical.SetPrice(unit.Price{Value: 1800, Currency: "EUR"})

	installTiles := core.NewActivity("Install tiles", "Tile installation in bathroom")
	installTiles.SetDuration(unit.Duration{Value: 4, Unit: unit.DurationUnitDay})
	installTiles.SetPrice(unit.Price{Value: 3500, Currency: "EUR"})

	// Level 4
	cutAndFitPipes := core.NewActivity("Cut and fit pipes", "Cut and fit pipes to length")
	cutAndFitPipes.SetDuration(unit.Duration{Value: 1.5, Unit: unit.DurationUnitDay})
	cutAndFitPipes.SetPrice(unit.Price{Value: 300, Currency: "EUR"})

	weldJoints := core.NewActivity("Weld joints", "Weld pipe joints")
	weldJoints.SetDuration(unit.Duration{Value: 1, Unit: unit.DurationUnitDay})
	weldJoints.SetPrice(unit.Price{Value: 250, Currency: "EUR"})

	runCables := core.NewActivity("Run cables", "Run electrical cables through walls")
	runCables.SetDuration(unit.Duration{Value: 1, Unit: unit.DurationUnitDay})
	runCables.SetPrice(unit.Price{Value: 200, Currency: "EUR"})

	mountSwitches := core.NewActivity("Mount switches", "Mount light switches and outlets")
	mountSwitches.SetDuration(unit.Duration{Value: 0.5, Unit: unit.DurationUnitDay})
	mountSwitches.SetPrice(unit.Price{Value: 150, Currency: "EUR"})

	prepareSurface := core.NewActivity("Prepare surface", "Prepare floor surface for tiling")
	prepareSurface.SetDuration(unit.Duration{Value: 1, Unit: unit.DurationUnitDay})
	prepareSurface.SetPrice(unit.Price{Value: 180, Currency: "EUR"})

	applyAdhesiveAndLayTiles := core.NewActivity("Apply adhesive and lay tiles", "Apply adhesive and lay floor tiles")
	applyAdhesiveAndLayTiles.SetDuration(unit.Duration{Value: 2, Unit: unit.DurationUnitDay})
	applyAdhesiveAndLayTiles.SetPrice(unit.Price{Value: 400, Currency: "EUR"})

	// Level 5
	measureAndMark := core.NewActivity("Measure and mark", "Measure and mark pipe cut points")
	measureAndMark.SetDuration(unit.Duration{Value: 1, Unit: unit.DurationUnitDay})
	measureAndMark.SetPrice(unit.Price{Value: 100, Currency: "EUR"})

	applyGrout := core.NewActivity("Apply grout", "Apply grout between tiles")
	applyGrout.SetDuration(unit.Duration{Value: 0.5, Unit: unit.DurationUnitDay})
	applyGrout.SetPrice(unit.Price{Value: 80, Currency: "EUR"})

	// Materials
	// installPipes: pipes (complex), cement (measurable), screws (countable)
	pipeUnit := material.NewMeasurableMaterial("Pipe 2m", "Copper pipe 2m", unit.Price{Value: 50, Currency: "EUR"}, unit.MeasurableQuantity{Value: 2, Unit: unit.UnitMeter})
	pipes := material.NewComplexMaterial("Pipes", "Plumbing pipes", unit.Price{Value: 100, Currency: "EUR"}, 5, pipeUnit)
	cement := material.NewMeasurableMaterial("Cement", "Bags of cement", unit.Price{Value: 15, Currency: "EUR"}, unit.MeasurableQuantity{Value: 50, Unit: unit.UnitKilogram})
	screwsPipes := material.NewCountableMaterial("Screws", "Pipe mounting screws", unit.Price{Value: 0.5, Currency: "EUR"}, 80)

	// installElectrical: cable (measurable), switches (countable)
	electricalCable := material.NewMeasurableMaterial("Electrical cable", "Copper electrical cable", unit.Price{Value: 2, Currency: "EUR"}, unit.MeasurableQuantity{Value: 100, Unit: unit.UnitMeter})
	switches := material.NewCountableMaterial("Light switches", "Double light switches", unit.Price{Value: 25, Currency: "EUR"}, 4)

	// installTiles: tiles (measurable), grout (measurable), screws (countable)
	tiles := material.NewMeasurableMaterial("Tiles", "Ceramic floor tiles", unit.Price{Value: 35, Currency: "EUR"}, unit.MeasurableQuantity{Value: 15, Unit: unit.UnitSquareMeter})
	grout := material.NewMeasurableMaterial("Grout", "Tile grout", unit.Price{Value: 8, Currency: "EUR"}, unit.MeasurableQuantity{Value: 10, Unit: unit.UnitKilogram})
	screwsTiles := material.NewCountableMaterial("Tile anchors", "Wall anchors for tiles", unit.Price{Value: 0.2, Currency: "EUR"}, 50)

	// Build tree
	installPipes.AddActivity(cutAndFitPipes)
	installPipes.AddActivity(weldJoints)
	cutAndFitPipes.AddActivity(measureAndMark)

	installElectrical.AddActivity(runCables)
	installElectrical.AddActivity(mountSwitches)

	installTiles.AddActivity(prepareSurface)
	installTiles.AddActivity(applyAdhesiveAndLayTiles)
	applyAdhesiveAndLayTiles.AddActivity(applyGrout)

	kitchen.AddActivity(installPipes)
	kitchen.AddActivity(installElectrical)
	bathroom.AddActivity(installTiles)
	homeRenovation.AddActivity(kitchen)
	homeRenovation.AddActivity(bathroom)

	// Add materials to activities
	installPipes.AddComplexMaterial(pipes)
	installPipes.AddMeasurableMaterial(cement)
	installPipes.AddCountableMaterial(screwsPipes)
	installElectrical.AddMeasurableMaterial(electricalCable)
	installElectrical.AddCountableMaterial(switches)
	installTiles.AddMeasurableMaterial(tiles)
	installTiles.AddMeasurableMaterial(grout)
	installTiles.AddCountableMaterial(screwsTiles)

	core.PrettyPrint([]*core.Activity{homeRenovation}, homeRenovation.CalculateCriticalPath())
	fmt.Printf("\nTotal price: %.2f EUR\n", homeRenovation.CalculatePrice())
	fmt.Printf("Total duration: %.0f days\n", homeRenovation.CalculateDuration())
	fmt.Printf("Total quantity: %d\n", homeRenovation.CalculateQuantity())
}
