package human

import "explosio/core/unit"

type HumanResource struct {
	Name        string
	Description string
	Duration    unit.Duration
	Price       unit.Price
}

func NewHumanResource(name string, description string, duration unit.Duration, price unit.Price) *HumanResource {
	return &HumanResource{Name: name, Description: description, Duration: duration, Price: price}
}

func (h *HumanResource) CalculatePrice() float64 {
	return h.Price.Value
}

func (h *HumanResource) CalculateDuration() float64 {
	return h.Duration.Value
}

func (h *HumanResource) CalculateHourlyRate() float64 {
	return h.Price.Value / h.Duration.ToHours()
}

func (h *HumanResource) CalculateDailyRate() float64 {
	return h.Price.Value / h.Duration.ToHours() * unit.WORKING_HOURS_PER_DAY
}
