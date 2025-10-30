package calculations

import "Lab1/internal/app/ds"

// CalculateHeatLoss рассчитывает теплопотери для материала
// material - теплоизоляционный материал
// area - площадь утепления в м²
// indoorTemp - температура внутри помещения в °C
// outdoorTemp - температура снаружи в °C
// возвращает теплопотери в Ваттах
func CalculateHeatLoss(material ds.Material, area float64, indoorTemp, outdoorTemp float64) float64 {
	// Термическое сопротивление материала
	thermalResistance := material.Thickness / material.Lambda
	if thermalResistance == 0 {
		return 0
	}

	// Разница температур
	deltaT := indoorTemp - outdoorTemp

	// Теплопотери по формуле: Q = (ΔT * S) / R
	// где ΔT - разница температур, S - площадь, R - термическое сопротивление
	heatLoss := (area * deltaT) / thermalResistance
	return heatLoss
}

// CalculateSavings рассчитывает экономию от использования теплоизоляции
// heatLossBefore - теплопотери до утепления (Вт)
// heatLossAfter - теплопотери после утепления (Вт)
// energyCost - стоимость энергии (руб/кВт·ч)
// возвращает экономию в руб/месяц
func CalculateSavings(heatLossBefore, heatLossAfter float64, energyCost float64) float64 {
	// Снижение теплопотерь
	heatLossReduction := heatLossBefore - heatLossAfter

	if heatLossReduction <= 0 {
		return 0
	}

	// Перевод в киловатт-часы за месяц:
	// 1 Вт = 1 Дж/с
	// За 1 час: 1 Вт * 3600 с = 3600 Дж = 0.001 кВт·ч
	// За месяц (30 дней): 0.001 кВт·ч/час * 24 часа * 30 дней

	monthlyEnergySavings := (heatLossReduction * 24 * 30) / 1000 // кВт·ч/месяц
	monthlySavings := monthlyEnergySavings * energyCost          // руб/месяц

	return monthlySavings
}

// CalculateTotalSavings рассчитывает общую экономию для заявки
// application - заявка на расчет
// appMaterials - материалы в заявке
// возвращает общую экономию в руб/месяц
func CalculateTotalSavings(application *ds.MaterialsApplication, appMaterials []ds.ApplicationMaterial) float64 {
	if application.Status != "завершён" {
		return 0
	}

	totalSavings := 0.0
	energyCost := 5.0 // руб/кВт·ч

	// Базовая стена: кирпич толщиной 0.5м с коэффициентом теплопроводности 0.7 Вт/(м·K)
	baseWall := ds.Material{
		Lambda:    0.7,
		Thickness: 0.5,
	}

	for _, appMaterial := range appMaterials {
		material := appMaterial.Material

		// Теплопотери базовой стены (до утепления)
		baseHeatLoss := CalculateHeatLoss(
			baseWall,
			appMaterial.Area,
			application.IndoorTemp,
			application.OutdoorTemp,
		)

		// Теплопотери утепленной стены (после утепления)
		insulatedHeatLoss := CalculateHeatLoss(
			material,
			appMaterial.Area,
			application.IndoorTemp,
			application.OutdoorTemp,
		)

		// Экономия для этого материала
		materialSavings := CalculateSavings(baseHeatLoss, insulatedHeatLoss, energyCost)
		totalSavings += materialSavings
	}

	return totalSavings
}

// CalculateHeatLossReduction рассчитывает снижение теплопотерь для одного материала
// material - теплоизоляционный материал
// area - площадь утепления
// indoorTemp - температура внутри
// outdoorTemp - температура снаружи
// возвращает снижение теплопотерь в Ваттах
func CalculateHeatLossReduction(material ds.Material, area float64, indoorTemp, outdoorTemp float64) float64 {
	// Базовая стена (кирпич)
	baseWall := ds.Material{Lambda: 0.7, Thickness: 0.5}

	baseHeatLoss := CalculateHeatLoss(baseWall, area, indoorTemp, outdoorTemp)
	insulatedHeatLoss := CalculateHeatLoss(material, area, indoorTemp, outdoorTemp)

	return baseHeatLoss - insulatedHeatLoss
}

// CalculatePaybackPeriod рассчитывает срок окупаемости теплоизоляции
// materialCost - стоимость материала (руб)
// monthlySavings - месячная экономия (руб/месяц)
// возвращает срок окупаемости в месяцах
func CalculatePaybackPeriod(materialCost, monthlySavings float64) float64 {
	if monthlySavings <= 0 {
		return 0
	}
	return materialCost / monthlySavings
}
