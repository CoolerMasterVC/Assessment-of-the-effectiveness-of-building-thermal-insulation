package repository

import (
	"Lab1/internal/models"
	"fmt"
	"strings"
)

type Repository struct {
	materials []models.Material
	cart      models.Cart
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
	}

	cart := models.Cart{
		ID:    1,
		Items: []models.CartItem{},
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

func (r *Repository) GetCart() (models.Cart, error) {
	totalSavings := 0.0
	for _, item := range r.cart.Items {
		material, _ := r.GetMaterialByID(item.MaterialID)
		monthlySavings := item.Area * (0.1 / material.Lambda) * 24 * 30 * 0.5
		totalSavings += monthlySavings
	}
	r.cart.TotalSavings = totalSavings
	return r.cart, nil
}

func (r *Repository) AddToCart(materialID int, area float64) error {
	for i, item := range r.cart.Items {
		if item.MaterialID == materialID {
			r.cart.Items[i].Area += area
			return nil
		}
	}

	r.cart.Items = append(r.cart.Items, models.CartItem{
		MaterialID: materialID,
		Area:       area,
	})
	return nil
}
