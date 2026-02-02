package ui

import (
	"explosio/lib"
	"sync"
)

// AppState gestisce lo stato globale dell'applicazione
type AppState struct {
	mu        sync.RWMutex
	project   *lib.Project
	suppliers []*lib.Supplier
	engine    *lib.AnalysisEngine
	
	// Callbacks per aggiornare l'UI quando cambia lo stato
	onProjectChanged []func(*lib.Project)
	onSuppliersChanged []func([]*lib.Supplier)
}

// NewAppState crea un nuovo AppState
func NewAppState() *AppState {
	return &AppState{
		project:   lib.NewProject(),
		suppliers: make([]*lib.Supplier, 0),
		engine:    &lib.AnalysisEngine{},
		onProjectChanged: make([]func(*lib.Project), 0),
		onSuppliersChanged: make([]func([]*lib.Supplier), 0),
	}
}

// GetProject restituisce il progetto corrente
func (a *AppState) GetProject() *lib.Project {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.project
}

// SetProject imposta il progetto corrente e notifica i listener
func (a *AppState) SetProject(project *lib.Project) {
	a.mu.Lock()
	a.project = project
	a.mu.Unlock()
	
	// Notifica i listener
	for _, callback := range a.onProjectChanged {
		callback(project)
	}
}

// GetSuppliers restituisce la lista dei fornitori
func (a *AppState) GetSuppliers() []*lib.Supplier {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]*lib.Supplier, len(a.suppliers))
	copy(result, a.suppliers)
	return result
}

// AddSupplier aggiunge un fornitore
func (a *AppState) AddSupplier(supplier *lib.Supplier) {
	a.mu.Lock()
	a.suppliers = append(a.suppliers, supplier)
	suppliers := make([]*lib.Supplier, len(a.suppliers))
	copy(suppliers, a.suppliers)
	a.mu.Unlock()
	
	// Notifica i listener
	for _, callback := range a.onSuppliersChanged {
		callback(suppliers)
	}
}

// RemoveSupplier rimuove un fornitore per nome
func (a *AppState) RemoveSupplier(name string) bool {
	a.mu.Lock()
	found := false
	newSuppliers := make([]*lib.Supplier, 0, len(a.suppliers))
	for _, s := range a.suppliers {
		if s.Name != name {
			newSuppliers = append(newSuppliers, s)
		} else {
			found = true
		}
	}
	a.suppliers = newSuppliers
	suppliers := make([]*lib.Supplier, len(a.suppliers))
	copy(suppliers, a.suppliers)
	a.mu.Unlock()
	
	if found {
		// Notifica i listener
		for _, callback := range a.onSuppliersChanged {
			callback(suppliers)
		}
	}
	return found
}

// GetSupplierByName restituisce un fornitore per nome
func (a *AppState) GetSupplierByName(name string) *lib.Supplier {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, s := range a.suppliers {
		if s.Name == name {
			return s
		}
	}
	return nil
}

// SetSuppliers imposta la lista completa dei fornitori
func (a *AppState) SetSuppliers(suppliers []*lib.Supplier) {
	a.mu.Lock()
	a.suppliers = make([]*lib.Supplier, len(suppliers))
	copy(a.suppliers, suppliers)
	suppliersCopy := make([]*lib.Supplier, len(suppliers))
	copy(suppliersCopy, suppliers)
	a.mu.Unlock()
	
	// Notifica i listener
	for _, callback := range a.onSuppliersChanged {
		callback(suppliersCopy)
	}
}

// GetEngine restituisce il motore di analisi
func (a *AppState) GetEngine() *lib.AnalysisEngine {
	return a.engine
}

// OnProjectChanged registra un callback che viene chiamato quando cambia il progetto
func (a *AppState) OnProjectChanged(callback func(*lib.Project)) {
	a.mu.Lock()
	a.onProjectChanged = append(a.onProjectChanged, callback)
	a.mu.Unlock()
}

// OnSuppliersChanged registra un callback che viene chiamato quando cambiano i fornitori
func (a *AppState) OnSuppliersChanged(callback func([]*lib.Supplier)) {
	a.mu.Lock()
	a.onSuppliersChanged = append(a.onSuppliersChanged, callback)
	a.mu.Unlock()
}

// ComputeCPM calcola il CPM per il progetto corrente
func (a *AppState) ComputeCPM() {
	project := a.GetProject()
	if project != nil && project.Root != nil {
		a.engine.ComputeCPM(project.Root)
		// Notifica che il progetto è cambiato (per aggiornare le visualizzazioni)
		a.SetProject(project)
	}
}
