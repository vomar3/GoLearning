package main

import "fmt"

func createTemperatureAdjuster() (func(change float64) float64, float64) {
	baseTemperature := 90.0

	adjustTemperature := func(change float64) float64 {
		baseTemperature += change
		return baseTemperature
	}

	return adjustTemperature, baseTemperature
}

func main() {
	// Сложно, но пробуем

	// Тут мы вызываем функцию, после чего мы НИ РАЗУ не вызываем ее с самого начала
	// По сути, значение baseTemperature должно перестать существовать в памяти, как локалка
	// Но т.к. мы эту переменную используем в последующей функции adjustTemperature, то
	// Значение переменной baseTemperature остается в памяти (измененное)
	// Из-за этого мы при дальнейшем вызове adjustTemp изменяем значение baseTemperature
	adjustTemp, originalTemp := createTemperatureAdjuster()
	fmt.Printf("Original temperature is %.1f\n", originalTemp)
	fmt.Printf("Adjusted Temp +1.5: %.1fC\n", adjustTemp(1.5))
	fmt.Printf("Adjusted Temp -3.0: %.1fC\n", adjustTemp(-3.0))
	fmt.Printf("Adjusted Temp +5.0: %.1fC\n", adjustTemp(5.0))
	fmt.Printf("Original temperature is %.1f\n", originalTemp) // Не меняется в прицнипе

	// 90.0
	// +1.5 -> 91.5
	// -3.0 -> 88.5
	// +5.0 -> 93.5
}
