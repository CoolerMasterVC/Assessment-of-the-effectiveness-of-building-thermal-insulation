package repository

import (
	"Lab1/internal/models"
	"fmt"
	"strings"
)

type Repository struct {
	materials []models.Material
	cart      models.ApplicationMaterials
}

func NewRepository() (*Repository, error) {
	materials := []models.Material{
		{
			ID:          1,
			Name:        "Минеральная вата",
			PricePerM2:  650,
			Lambda:      0.045,
			Description: "Эффективный утеплитель с отличными звукоизоляционными свойствами",
			ImageURL:    "http://localhost:9000/images/mineral_wool.jpg",
		},
		{
			ID:          2,
			Name:        "Пенополистирол",
			PricePerM2:  450,
			Lambda:      0.038,
			Description: "Легкий и влагостойкий материал для утепления",
			ImageURL:    "http://localhost:9000/images/polystyrene.jpg",
		},
		{
			ID:          3,
			Name:        "PIR-плиты",
			PricePerM2:  1200,
			Lambda:      0.028,
			Description: "Современный высокоэффективный утеплитель с низкой теплопроводностью",
			ImageURL:    "http://localhost:9000/images/pir_plates.jpg",
		},
		{
			ID:          4,
			Name:        "Экструдированный пенополистирол",
			PricePerM2:  550,
			Lambda:      0.033,
			Description: "Прочный влагостойкий утеплитель для фундаментов и фасадов",
			ImageURL:    "http://localhost:9000/images/xps.jpg",
		},
	}

	cart := models.ApplicationMaterials{
		ID:    1,
		Items: []models.CartMaterial{},
	}

	return &Repository{
		materials: materials,
		cart:      cart,
	}, nil
}

func (r *Repository) GetAllMaterials() ([]models.Material, error) {
	return r.materials, nil
}

func (r *Repository) GetMaterialByID(id int) (models.Material, error) {
	for _, material := range r.materials {
		if material.ID == id {
			return material, nil
		}
	}
	return models.Material{}, fmt.Errorf("материал не найден")
}

func (r *Repository) GetMaterialsByName(name string) ([]models.Material, error) {
	var result []models.Material
	for _, material := range r.materials {
		if strings.Contains(strings.ToLower(material.Name), strings.ToLower(name)) {
			result = append(result, material)
		}
	}
	return result, nil
}

func (r *Repository) GetCart(id int) (models.ApplicationMaterials, error) {
	cart := models.ApplicationMaterials{
		ID:          id,
		IndoorTemp:  22.0,
		OutdoorTemp: -15.0,
		Items: []models.CartMaterial{
			{MaterialID: 1, Area: 15.5},
			{MaterialID: 3, Area: 22.0},
		},
	}

	totalArea := 0.0
	for _, item := range cart.Items {
		totalArea += item.Area
	}
	cart.TotalArea = totalArea

	// Статичное значение вместо расчёта
	cart.TotalSavings = 4900.0

	return cart, nil
}
