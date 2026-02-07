package unit

import "fmt"

type Price struct {
	Value    float64
	Currency string
}

func NewPrice(value float64, currency string) *Price {
	return &Price{Value: value, Currency: currency}
}

func (p *Price) String() string {
	return fmt.Sprintf("%.2f %s", p.Value, p.Currency)
}
