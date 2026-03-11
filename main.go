// Package main is the entry point of the explosio application.
package main

import (
	"explosio/core"
	"explosio/core/material"
	"explosio/core/resource/asset"
	"explosio/core/resource/human"
	"explosio/core/unit"
	"fmt"
	"log"
)

// must panics if err is non-nil. Used for Build() calls in demo.
func must[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}

// main builds a sample activity tree (home renovation), prints it with PrettyPrint, and shows totals and critical path.
func main() {
	// Level 1 (root)
	homeRenovation := must(core.NewActivityBuilder().
		WithName("Home Renovation").
		WithDescription("Complete home renovation project").
		WithDuration(*unit.NewDuration(40, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(50000, "EUR")).
		Build())

	// Level 2
	kitchen := must(core.NewActivityBuilder().
		WithName("Kitchen Renovation").
		WithDescription("Kitchen remodeling").
		WithDuration(*unit.NewDuration(15, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(20000, "EUR")).
		Build())

	bathroom := must(core.NewActivityBuilder().
		WithName("Bathroom Renovation").
		WithDescription("Bathroom remodeling").
		WithDuration(*unit.NewDuration(10, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(12000, "EUR")).
		Build())

	// Level 3
	installPipes := must(core.NewActivityBuilder().
		WithName("Install pipes").
		WithDescription("Plumbing installation in kitchen").
		WithDuration(*unit.NewDuration(3, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(2500, "EUR")).
		Build())

	installElectrical := must(core.NewActivityBuilder().
		WithName("Install electrical").
		WithDescription("Electrical wiring in kitchen").
		WithDuration(*unit.NewDuration(2, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(1800, "EUR")).
		Build())

	installTiles := must(core.NewActivityBuilder().
		WithName("Install tiles").
		WithDescription("Tile installation in bathroom").
		WithDuration(*unit.NewDuration(4, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(3500, "EUR")).
		Build())

	// Level 4
	cutAndFitPipes := must(core.NewActivityBuilder().
		WithName("Cut and fit pipes").
		WithDescription("Cut and fit pipes to length").
		WithDuration(*unit.NewDuration(1.5, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(300, "EUR")).
		Build())

	weldJoints := must(core.NewActivityBuilder().
		WithName("Weld joints").
		WithDescription("Weld pipe joints").
		WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(250, "EUR")).
		Build())

	runCables := must(core.NewActivityBuilder().
		WithName("Run cables").
		WithDescription("Run electrical cables through walls").
		WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(200, "EUR")).
		Build())

	mountSwitches := must(core.NewActivityBuilder().
		WithName("Mount switches").
		WithDescription("Mount light switches and outlets").
		WithDuration(*unit.NewDuration(0.5, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(150, "EUR")).
		Build())

	prepareSurface := must(core.NewActivityBuilder().
		WithName("Prepare surface").
		WithDescription("Prepare floor surface for tiling").
		WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(180, "EUR")).
		Build())

	applyAdhesiveAndLayTiles := must(core.NewActivityBuilder().
		WithName("Apply adhesive and lay tiles").
		WithDescription("Apply adhesive and lay floor tiles").
		WithDuration(*unit.NewDuration(2, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(400, "EUR")).
		Build())

	// Level 5
	measureAndMark := must(core.NewActivityBuilder().
		WithName("Measure and mark").
		WithDescription("Measure and mark pipe cut points").
		WithDuration(*unit.NewDuration(1, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(100, "EUR")).
		Build())

	applyGrout := must(core.NewActivityBuilder().
		WithName("Apply grout").
		WithDescription("Apply grout between tiles").
		WithDuration(*unit.NewDuration(0.5, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(80, "EUR")).
		Build())

	// Materials: installPipes — pipes (complex), cement (measurable), screws (countable)
	pipeUnit := must(material.NewMeasurableMaterialBuilder().
		WithName("Pipe 2m").
		WithDescription("Copper pipe 2m").
		WithPrice(*unit.NewPrice(50, "EUR")).
		WithQuantity(*unit.NewMeasurableQuantity(2, unit.UnitMeter)).
		Build())
	pipes := must(material.NewComplexMaterialBuilder().
		WithName("Pipes").
		WithDescription("Plumbing pipes").
		WithPrice(*unit.NewPrice(100, "EUR")).
		WithUnitQuantity(5).
		WithMeasurableMaterial(pipeUnit).
		Build())
	cement := must(material.NewMeasurableMaterialBuilder().
		WithName("Cement").
		WithDescription("Bags of cement").
		WithPrice(*unit.NewPrice(15, "EUR")).
		WithQuantity(*unit.NewMeasurableQuantity(50, unit.UnitKilogram)).
		Build())
	screwsPipes := must(material.NewCountableMaterialBuilder().
		WithName("Screws").
		WithDescription("Pipe mounting screws").
		WithPrice(*unit.NewPrice(0.5, "EUR")).
		WithQuantity(80).
		Build())

	// installElectrical: cable (measurable), switches (countable)
	electricalCable := must(material.NewMeasurableMaterialBuilder().
		WithName("Electrical cable").
		WithDescription("Copper electrical cable").
		WithPrice(*unit.NewPrice(2, "EUR")).
		WithQuantity(*unit.NewMeasurableQuantity(100, unit.UnitMeter)).
		Build())
	switches := must(material.NewCountableMaterialBuilder().
		WithName("Light switches").
		WithDescription("Double light switches").
		WithPrice(*unit.NewPrice(25, "EUR")).
		WithQuantity(4).
		Build())

	// installTiles: tiles (measurable), grout (measurable), screws (countable)
	tiles := must(material.NewMeasurableMaterialBuilder().
		WithName("Tiles").
		WithDescription("Ceramic floor tiles").
		WithPrice(*unit.NewPrice(35, "EUR")).
		WithQuantity(*unit.NewMeasurableQuantity(15, unit.UnitSquareMeter)).
		Build())
	grout := must(material.NewMeasurableMaterialBuilder().
		WithName("Grout").
		WithDescription("Tile grout").
		WithPrice(*unit.NewPrice(8, "EUR")).
		WithQuantity(*unit.NewMeasurableQuantity(10, unit.UnitKilogram)).
		Build())
	screwsTiles := must(material.NewCountableMaterialBuilder().
		WithName("Tile anchors").
		WithDescription("Wall anchors for tiles").
		WithPrice(*unit.NewPrice(0.2, "EUR")).
		WithQuantity(50).
		Build())

	// Human resources
	plumber := must(human.NewHumanResourceBuilder().
		WithName("Plumber").
		WithDescription("Licensed plumber for pipe installation").
		WithDuration(*unit.NewDuration(3, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(800, "EUR")).
		Build())
	electrician := must(human.NewHumanResourceBuilder().
		WithName("Electrician").
		WithDescription("Certified electrician for wiring").
		WithDuration(*unit.NewDuration(2, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(600, "EUR")).
		Build())
	tiler := must(human.NewHumanResourceBuilder().
		WithName("Tiler").
		WithDescription("Tile layer for bathroom").
		WithDuration(*unit.NewDuration(4, unit.DurationUnitDay)).
		WithPrice(*unit.NewPrice(900, "EUR")).
		Build())

	// Assets (equipment / tools)
	pipeCutter := must(asset.NewAssetBuilder().
		WithName("Pipe cutter").
		WithDescription("Professional pipe cutting tool").
		WithPrice(*unit.NewPrice(120, "EUR")).
		WithDuration(*unit.NewDuration(0, unit.DurationUnitDay)).
		Build())
	weldingKit := must(asset.NewAssetBuilder().
		WithName("Welding kit").
		WithDescription("Portable welding equipment").
		WithPrice(*unit.NewPrice(350, "EUR")).
		WithDuration(*unit.NewDuration(0, unit.DurationUnitDay)).
		Build())
	drill := must(asset.NewAssetBuilder().
		WithName("Drill").
		WithDescription("Cordless drill for mounting").
		WithPrice(*unit.NewPrice(80, "EUR")).
		WithDuration(*unit.NewDuration(0, unit.DurationUnitDay)).
		Build())
	tileCutter := must(asset.NewAssetBuilder().
		WithName("Tile cutter").
		WithDescription("Manual tile cutter").
		WithPrice(*unit.NewPrice(60, "EUR")).
		WithDuration(*unit.NewDuration(0, unit.DurationUnitDay)).
		Build())

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

	// Add human resources to activities
	installPipes.AddHumanResource(plumber)
	installElectrical.AddHumanResource(electrician)
	installTiles.AddHumanResource(tiler)

	// Add assets to activities
	installPipes.AddAsset(pipeCutter)
	installPipes.AddAsset(weldingKit)
	installElectrical.AddAsset(drill)
	installTiles.AddAsset(tileCutter)

	core.PrettyPrint([]*core.Activity{homeRenovation}, homeRenovation.CalculateCriticalPath())
	fmt.Printf("\nTotal price: %.2f EUR\n", homeRenovation.CalculatePrice())
	fmt.Printf("Total duration: %.0f days\n", homeRenovation.CalculateDuration())
	meas := homeRenovation.GetMeasurableMaterials()
	countable := homeRenovation.GetCountableMaterials()
	complexMat := homeRenovation.GetComplexMaterials()
	hr := homeRenovation.GetHumanResources()
	assets := homeRenovation.GetAssets()
	fmt.Printf("Materials: %d measurable, %d countable, %d complex | Human resources: %d | Assets: %d\n", len(meas), len(countable), len(complexMat), len(hr), len(assets))
}
