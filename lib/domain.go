package lib

type ResourceType string

type HumanResource struct {
	Role        string
	Description string
	CostPerH    float64
	Quantity    float64
}

type MaterialResource struct {
	Name        string
	Description string
	UnitCost    float64
	Quantity    float64
}

type Asset struct {
	Name        string
	Description string
	CostPerUse  float64
	Quantity    float64
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
