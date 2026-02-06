package domain

import "fmt"

// Resource è l'interfaccia comune per tutte le risorse
type Resource interface {
	GetCost(duration int) float64
	GetQuantity() float64
}

type HumanResource struct {
	Role        string
	Description string
	CostPerH    float64
	Quantity    float64
}

func (h HumanResource) GetCost(duration int) float64 {
	return (h.CostPerH / float64(MinutesPerHour)) * float64(duration) * h.Quantity
}

func (h HumanResource) GetQuantity() float64 {
	return h.Quantity
}

func (h HumanResource) Validate() error {
	if h.CostPerH < 0 {
		return fmt.Errorf("%w: human resource %s has negative cost per hour", ErrNegativeCost, h.Role)
	}
	if h.Quantity < 0 {
		return fmt.Errorf("%w: human resource %s has negative quantity", ErrNegativeQuantity, h.Role)
	}
	return nil
}

type MaterialResource struct {
	Name        string
	Description string
	UnitCost    float64
	Quantity    float64
}

func (m MaterialResource) GetCost(duration int) float64 {
	return m.UnitCost * m.Quantity
}

func (m MaterialResource) GetQuantity() float64 {
	return m.Quantity
}

func (m MaterialResource) Validate() error {
	if m.UnitCost < 0 {
		return fmt.Errorf("%w: material resource %s has negative unit cost", ErrNegativeCost, m.Name)
	}
	if m.Quantity < 0 {
		return fmt.Errorf("%w: material resource %s has negative quantity", ErrNegativeQuantity, m.Name)
	}
	return nil
}

type Asset struct {
	Name        string
	Description string
	CostPerUse  float64
	Quantity    float64
}

func (a Asset) GetCost(duration int) float64 {
	return a.CostPerUse * a.Quantity
}

func (a Asset) GetQuantity() float64 {
	return a.Quantity
}

func (a Asset) Validate() error {
	if a.CostPerUse < 0 {
		return fmt.Errorf("%w: asset %s has negative cost per use", ErrNegativeCost, a.Name)
	}
	if a.Quantity < 0 {
		return fmt.Errorf("%w: asset %s has negative quantity", ErrNegativeQuantity, a.Name)
	}
	return nil
}

// Activity rappresenta un'attività nel progetto con risorse allocate e relazioni di dipendenza.
type Activity struct {
	ID          string
	Name        string
	Description string
	Duration    int

	MinDuration   int
	CrashCostStep float64

	Humans    []HumanResource
	Materials []MaterialResource
	Assets    []Asset

	Next          []string
	SubActivities []*Activity

	ES, EF, LS, LF int
	Slack          int
}

// ValidateBasic valida i campi base dell'Activity (durata, minDuration)
func (a *Activity) ValidateBasic() error {
	if a == nil {
		return fmt.Errorf("%w: activity is nil", ErrInvalidActivity)
	}
	if a.Duration < 0 {
		return fmt.Errorf("%w: activity %s has negative duration", ErrInvalidActivity, a.ID)
	}
	if a.MinDuration < 0 {
		return fmt.Errorf("%w: activity %s has negative min duration", ErrInvalidActivity, a.ID)
	}
	if a.MinDuration > a.Duration {
		return fmt.Errorf("%w: activity %s has min duration (%d) greater than duration (%d)", ErrInvalidActivity, a.ID, a.MinDuration, a.Duration)
	}
	return nil
}
