package converter

import "fmt"

var temperatureUnits = map[string]float64{
	"C": 1,
	"F": 5.0 / 9.0,
	"K": 1.0 / 273.15,
}

func ConvertTemperature(value float64, fromUnit string, toUnit string) (float64, error) {
	unitFrom, ok1 := temperatureUnits[fromUnit]
	unitTo, ok2 := temperatureUnits[toUnit]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid units")
	}

	switch unitFrom {
	case temperatureUnits["C"]:
		switch unitTo {
		case temperatureUnits["F"]:
			return celsiusToFahrenheit(value), nil
		case temperatureUnits["K"]:
			return celsiusToKelvin(value), nil
		}
	case temperatureUnits["F"]:
		switch unitTo {
		case temperatureUnits["C"]:
			return fahrenheitToCelsius(value), nil
		case temperatureUnits["K"]:
			return fahrenheitToKelvin(value), nil
		}
	case temperatureUnits["K"]:
		switch unitTo {
		case temperatureUnits["C"]:
			return kelvinToCelsius(value), nil
		case temperatureUnits["F"]:
			return kelvinToFahrenheit(value), nil
		}
	}

	return 0, fmt.Errorf("invalid conversion")
}

func GetTemperatureUnits() []string {
	units := make([]string, 0, len(temperatureUnits))
	for unit := range temperatureUnits {
		units = append(units, unit)
	}

	return units
}

// CelsiusToFahrenheit converts Celsius to Fahrenheit
func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*9/5 + 32
}

// FahrenheitToCelsius converts Fahrenheit to Celsius
func fahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}

// CelsiusToKelvin converts Celsius to Kelvin
func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

// KelvinToCelsius converts Kelvin to Celsius
func kelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}

// FahrenheitToKelvin converts Fahrenheit to Kelvin
func fahrenheitToKelvin(fahrenheit float64) float64 {
	return (fahrenheit-32)*5/9 + 273.15
}

// KelvinToFahrenheit converts Kelvin to Fahrenheit
func kelvinToFahrenheit(kelvin float64) float64 {
	return (kelvin-273.15)*9/5 + 32
}
