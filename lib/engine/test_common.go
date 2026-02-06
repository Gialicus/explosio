package engine

import (
	"math"
	"testing"
)

const floatEpsilon = 1e-9

func assertFloatEqual(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > floatEpsilon {
		t.Errorf("got %v, want %v", got, want)
	}
}
