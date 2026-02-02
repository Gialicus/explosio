package lib

import (
	"math"
	"testing"
)

const floatEpsilon = 1e-9

// assertFloatEqual verifica che due float64 siano uguali entro un epsilon
func assertFloatEqual(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > floatEpsilon {
		t.Errorf("got %v, want %v", got, want)
	}
}

// buildActivity crea un'attività minimale per test (senza DSL)
func buildActivity(id, name string, duration int, subs []*Activity) *Activity {
	a := &Activity{
		ID:            id,
		Name:          name,
		Description:   name,
		Duration:      duration,
		MinDuration:   duration,
		SubActivities: subs,
	}
	for _, s := range subs {
		s.Next = append(s.Next, a.ID)
	}
	return a
}
