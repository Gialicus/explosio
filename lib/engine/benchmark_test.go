package engine

import (
	"explosio/lib/domain"
	"testing"
)

func buildDeepTree(depth, width int) *domain.Activity {
	if depth == 0 {
		return domain.BuildActivityForTest("LEAF", "Leaf", 1, nil)
	}
	subs := make([]*domain.Activity, width)
	for i := 0; i < width; i++ {
		subs[i] = buildDeepTree(depth-1, width)
	}
	return domain.BuildActivityForTest("NODE", "Node", 1, subs)
}

func buildWideTree(count int) *domain.Activity {
	subs := make([]*domain.Activity, count)
	for i := 0; i < count; i++ {
		subs[i] = domain.BuildActivityForTest("ACT", "Activity", 1, nil)
	}
	return domain.BuildActivityForTest("ROOT", "Root", 1, subs)
}

func BenchmarkComputeCPM_Shallow(b *testing.B) {
	root := buildWideTree(10)
	eng := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.ComputeCPM(root)
	}
}

func BenchmarkComputeCPM_Deep(b *testing.B) {
	root := buildDeepTree(5, 2)
	eng := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.ComputeCPM(root)
	}
}

func BenchmarkComputeCPM_Wide(b *testing.B) {
	root := buildWideTree(100)
	eng := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.ComputeCPM(root)
	}
}

func BenchmarkGetTotalCost_Simple(b *testing.B) {
	root := domain.BuildActivityForTest("A", "A", 10, nil)
	root.Humans = []domain.HumanResource{{Role: "Worker", CostPerH: 20, Quantity: 1}}
	root.Materials = []domain.MaterialResource{{Name: "Material", UnitCost: 5, Quantity: 2}}
	eng := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eng.GetTotalCost(root)
	}
}

func BenchmarkGetTotalCost_Complex(b *testing.B) {
	root := buildDeepTree(4, 3)
	var addResources func(*domain.Activity)
	addResources = func(a *domain.Activity) {
		a.Humans = []domain.HumanResource{{Role: "Worker", CostPerH: 20, Quantity: 1}}
		a.Materials = []domain.MaterialResource{{Name: "Material", UnitCost: 5, Quantity: 2}}
		for _, sub := range a.SubActivities {
			addResources(sub)
		}
	}
	addResources(root)
	eng := &AnalysisEngine{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eng.GetTotalCost(root)
	}
}

func BenchmarkGetCriticalPath(b *testing.B) {
	root := buildDeepTree(4, 3)
	eng := &AnalysisEngine{}
	eng.ComputeCPM(root)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eng.GetCriticalPath(root)
	}
}

func BenchmarkActivitiesByES(b *testing.B) {
	root := buildWideTree(50)
	eng := &AnalysisEngine{}
	eng.ComputeCPM(root)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eng.ActivitiesByES(root)
	}
}
