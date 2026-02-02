package lib

import "testing"

// buildDeepTree crea un albero di attività con profondità depth e larghezza width
func buildDeepTree(depth, width int) *Activity {
	if depth == 0 {
		return buildActivity("LEAF", "Leaf", 1, nil)
	}
	subs := make([]*Activity, width)
	for i := 0; i < width; i++ {
		subs[i] = buildDeepTree(depth-1, width)
	}
	return buildActivity("NODE", "Node", 1, subs)
}

// buildWideTree crea un albero largo con molte attività allo stesso livello
func buildWideTree(count int) *Activity {
	subs := make([]*Activity, count)
	for i := 0; i < count; i++ {
		subs[i] = buildActivity("ACT", "Activity", 1, nil)
	}
	return buildActivity("ROOT", "Root", 1, subs)
}

func BenchmarkComputeCPM_Shallow(b *testing.B) {
	root := buildWideTree(10)
	engine := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.ComputeCPM(root)
	}
}

func BenchmarkComputeCPM_Deep(b *testing.B) {
	root := buildDeepTree(5, 2)
	engine := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.ComputeCPM(root)
	}
}

func BenchmarkComputeCPM_Wide(b *testing.B) {
	root := buildWideTree(100)
	engine := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.ComputeCPM(root)
	}
}

func BenchmarkGetTotalCost_Simple(b *testing.B) {
	root := buildActivity("A", "A", 10, nil)
	root.Humans = []HumanResource{{Role: "Worker", CostPerH: 20, Quantity: 1}}
	root.Materials = []MaterialResource{{Name: "Material", UnitCost: 5, Quantity: 2}}
	engine := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetTotalCost(root)
	}
}

func BenchmarkGetTotalCost_Complex(b *testing.B) {
	root := buildDeepTree(4, 3)
	// Aggiungi risorse a tutte le attività
	var addResources func(*Activity)
	addResources = func(a *Activity) {
		a.Humans = []HumanResource{{Role: "Worker", CostPerH: 20, Quantity: 1}}
		a.Materials = []MaterialResource{{Name: "Material", UnitCost: 5, Quantity: 2}}
		for _, sub := range a.SubActivities {
			addResources(sub)
		}
	}
	addResources(root)
	engine := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetTotalCost(root)
	}
}

func BenchmarkCalculateSupplierRequirements_Simple(b *testing.B) {
	supplier := &Supplier{
		Name:              "Test",
		AvailableQuantity: 1000,
		Period:            PeriodDay,
	}
	root := buildActivity("A", "A", 1, nil)
	root.Materials = []MaterialResource{
		{Name: "Material", UnitCost: 1, Quantity: 10, Supplier: supplier},
	}
	engine := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.CalculateSupplierRequirements(root, 100, PeriodDay)
	}
}

func BenchmarkCalculateSupplierRequirements_Multiple(b *testing.B) {
	suppliers := make([]*Supplier, 10)
	for i := 0; i < 10; i++ {
		suppliers[i] = &Supplier{
			Name:              "Supplier",
			AvailableQuantity: 1000,
			Period:            PeriodDay,
		}
	}
	root := buildWideTree(10)
	// Aggiungi fornitori a tutte le attività
	var addSuppliers func(*Activity, int)
	addSuppliers = func(a *Activity, idx int) {
		a.Materials = []MaterialResource{
			{Name: "Material", UnitCost: 1, Quantity: 10, Supplier: suppliers[idx%10]},
		}
		for i, sub := range a.SubActivities {
			addSuppliers(sub, (idx+i)%10)
		}
	}
	addSuppliers(root, 0)
	engine := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.CalculateSupplierRequirements(root, 100, PeriodDay)
	}
}

func BenchmarkGetCriticalPath(b *testing.B) {
	root := buildDeepTree(4, 3)
	engine := &AnalysisEngine{}
	engine.ComputeCPM(root)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetCriticalPath(root)
	}
}

func BenchmarkActivitiesByES(b *testing.B) {
	root := buildWideTree(50)
	engine := &AnalysisEngine{}
	engine.ComputeCPM(root)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.ActivitiesByES(root)
	}
}
