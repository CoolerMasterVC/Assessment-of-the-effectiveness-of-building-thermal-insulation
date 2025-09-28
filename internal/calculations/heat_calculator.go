// internal/calculations/heat_calculator.go
package calculations

import "Lab1/internal/app/ds"

func CalculateHeatLoss(material ds.Material, area float64, indoorTemp, outdoorTemp float64) float64 {
	thermalResistance := material.Thickness / material.Lambda
	if thermalResistance == 0 {
		return 0
	}
	heatLoss := (area * (indoorTemp - outdoorTemp)) / thermalResistance
	return heatLoss
}

func CalculateSavings(heatLossBefore, heatLossAfter float64, energyCost float64) float64 {
	heatLossReduction := heatLossBefore - heatLossAfter
	monthlySavings := (heatLossReduction * 24 * 30 * energyCost) / 1000
	return monthlySavings
}

func CalculateTotalSavings(application *ds.Application, appMaterials []ds.ApplicationMaterial) float64 {
	if application.Status != "завершён" {
		return 0
	}

	totalSavings := 0.0
	energyCost := 5.0

	for _, appMaterial := range appMaterials {
		material := appMaterial.Material
		baseHeatLoss := CalculateHeatLoss(ds.Material{Lambda: 0.7, Thickness: 0.5},
			appMaterial.Area, application.IndoorTemp, application.OutdoorTemp)
		insulatedHeatLoss := CalculateHeatLoss(material, appMaterial.Area,
			application.IndoorTemp, application.OutdoorTemp)

		savings := CalculateSavings(baseHeatLoss, insulatedHeatLoss, energyCost)
		totalSavings += savings
	}

	return totalSavings
}
