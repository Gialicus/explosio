package lib

import "fmt"

// Errori tipizzati per validazione e operazioni

var (
	// ErrInvalidPeriod indica che un PeriodType non è valido
	ErrInvalidPeriod = fmt.Errorf("invalid period type")

	// ErrNegativeQuantity indica che una quantità è negativa
	ErrNegativeQuantity = fmt.Errorf("quantity cannot be negative")

	// ErrNegativeCost indica che un costo è negativo
	ErrNegativeCost = fmt.Errorf("cost cannot be negative")

	// ErrInvalidActivity indica che un'attività non è valida
	ErrInvalidActivity = fmt.Errorf("invalid activity")

	// ErrInvalidSupplier indica che un fornitore non è valido
	ErrInvalidSupplier = fmt.Errorf("invalid supplier")
)

// ValidationErrors raccoglie multiple errori di validazione
type ValidationErrors struct {
	Errors []error
}

// Error implementa l'interfaccia error
func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "no validation errors"
	}
	if len(ve.Errors) == 1 {
		return ve.Errors[0].Error()
	}
	return fmt.Sprintf("%d validation errors: %v", len(ve.Errors), ve.Errors[0])
}

// Unwrap restituisce gli errori sottostanti (per errors.Is e errors.As)
func (ve *ValidationErrors) Unwrap() []error {
	return ve.Errors
}

// Add aggiunge un errore alla collezione
func (ve *ValidationErrors) Add(err error) {
	if err != nil {
		ve.Errors = append(ve.Errors, err)
	}
}

// HasErrors restituisce true se ci sono errori
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}
