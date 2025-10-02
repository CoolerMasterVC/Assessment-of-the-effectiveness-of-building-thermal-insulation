// internal/calculations/heat_calculator.go
package calculations

import "Lab1/internal/app/ds"

func CalculateHeatLoss(material ds.Material, area float64, indoorTemp, outdoorTemp float64) float64 {
	thermalResistance := material.Thickness / material.Lambda
	if thermalResistance == 0 {
		return 0
	}
	// Это МГНОВЕННЫЕ теплопотери в Ваттах
	heatLoss := (area * (indoorTemp - outdoorTemp)) / thermalResistance
	return heatLoss
}

func CalculateSavings(heatLossBefore, heatLossAfter float64, energyCost float64) float64 {
	heatLossReduction := heatLossBefore - heatLossAfter

	// ПРАВИЛЬНЫЙ расчет: переводим в киловатт-часы за месяц
	// heatLossReduction в Ваттах, поэтому:
	// 1 Вт = 1 Дж/с
	// За 1 час: 1 Вт * 3600 с = 3600 Дж = 0.001 кВт·ч
	// За месяц (30 дней): 0.001 кВт·ч/час * 24 часа * 30 дней

	monthlyEnergySavings := (heatLossReduction * 24 * 30) / 1000 // кВт·ч/месяц
	monthlySavings := monthlyEnergySavings * energyCost          // руб/месяц

	return monthlySavings
}

func CalculateTotalSavings(application *ds.MaterialsApplication, appMaterials []ds.ApplicationMaterial) float64 {
	if application.Status != "завершён" {
		return 0
	}

	totalSavings := 0.0
	energyCost := 5.0 // руб/кВт·ч

	// Базовая стена: кирпич 0.5м
	baseWall := ds.Material{Lambda: 0.7, Thickness: 0.5}

	for _, appMaterial := range appMaterials {
		material := appMaterial.Material

		baseHeatLoss := CalculateHeatLoss(baseWall, appMaterial.Area,
			application.IndoorTemp, application.OutdoorTemp)
		insulatedHeatLoss := CalculateHeatLoss(material, appMaterial.Area,
			application.IndoorTemp, application.OutdoorTemp)

		heatLossReduction := baseHeatLoss - insulatedHeatLoss

		// ПРАВИЛЬНЫЙ расчет экономии
		monthlyEnergySavings := (heatLossReduction * 24 * 30) / 1000 // кВт·ч/месяц
		monthlySavings := monthlyEnergySavings * energyCost          // руб/месяц

		totalSavings += monthlySavings
	}

	return totalSavings
}
