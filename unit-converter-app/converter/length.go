package converter

import "fmt"

var lengthUnits = map[string]float64{
	"mm": 0.001,
	"cm": 0.01,
	"m":  1,
	"km": 1000,
	"in": 0.0254,
	"ft": 0.3048,
	"yd": 0.9144,
	"mi": 1609.34,
}

func ConvertLength(value float64, fromUnit string, toUnit string) (float64, error) {
	// convert from unit to meters
	fromRate, ok1 := lengthUnits[fromUnit]
	toRate, ok2 := lengthUnits[toUnit]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid units")
	}

	// first convert from unit to meters then to target unit
	meters := value * fromRate
	result := meters / toRate

	return result, nil
}

func GetLengthUnits() []string {
	units := make([]string, 0, len(lengthUnits))

	for unit := range lengthUnits {
		units = append(units, unit)
	}

	return units
}
