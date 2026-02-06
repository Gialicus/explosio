package engine

import (
	"errors"
	"explosio/lib/domain"
	"testing"
)

func TestValidateSupplierUsage_Valid(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier := &domain.Supplier{
		Name:              "Test Supplier",
		AvailableQuantity: 100,
		Period:            domain.PeriodDay,
	}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	root.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 50, Supplier: supplier},
	}
	err := eng.ValidateSupplierUsage(root)
	if err != nil {
		t.Errorf("ValidateSupplierUsage valid: want nil, got %v", err)
	}
}

func TestValidateSupplierUsage_ExceedsCapacity(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier := &domain.Supplier{
		Name:              "Test Supplier",
		AvailableQuantity: 50,
		Period:            domain.PeriodDay,
	}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	root.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 100, Supplier: supplier},
	}
	err := eng.ValidateSupplierUsage(root)
	if err == nil {
		t.Fatal("ValidateSupplierUsage exceeds capacity: want error, got nil")
	}
	var ve *domain.ValidationErrors
	if !errors.As(err, &ve) || len(ve.Errors) != 1 {
		t.Errorf("ValidateSupplierUsage exceeds capacity: want ValidationErrors with 1 error, got %v", err)
	}
}

func TestValidateSupplierUsage_WithHuman(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier := &domain.Supplier{
		Name:              "Agency",
		AvailableQuantity: 5,
		Period:            domain.PeriodDay,
	}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	root.Humans = []domain.HumanResource{
		{Role: "Worker", CostPerH: 20, Quantity: 3, Supplier: supplier},
	}
	err := eng.ValidateSupplierUsage(root)
	if err != nil {
		t.Errorf("ValidateSupplierUsage with human: want nil, got %v", err)
	}
}

func TestValidateSupplierUsage_WithAsset(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier := &domain.Supplier{
		Name:              "Rental",
		AvailableQuantity: 10,
		Period:            domain.PeriodWeek,
	}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	root.Assets = []domain.Asset{
		{Name: "Machine", CostPerUse: 50, Quantity: 2, Supplier: supplier},
	}
	err := eng.ValidateSupplierUsage(root)
	if err != nil {
		t.Errorf("ValidateSupplierUsage with asset: want nil, got %v", err)
	}
}

func TestValidateSupplierUsage_Tree(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier := &domain.Supplier{
		Name:              "Supplier",
		AvailableQuantity: 10,
		Period:            domain.PeriodDay,
	}
	child := domain.BuildActivityForTest("B", "B", 1, nil)
	child.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 5, Supplier: supplier},
	}
	root := domain.BuildActivityForTest("A", "A", 1, []*domain.Activity{child})
	root.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 6, Supplier: supplier},
	}
	err := eng.ValidateSupplierUsage(root)
	if err != nil {
		t.Errorf("ValidateSupplierUsage tree: want nil, got %v", err)
	}
}

func TestCalculateSupplierRequirements_Empty(t *testing.T) {
	eng := &AnalysisEngine{}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	requirements := eng.CalculateSupplierRequirements(root, 100, domain.PeriodDay)
	if len(requirements) != 0 {
		t.Errorf("CalculateSupplierRequirements empty: want 0, got %d", len(requirements))
	}
}

func TestCalculateSupplierRequirements_WithMaterial(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier := &domain.Supplier{
		Name:              "Test Supplier",
		AvailableQuantity: 1000,
		Period:            domain.PeriodDay,
	}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	root.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 10, Supplier: supplier},
	}
	requirements := eng.CalculateSupplierRequirements(root, 50, domain.PeriodDay)
	if len(requirements) != 1 {
		t.Fatalf("CalculateSupplierRequirements: want 1 requirement, got %d", len(requirements))
	}
	req := requirements[0]
	if req.SupplierName != "Test Supplier" {
		t.Errorf("SupplierName want 'Test Supplier', got '%s'", req.SupplierName)
	}
	assertFloatEqual(t, req.RequiredQuantity, 500)
	assertFloatEqual(t, req.SuppliersNeeded, 0.5)
	if !req.IsFeasible {
		t.Error("IsFeasible want true")
	}
}

func TestCalculateSupplierRequirements_PeriodConversion(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier := &domain.Supplier{
		Name:              "Weekly Supplier",
		AvailableQuantity: 700,
		Period:            domain.PeriodWeek,
	}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	root.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 10, Supplier: supplier},
	}
	requirements := eng.CalculateSupplierRequirements(root, 100, domain.PeriodDay)
	if len(requirements) != 1 {
		t.Fatalf("CalculateSupplierRequirements: want 1 requirement, got %d", len(requirements))
	}
	req := requirements[0]
	assertFloatEqual(t, req.RequiredQuantity, 7000)
	assertFloatEqual(t, req.SuppliersNeeded, 10.0)
}

func TestCalculateSupplierRequirements_MultipleResources(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier1 := &domain.Supplier{
		Name:              "Supplier 1",
		AvailableQuantity: 1000,
		Period:            domain.PeriodDay,
	}
	supplier2 := &domain.Supplier{
		Name:              "Supplier 2",
		AvailableQuantity: 500,
		Period:            domain.PeriodDay,
	}
	root := domain.BuildActivityForTest("A", "A", 1, nil)
	root.Materials = []domain.MaterialResource{
		{Name: "Material 1", UnitCost: 1, Quantity: 10, Supplier: supplier1},
		{Name: "Material 2", UnitCost: 1, Quantity: 5, Supplier: supplier2},
	}
	requirements := eng.CalculateSupplierRequirements(root, 50, domain.PeriodDay)
	if len(requirements) != 2 {
		t.Fatalf("CalculateSupplierRequirements: want 2 requirements, got %d", len(requirements))
	}
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
	eng := &AnalysisEngine{}
	requirements := eng.CalculateSupplierRequirements(nil, 100, domain.PeriodDay)
	if requirements != nil {
		t.Errorf("CalculateSupplierRequirements(nil): want nil, got %v", requirements)
	}
}

func TestSupplier_GetCapacityForPeriod(t *testing.T) {
	supplier := &domain.Supplier{
		Name:              "Test",
		AvailableQuantity: 100,
		Period:            domain.PeriodDay,
	}
	capacity := supplier.GetCapacityForPeriod(domain.PeriodDay)
	assertFloatEqual(t, capacity, 100)
	capacity = supplier.GetCapacityForPeriod(domain.PeriodWeek)
	assertFloatEqual(t, capacity, 700)
	capacity = supplier.GetCapacityForPeriod(domain.PeriodMonth)
	assertFloatEqual(t, capacity, 3000)
}

func TestSupplier_GetDailyCapacity(t *testing.T) {
	supplier := &domain.Supplier{
		Name:              "Test",
		AvailableQuantity: 700,
		Period:            domain.PeriodWeek,
	}
	capacity := supplier.GetDailyCapacity()
	assertFloatEqual(t, capacity, 100)
}

func TestSupplier_Validate_NegativeQuantity(t *testing.T) {
	supplier := &domain.Supplier{
		Name:              "Test",
		AvailableQuantity: -10,
		Period:            domain.PeriodDay,
	}
	err := supplier.Validate()
	if err == nil {
		t.Fatal("Validate() should return error for negative quantity")
	}
	if !errors.Is(err, domain.ErrNegativeQuantity) {
		t.Errorf("Validate() should return ErrNegativeQuantity, got %v", err)
	}
}

func TestSupplier_Validate_InvalidPeriod(t *testing.T) {
	supplier := &domain.Supplier{
		Name:              "Test",
		AvailableQuantity: 100,
		Period:            domain.PeriodType("invalid"),
	}
	err := supplier.Validate()
	if err == nil {
		t.Fatal("Validate() should return error for invalid period")
	}
	if !errors.Is(err, domain.ErrInvalidPeriod) {
		t.Errorf("Validate() should return ErrInvalidPeriod, got %v", err)
	}
}

func TestHumanResource_Validate_NegativeCost(t *testing.T) {
	hr := domain.HumanResource{
		Role:     "Worker",
		CostPerH: -10,
		Quantity: 1,
	}
	err := hr.Validate()
	if err == nil {
		t.Fatal("Validate() should return error for negative cost")
	}
	if !errors.Is(err, domain.ErrNegativeCost) {
		t.Errorf("Validate() should return ErrNegativeCost, got %v", err)
	}
}

func TestMaterialResource_Validate_NegativeQuantity(t *testing.T) {
	mr := domain.MaterialResource{
		Name:     "Material",
		UnitCost: 10,
		Quantity: -5,
	}
	err := mr.Validate()
	if err == nil {
		t.Fatal("Validate() should return error for negative quantity")
	}
	if !errors.Is(err, domain.ErrNegativeQuantity) {
		t.Errorf("Validate() should return ErrNegativeQuantity, got %v", err)
	}
}

func TestValidateSupplierUsage_DeepTree(t *testing.T) {
	eng := &AnalysisEngine{}
	supplier := &domain.Supplier{
		Name:              "Test",
		AvailableQuantity: 100,
		Period:            domain.PeriodDay,
	}
	deep := domain.BuildActivityForTest("D", "Deep", 1, nil)
	deep.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 10, Supplier: supplier},
	}
	mid := domain.BuildActivityForTest("M", "Mid", 1, []*domain.Activity{deep})
	mid.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 20, Supplier: supplier},
	}
	root := domain.BuildActivityForTest("R", "Root", 1, []*domain.Activity{mid})
	root.Materials = []domain.MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 30, Supplier: supplier},
	}
	err := eng.ValidateSupplierUsage(root)
	if err != nil {
		t.Errorf("ValidateSupplierUsage deep tree: want nil, got %v", err)
	}
}
