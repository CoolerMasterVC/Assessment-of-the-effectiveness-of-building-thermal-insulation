// internal/app/repository/material.go
package repository

import (
	"Lab1/internal/app/ds"
	"strings"
)

func (r *Repository) GetMaterials() ([]*ds.Material, error) {
	var materials []*ds.Material
	err := r.db.Where("status = ?", "действует").Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}

func (r *Repository) GetMaterialByID(id uint) (*ds.Material, error) {
	var material ds.Material
	err := r.db.Where("id = ? AND status = ?", id, "действует").First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *Repository) SearchMaterials(query string) ([]*ds.Material, error) {
	var materials []*ds.Material
	searchQuery := "%" + strings.ToLower(query) + "%"
	err := r.db.Where("(LOWER(name) LIKE ? OR LOWER(description) LIKE ?) AND status = ?",
		searchQuery, searchQuery, "действует").Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}
