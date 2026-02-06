package domain

import "fmt"

// PeriodType rappresenta il tipo di periodo per la capacità del fornitore.
// Esempi di utilizzo:
//
//	supplier := NewSupplier("Fornitore", "Desc", 10.0, 100, PeriodDay)
//	capacity := supplier.GetCapacityForPeriod(PeriodWeek)
type PeriodType string

const (
	PeriodMinute PeriodType = "minuto"
	PeriodHour   PeriodType = "ora"
	PeriodDay    PeriodType = "giorno"
	PeriodWeek   PeriodType = "settimana"
	PeriodMonth  PeriodType = "mese"
	PeriodYear   PeriodType = "anno"
)

// IsValid verifica se il PeriodType è valido
func (p PeriodType) IsValid() bool {
	switch p {
	case PeriodMinute, PeriodHour, PeriodDay, PeriodWeek, PeriodMonth, PeriodYear:
		return true
	default:
		return false
	}
}

// ToMinutes converte un PeriodType in minuti.
// Restituisce 0 se il periodo non è valido.
func (p PeriodType) ToMinutes() int {
	switch p {
	case PeriodMinute:
		return 1
	case PeriodHour:
		return MinutesPerHour
	case PeriodDay:
		return MinutesPerDay
	case PeriodWeek:
		return MinutesPerWeek
	case PeriodMonth:
		return MinutesPerMonth
	case PeriodYear:
		return MinutesPerYear
	default:
		return 0 // periodo invalido
	}
}

// String restituisce la rappresentazione stringa del periodo
func (p PeriodType) String() string {
	return string(p)
}

// Resource è l'interfaccia comune per tutte le risorse
type Resource interface {
	GetCost(duration int) float64
	GetQuantity() float64
	GetSupplier() *Supplier
}

// Supplier rappresenta un fornitore esterno che può fornire risorse con una capacità limitata per periodo.
type Supplier struct {
	Name              string
	Description       string
	UnitCost          float64
	AvailableQuantity float64
	Period            PeriodType
}

// Validate valida i campi del Supplier
func (s *Supplier) Validate() error {
	if s == nil {
		return nil
	}
	if s.AvailableQuantity < 0 {
		return fmt.Errorf("%w: supplier %s has negative available quantity", ErrNegativeQuantity, s.Name)
	}
	if !s.Period.IsValid() {
		return fmt.Errorf("%w: supplier %s has invalid period %s", ErrInvalidPeriod, s.Name, s.Period)
	}
	return nil
}

// GetDailyCapacity converte la capacità disponibile in capacità giornaliera
func (s *Supplier) GetDailyCapacity() float64 {
	if s == nil {
		return 0
	}
	periodMinutes := s.Period.ToMinutes()
	if periodMinutes == 0 {
		return 0
	}
	return (s.AvailableQuantity / float64(periodMinutes)) * float64(MinutesPerDay)
}

// GetCapacityForPeriod converte la capacità disponibile in un periodo specifico
func (s *Supplier) GetCapacityForPeriod(targetPeriod PeriodType) float64 {
	if s == nil {
		return 0
	}
	sourceMinutes := s.Period.ToMinutes()
	targetMinutes := targetPeriod.ToMinutes()
	if sourceMinutes == 0 || targetMinutes == 0 {
		return 0
	}
	return (s.AvailableQuantity / float64(sourceMinutes)) * float64(targetMinutes)
}

type HumanResource struct {
	Role        string
	Description string
	CostPerH    float64
	Quantity    float64
	Supplier    *Supplier
}

func (h HumanResource) GetCost(duration int) float64 {
	return (h.CostPerH / float64(MinutesPerHour)) * float64(duration) * h.Quantity
}

func (h HumanResource) GetQuantity() float64 {
	return h.Quantity
}

func (h HumanResource) GetSupplier() *Supplier {
	return h.Supplier
}

func (h HumanResource) Validate() error {
	if h.CostPerH < 0 {
		return fmt.Errorf("%w: human resource %s has negative cost per hour", ErrNegativeCost, h.Role)
	}
	if h.Quantity < 0 {
		return fmt.Errorf("%w: human resource %s has negative quantity", ErrNegativeQuantity, h.Role)
	}
	if h.Supplier != nil {
		if err := h.Supplier.Validate(); err != nil {
			return fmt.Errorf("human resource %s: %w", h.Role, err)
		}
	}
	return nil
}

type MaterialResource struct {
	Name        string
	Description string
	UnitCost    float64
	Quantity    float64
	Supplier    *Supplier
}

func (m MaterialResource) GetCost(duration int) float64 {
	return m.UnitCost * m.Quantity
}

func (m MaterialResource) GetQuantity() float64 {
	return m.Quantity
}

func (m MaterialResource) GetSupplier() *Supplier {
	return m.Supplier
}

func (m MaterialResource) Validate() error {
	if m.UnitCost < 0 {
		return fmt.Errorf("%w: material resource %s has negative unit cost", ErrNegativeCost, m.Name)
	}
	if m.Quantity < 0 {
		return fmt.Errorf("%w: material resource %s has negative quantity", ErrNegativeQuantity, m.Name)
	}
	if m.Supplier != nil {
		if err := m.Supplier.Validate(); err != nil {
			return fmt.Errorf("material resource %s: %w", m.Name, err)
		}
	}
	return nil
}

type Asset struct {
	Name        string
	Description string
	CostPerUse  float64
	Quantity    float64
	Supplier    *Supplier
}

func (a Asset) GetCost(duration int) float64 {
	return a.CostPerUse * a.Quantity
}

func (a Asset) GetQuantity() float64 {
	return a.Quantity
}

func (a Asset) GetSupplier() *Supplier {
	return a.Supplier
}

func (a Asset) Validate() error {
	if a.CostPerUse < 0 {
		return fmt.Errorf("%w: asset %s has negative cost per use", ErrNegativeCost, a.Name)
	}
	if a.Quantity < 0 {
		return fmt.Errorf("%w: asset %s has negative quantity", ErrNegativeQuantity, a.Name)
	}
	if a.Supplier != nil {
		if err := a.Supplier.Validate(); err != nil {
			return fmt.Errorf("asset %s: %w", a.Name, err)
		}
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
