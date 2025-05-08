package converter

import "fmt"

var weightUnits = map[string]float64{
	"mg": 0.001,
	"g":  1,
	"kg": 1000,
	"lb": 453.592,
	"oz": 28.3495,
}

func ConvertWeight(value float64, fromUnit string, toUnit string) (float64, error) {
	// convert from unit to grams
	weightFrom, ok1 := weightUnits[fromUnit]
	weightTo, ok2 := weightUnits[toUnit]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid units")
	}

	return value * weightFrom / weightTo, nil
}

func GetWeightUnits() []string {
	units := make([]string, 0, len(weightUnits))

	for unit := range weightUnits {
		units = append(units, unit)
	}

	return units
}
