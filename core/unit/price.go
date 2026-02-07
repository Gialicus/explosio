package unit

import "fmt"

// Price represents an amount with currency.
type Price struct {
	Value    float64
	Currency string
}

// NewPrice creates a price with value and currency code.
func NewPrice(value float64, currency string) *Price {
	return &Price{Value: value, Currency: currency}
}

// String formats the price for output (e.g. "10.50 EUR").
func (p *Price) String() string {
	return fmt.Sprintf("%.2f %s", p.Value, p.Currency)
}
