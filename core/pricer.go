package core

// Pricer is implemented by any type that can calculate its price.
// Activity, Asset, HumanResource, ComplexMaterial, CountableMaterial, and MeasurableMaterial implement this interface.
type Pricer interface {
	CalculatePrice() float64
}
