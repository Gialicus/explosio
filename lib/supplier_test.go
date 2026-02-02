package lib

import "testing"

func TestValidateSupplierUsage_Valid(t *testing.T) {
	engine := &AnalysisEngine{}
	supplier := &Supplier{
		Name:              "Test Supplier",
		AvailableQuantity: 100,
		Period:            PeriodDay,
	}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 50, Supplier: supplier},
	}
	errors := engine.ValidateSupplierUsage(root)
	if len(errors) != 0 {
		t.Errorf("ValidateSupplierUsage valid: want 0 errors, got %d", len(errors))
	}
}

func TestValidateSupplierUsage_ExceedsCapacity(t *testing.T) {
	engine := &AnalysisEngine{}
	supplier := &Supplier{
		Name:              "Test Supplier",
		AvailableQuantity: 50,
		Period:            PeriodDay,
	}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 100, Supplier: supplier},
	}
	errors := engine.ValidateSupplierUsage(root)
	if len(errors) != 1 {
		t.Fatalf("ValidateSupplierUsage exceeds capacity: want 1 error, got %d", len(errors))
	}
	if errors[0] == nil {
		t.Error("error should not be nil")
	}
}

func TestValidateSupplierUsage_WithHuman(t *testing.T) {
	engine := &AnalysisEngine{}
	supplier := &Supplier{
		Name:              "Agency",
		AvailableQuantity: 5,
		Period:            PeriodDay,
	}
	root := buildActivity("A", "A", 1, nil)
	root.Humans = []HumanResource{
		{Role: "Worker", CostPerH: 20, Quantity: 3, Supplier: supplier},
	}
	errors := engine.ValidateSupplierUsage(root)
	if len(errors) != 0 {
		t.Errorf("ValidateSupplierUsage with human: want 0 errors, got %d", len(errors))
	}
}

func TestValidateSupplierUsage_WithAsset(t *testing.T) {
	engine := &AnalysisEngine{}
	supplier := &Supplier{
		Name:              "Rental",
		AvailableQuantity: 10,
		Period:            PeriodWeek,
	}
	root := buildActivity("A", "A", 1, nil)
	root.Assets = []Asset{
		{Name: "Machine", CostPerUse: 50, Quantity: 2, Supplier: supplier},
	}
	errors := engine.ValidateSupplierUsage(root)
	if len(errors) != 0 {
		t.Errorf("ValidateSupplierUsage with asset: want 0 errors, got %d", len(errors))
	}
}

func TestValidateSupplierUsage_Tree(t *testing.T) {
	engine := &AnalysisEngine{}
	supplier := &Supplier{
		Name:              "Supplier",
		AvailableQuantity: 10,
		Period:            PeriodDay,
	}
	child := buildActivity("B", "B", 1, nil)
	child.Materials = []MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 5, Supplier: supplier},
	}
	root := buildActivity("A", "A", 1, []*Activity{child})
	root.Materials = []MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 6, Supplier: supplier},
	}
	errors := engine.ValidateSupplierUsage(root)
	if len(errors) != 0 {
		t.Errorf("ValidateSupplierUsage tree: want 0 errors, got %d", len(errors))
	}
}

func TestCalculateSupplierRequirements_Empty(t *testing.T) {
	engine := &AnalysisEngine{}
	root := buildActivity("A", "A", 1, nil)
	requirements := engine.CalculateSupplierRequirements(root, 100, PeriodDay)
	if len(requirements) != 0 {
		t.Errorf("CalculateSupplierRequirements empty: want 0, got %d", len(requirements))
	}
}

func TestCalculateSupplierRequirements_WithMaterial(t *testing.T) {
	engine := &AnalysisEngine{}
	supplier := &Supplier{
		Name:              "Test Supplier",
		AvailableQuantity: 1000,
		Period:            PeriodDay,
	}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 10, Supplier: supplier},
	}
	requirements := engine.CalculateSupplierRequirements(root, 50, PeriodDay)
	if len(requirements) != 1 {
		t.Fatalf("CalculateSupplierRequirements: want 1 requirement, got %d", len(requirements))
	}
	req := requirements[0]
	if req.SupplierName != "Test Supplier" {
		t.Errorf("SupplierName want 'Test Supplier', got '%s'", req.SupplierName)
	}
	// 50 prodotti * 10 unità = 500 unità/giorno
	assertFloatEqual(t, req.RequiredQuantity, 500)
	assertFloatEqual(t, req.SuppliersNeeded, 0.5) // 500 / 1000
	if !req.IsFeasible {
		t.Error("IsFeasible want true")
	}
}

func TestCalculateSupplierRequirements_PeriodConversion(t *testing.T) {
	engine := &AnalysisEngine{}
	// Fornitore con capacità settimanale
	supplier := &Supplier{
		Name:              "Weekly Supplier",
		AvailableQuantity: 700, // 700 unità/settimana
		Period:            PeriodWeek,
	}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 10, Supplier: supplier},
	}
	// Produzione target: 100 prodotti/giorno
	requirements := engine.CalculateSupplierRequirements(root, 100, PeriodDay)
	if len(requirements) != 1 {
		t.Fatalf("CalculateSupplierRequirements: want 1 requirement, got %d", len(requirements))
	}
	req := requirements[0]
	// 100 prodotti/giorno * 10 unità = 1000 unità/giorno
	// Convertito in settimana: 1000 * 7 = 7000 unità/settimana
	assertFloatEqual(t, req.RequiredQuantity, 7000)
	assertFloatEqual(t, req.SuppliersNeeded, 10.0) // 7000 / 700
}

func TestCalculateSupplierRequirements_MultipleResources(t *testing.T) {
	engine := &AnalysisEngine{}
	supplier1 := &Supplier{
		Name:              "Supplier 1",
		AvailableQuantity: 1000,
		Period:            PeriodDay,
	}
	supplier2 := &Supplier{
		Name:              "Supplier 2",
		AvailableQuantity: 500,
		Period:            PeriodDay,
	}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{
		{Name: "Material 1", UnitCost: 1, Quantity: 10, Supplier: supplier1},
		{Name: "Material 2", UnitCost: 1, Quantity: 5, Supplier: supplier2},
	}
	requirements := engine.CalculateSupplierRequirements(root, 50, PeriodDay)
	if len(requirements) != 2 {
		t.Fatalf("CalculateSupplierRequirements: want 2 requirements, got %d", len(requirements))
	}
	// Verifica che entrambi i fornitori siano presenti
	supplierMap := make(map[string]SupplierRequirement)
	for _, req := range requirements {
		supplierMap[req.SupplierName] = req
	}
	if _, exists := supplierMap["Supplier 1"]; !exists {
		t.Error("Supplier 1 should be in requirements")
	}
	if _, exists := supplierMap["Supplier 2"]; !exists {
		t.Error("Supplier 2 should be in requirements")
	}
}

func TestCalculateSupplierRequirements_NilRoot(t *testing.T) {
	engine := &AnalysisEngine{}
	requirements := engine.CalculateSupplierRequirements(nil, 100, PeriodDay)
	if requirements != nil {
		t.Errorf("CalculateSupplierRequirements(nil): want nil, got %v", requirements)
	}
}

func TestSupplier_GetCapacityForPeriod(t *testing.T) {
	supplier := &Supplier{
		Name:              "Test",
		AvailableQuantity: 100,
		Period:            PeriodDay,
	}
	// Stesso periodo
	capacity := supplier.GetCapacityForPeriod(PeriodDay)
	assertFloatEqual(t, capacity, 100)
	// Conversione giorno -> settimana
	capacity = supplier.GetCapacityForPeriod(PeriodWeek)
	assertFloatEqual(t, capacity, 700) // 100 * 7
	// Conversione giorno -> mese (30 giorni)
	capacity = supplier.GetCapacityForPeriod(PeriodMonth)
	assertFloatEqual(t, capacity, 3000) // 100 * 30
}

func TestSupplier_GetDailyCapacity(t *testing.T) {
	supplier := &Supplier{
		Name:              "Test",
		AvailableQuantity: 700,
		Period:            PeriodWeek,
	}
	capacity := supplier.GetDailyCapacity()
	assertFloatEqual(t, capacity, 100) // 700 / 7
}
