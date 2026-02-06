package domain

import (
	"fmt"
	"sync/atomic"
)

// Project rappresenta un progetto con albero di attività e generazione ID.
type Project struct {
	idCounter int64
	Root      *Activity
}

// NewProject crea un nuovo progetto vuoto.
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

func (a *Activity) WithHuman(role, desc string, costH, qty float64) *Activity {
	hr := HumanResource{role, desc, costH, qty}
	_ = hr.Validate()
	a.Humans = append(a.Humans, hr)
	return a
}

func (a *Activity) WithMaterial(name, desc string, cost, qty float64) *Activity {
	mr := MaterialResource{name, desc, cost, qty}
	_ = mr.Validate()
	a.Materials = append(a.Materials, mr)
	return a
}

func (a *Activity) WithAsset(name, desc string, cost, qty float64) *Activity {
	asset := Asset{name, desc, cost, qty}
	_ = asset.Validate()
	a.Assets = append(a.Assets, asset)
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
