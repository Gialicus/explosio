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
	a.Humans = append(a.Humans, HumanResource{role, desc, costH, qty})
	return a
}

func (a *Activity) WithMaterial(name, desc string, cost, qty float64) *Activity {
	a.Materials = append(a.Materials, MaterialResource{name, desc, cost, qty})
	return a
}

func (a *Activity) WithAsset(name, desc string, cost, qty float64) *Activity {
	a.Assets = append(a.Assets, Asset{name, desc, cost, qty})
	return a
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
