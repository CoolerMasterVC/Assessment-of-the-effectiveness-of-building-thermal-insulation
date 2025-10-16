// internal/app/repository/application.go
package repository

import (
	"Lab1/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetUserDraft(userID uint) (*ds.MaterialsApplication, error) {
	var application ds.MaterialsApplication
	err := r.db.Where("creator_id = ? AND status = ?", userID, "черновик").First(&application).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *Repository) CreateDraft(userID uint) (*ds.MaterialsApplication, error) {
	application := ds.MaterialsApplication{
		Status:      "черновик",
		CreatorID:   userID,
		TotalArea:   0,
		IndoorTemp:  22,
		OutdoorTemp: -15,
	}
	err := r.db.Create(&application).Error
	return &application, err
}

func (r *Repository) AddMaterialToApplication(appID, materialID uint, area float64) error {
	var count int64
	r.db.Model(&ds.ApplicationMaterial{}).
		Where("application_id = ? AND material_id = ?", appID, materialID).
		Count(&count)

	if count > 0 {
		err := r.db.Model(&ds.ApplicationMaterial{}).
			Where("application_id = ? AND material_id = ?", appID, materialID).
			Update("area", area).Error
		if err != nil {
			return err
		}
	} else {
		appMaterial := ds.ApplicationMaterial{
			ApplicationID: appID,
			MaterialID:    materialID,
			Area:          area,
		}
		err := r.db.Create(&appMaterial).Error
		if err != nil {
			return err
		}
	}

	// Пересчитываем общую площадь
	return r.RecalculateApplicationArea(appID)
}

func (r *Repository) GetApplicationByID(id uint) (*ds.MaterialsApplication, error) {
	var application ds.MaterialsApplication
	err := r.db.Where("id = ?", id).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *Repository) GetApplicationMaterials(appID uint) ([]ds.ApplicationMaterial, error) {
	var appMaterials []ds.ApplicationMaterial
	err := r.db.
		Preload("Material").
		Where("application_id = ?", appID).
		Find(&appMaterials).Error
	if err != nil {
		return nil, err
	}
	return appMaterials, nil
}

func (r *Repository) DeleteApplication(appID uint) error {
	return r.db.Exec("UPDATE materials_applications SET status = 'удалён' WHERE id = ?", appID).Error
}

func (r *Repository) GetApplicationMaterialsCount(appID uint) int64 {
	var count int64
	r.db.Model(&ds.ApplicationMaterial{}).Where("application_id = ?", appID).Count(&count)
	return count
}

func (r *Repository) UpdateApplicationArea(appID uint, area float64) error {
	return r.db.Model(&ds.MaterialsApplication{}).
		Where("id = ?", appID).
		Update("total_area", area).Error
}

func (r *Repository) RecalculateApplicationArea(appID uint) error {
	var totalArea float64
	err := r.db.Model(&ds.ApplicationMaterial{}).
		Where("application_id = ?", appID).
		Select("COALESCE(SUM(area), 0)").
		Scan(&totalArea).Error
	if err != nil {
		return err
	}

	return r.db.Model(&ds.MaterialsApplication{}).
		Where("id = ?", appID).
		Update("total_area", totalArea).Error
}
