package unit

import "testing"

func TestNewPrice(t *testing.T) {
	p := NewPrice(10.50, "EUR")
	if p == nil {
		t.Fatal("NewPrice returned nil")
	}
	if p.Value != 10.50 || p.Currency != "EUR" {
		t.Errorf("NewPrice(10.50, EUR) = %+v", p)
	}
}

func TestPrice_String(t *testing.T) {
	p := NewPrice(10.5, "EUR")
	got := p.String()
	if got != "10.50 EUR" {
		t.Errorf("String() = %q, want \"10.50 EUR\"", got)
	}
}
