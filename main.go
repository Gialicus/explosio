package main

import (
	"explosio/core"
	"explosio/core/material"
	"explosio/core/unit"
)

func main() {
	activity := core.NewActivity("Install pipes in the basement", "Install pipes in the basement", unit.Duration{Value: 1, Unit: unit.DurationUnitHour})
	sub1 := core.NewActivity("Install a pipe", "Install a pipe 1 meter long in the kitchen", unit.Duration{Value: 1, Unit: unit.DurationUnitHour})
	sub1.AddMeasurableMaterial(material.NewMeasurableMaterial("Pipe", "Pipe 1 meter long", *unit.DefaultPrice(), *unit.DefaultMeasurableQuantity()))
	sub2 := core.NewActivity("Fix with Screws", "Fix with screws the pipe", unit.Duration{Value: 1, Unit: unit.DurationUnitHour})
	sub1.AddActivity(sub2)
	activity.AddActivity(sub1)
	activity.AddActivity(sub1)
	core.PrettyPrint([]*core.Activity{activity})
}
