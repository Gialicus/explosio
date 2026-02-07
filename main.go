package main

import (
	"explosio/core"
	"explosio/core/unit"
	"fmt"
)

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

	// Build tree
	kitchen.AddActivity(installPipes)
	kitchen.AddActivity(installElectrical)
	bathroom.AddActivity(installTiles)
	homeRenovation.AddActivity(kitchen)
	homeRenovation.AddActivity(bathroom)

	core.PrettyPrint([]*core.Activity{homeRenovation})
	fmt.Printf("\nTotal price: %.2f EUR\n", homeRenovation.CalculatePrice())
	fmt.Printf("Total duration: %.0f days\n", homeRenovation.CalculateDuration())
}
