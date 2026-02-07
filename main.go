package main

import (
	"explosio/core"
	"explosio/core/unit"
	"fmt"
)

func main() {
	activity := core.NewActivity("Install a pipe", "Install a pipe 1 meter long", unit.Duration{Value: 1, Unit: unit.DurationUnitHour}, unit.Price{Value: 100, Currency: "EUR"})
	activity.AddActivity(core.NewActivity("Install a pipe", "Install a pipe 1 meter long", unit.Duration{Value: 1, Unit: unit.DurationUnitHour}, unit.Price{Value: 100, Currency: "EUR"}))
	fmt.Println(activity)
}
