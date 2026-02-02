package lib

type ResourceType string

// PeriodType rappresenta il tipo di periodo per la capacità del fornitore
type PeriodType string

const (
	PeriodMinute PeriodType = "minuto"
	PeriodHour   PeriodType = "ora"
	PeriodDay    PeriodType = "giorno"
	PeriodWeek   PeriodType = "settimana"
	PeriodMonth  PeriodType = "mese"
	PeriodYear   PeriodType = "anno"
)

// ToMinutes converte un PeriodType in minuti
func (p PeriodType) ToMinutes() int {
	switch p {
	case PeriodMinute:
		return 1
	case PeriodHour:
		return 60
	case PeriodDay:
		return 1440
	case PeriodWeek:
		return 10080
	case PeriodMonth:
		return 43200 // ~30 giorni
	case PeriodYear:
		return 525600 // ~365 giorni
	default:
		return 1
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

type Supplier struct {
	Name              string
	Description       string
	UnitCost          float64
	AvailableQuantity float64
	Period            PeriodType
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
	// Converti in capacità giornaliera: (AvailableQuantity / periodo_in_minuti) * minuti_in_un_giorno
	return (s.AvailableQuantity / float64(periodMinutes)) * 1440.0
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
	// Converti: (AvailableQuantity / periodo_sorgente) * periodo_target
	return (s.AvailableQuantity / float64(sourceMinutes)) * float64(targetMinutes)
}

type HumanResource struct {
	Role        string
	Description string
	CostPerH    float64
	Quantity    float64
	Supplier    *Supplier // Fornitore opzionale (es. agenzia lavoro)
}

// GetCost calcola il costo della risorsa umana
// Il fornitore (se presente) serve solo per validare la capacità, non per il calcolo del costo
func (h HumanResource) GetCost(duration int) float64 {
	return (h.CostPerH / 60.0) * float64(duration) * h.Quantity
}

// GetQuantity restituisce la quantità della risorsa
func (h HumanResource) GetQuantity() float64 {
	return h.Quantity
}

// GetSupplier restituisce il fornitore associato
func (h HumanResource) GetSupplier() *Supplier {
	return h.Supplier
}

type MaterialResource struct {
	Name        string
	Description string
	UnitCost    float64
	Quantity    float64
	Supplier    *Supplier // Fornitore opzionale
}

// GetCost calcola il costo del materiale
// Il fornitore (se presente) serve solo per validare la capacità, non per il calcolo del costo
func (m MaterialResource) GetCost(duration int) float64 {
	return m.UnitCost * m.Quantity
}

// GetQuantity restituisce la quantità del materiale
func (m MaterialResource) GetQuantity() float64 {
	return m.Quantity
}

// GetSupplier restituisce il fornitore associato
func (m MaterialResource) GetSupplier() *Supplier {
	return m.Supplier
}

type Asset struct {
	Name        string
	Description string
	CostPerUse  float64
	Quantity    float64
	Supplier    *Supplier // Fornitore opzionale (es. noleggio)
}

// GetCost calcola il costo dell'asset
// Il fornitore (se presente) serve solo per validare la capacità, non per il calcolo del costo
func (a Asset) GetCost(duration int) float64 {
	return a.CostPerUse * a.Quantity
}

// GetQuantity restituisce la quantità dell'asset
func (a Asset) GetQuantity() float64 {
	return a.Quantity
}

// GetSupplier restituisce il fornitore associato
func (a Asset) GetSupplier() *Supplier {
	return a.Supplier
}

// Activity rappresenta il core del nostro dominio
type Activity struct {
	ID          string
	Name        string
	Description string
	Duration    int

	// Parametri per Crashing
	MinDuration   int
	CrashCostStep float64

	// Risorse allocate
	Humans    []HumanResource
	Materials []MaterialResource
	Assets    []Asset

	// Relazioni Grafo
	Next          []string
	SubActivities []*Activity

	// Dati calcolati (CPM)
	ES, EF, LS, LF int
	Slack          int
}
