package unit

type MeasurableUnit string

const (
	// Unit of measurement for length
	UnitMillimeter MeasurableUnit = "mm"
	UnitCentimeter MeasurableUnit = "cm"
	UnitDecimeter  MeasurableUnit = "dm"
	UnitMeter      MeasurableUnit = "m"
	UnitKilometer  MeasurableUnit = "km"
	// Unit of measurement for area
	UnitSquareMillimeter MeasurableUnit = "mm²"
	UnitSquareCentimeter MeasurableUnit = "cm²"
	UnitSquareDecimeter  MeasurableUnit = "dm²"
	UnitSquareMeter      MeasurableUnit = "m²"
	UnitSquareKilometer  MeasurableUnit = "km²"
	// Unit of measurement for volume
	UnitCubicMillimeter MeasurableUnit = "mm³"
	UnitCubicCentimeter MeasurableUnit = "cm³"
	UnitCubicDecimeter  MeasurableUnit = "dm³"
	UnitCubicMeter      MeasurableUnit = "m³"
	UnitCubicKilometer  MeasurableUnit = "km³"
	// Unit of measurement for weight
	UnitMilligram MeasurableUnit = "mg"
	UnitCentigram MeasurableUnit = "cg"
	UnitDecigram  MeasurableUnit = "dg"
	UnitGram      MeasurableUnit = "g"
	UnitKilogram  MeasurableUnit = "kg"
	UnitTon       MeasurableUnit = "t"
)
