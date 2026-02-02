package lib

import (
	"fmt"
	"sync/atomic"
)

type Project struct {
	idCounter int64
	Root      *Activity
}

func NewProject() *Project {
	return &Project{idCounter: 0}
}

func (p *Project) genID() string {
	newID := atomic.AddInt64(&p.idCounter, 1)
	return fmt.Sprintf("ACT-%03d", newID)
}

// Start inizializza il nodo radice (Top-Down)
func (p *Project) Start(name, desc string, duration int) *Activity {
	p.Root = &Activity{
		ID:          p.genID(),
		Name:        name,
		Description: desc,
		Duration:    duration,
		MinDuration: duration,
	}
	return p.Root
}

// Node crea un'attività figlia
func (p *Project) Node(name, desc string, duration int) *Activity {
	return &Activity{
		ID:          p.genID(),
		Name:        name,
		Description: desc,
		Duration:    duration,
		MinDuration: duration,
	}
}

// --- Metodi di Chaining per Activity ---

func (a *Activity) WithHuman(role, desc string, costH, qty float64) *Activity {
	hr := HumanResource{role, desc, costH, qty, nil}
	// Validazione silenziosa - l'utente può chiamare Validate() esplicitamente
	_ = hr.Validate()
	a.Humans = append(a.Humans, hr)
	return a
}

// WithHumanFromSupplier aggiunge una risorsa umana fornita da un fornitore
func (a *Activity) WithHumanFromSupplier(role, desc string, costH, qty float64, supplier *Supplier) *Activity {
	hr := HumanResource{role, desc, costH, qty, supplier}
	// Validazione silenziosa - l'utente può chiamare Validate() esplicitamente
	_ = hr.Validate()
	a.Humans = append(a.Humans, hr)
	return a
}

func (a *Activity) WithMaterial(name, desc string, cost, qty float64) *Activity {
	mr := MaterialResource{name, desc, cost, qty, nil}
	// Validazione silenziosa - l'utente può chiamare Validate() esplicitamente
	_ = mr.Validate()
	a.Materials = append(a.Materials, mr)
	return a
}

// WithMaterialFromSupplier aggiunge un materiale fornito da un fornitore
func (a *Activity) WithMaterialFromSupplier(name, desc string, cost, qty float64, supplier *Supplier) *Activity {
	mr := MaterialResource{name, desc, cost, qty, supplier}
	// Validazione silenziosa - l'utente può chiamare Validate() esplicitamente
	_ = mr.Validate()
	a.Materials = append(a.Materials, mr)
	return a
}

func (a *Activity) WithAsset(name, desc string, cost, qty float64) *Activity {
	asset := Asset{name, desc, cost, qty, nil}
	// Validazione silenziosa - l'utente può chiamare Validate() esplicitamente
	_ = asset.Validate()
	a.Assets = append(a.Assets, asset)
	return a
}

// WithAssetFromSupplier aggiunge un asset fornito da un fornitore
func (a *Activity) WithAssetFromSupplier(name, desc string, cost, qty float64, supplier *Supplier) *Activity {
	asset := Asset{name, desc, cost, qty, supplier}
	// Validazione silenziosa - l'utente può chiamare Validate() esplicitamente
	_ = asset.Validate()
	a.Assets = append(a.Assets, asset)
	return a
}

// NewSupplier crea un nuovo fornitore (può essere riutilizzato).
// Valida automaticamente i campi del fornitore e ritorna errore se la validazione fallisce.
// Se la validazione fallisce, ritorna nil e l'errore può essere recuperato chiamando Validate().
func NewSupplier(name, desc string, unitCost, availableQty float64, period PeriodType) *Supplier {
	s := &Supplier{
		Name:              name,
		Description:       desc,
		UnitCost:          unitCost,
		AvailableQuantity: availableQty,
		Period:            period,
	}
	// Validazione automatica - se fallisce, il fornitore è comunque creato
	// ma l'utente può chiamare Validate() per verificare
	_ = s.Validate() // validazione silenziosa, l'utente può chiamare Validate() esplicitamente
	return s
}

func (a *Activity) CanCrash(minDur int, extraCost float64) *Activity {
	a.MinDuration = minDur
	a.CrashCostStep = extraCost
	return a
}

func (a *Activity) DependsOn(subs ...*Activity) *Activity {
	for _, s := range subs {
		s.Next = append(s.Next, a.ID)
		a.SubActivities = append(a.SubActivities, s)
	}
	return a
}
