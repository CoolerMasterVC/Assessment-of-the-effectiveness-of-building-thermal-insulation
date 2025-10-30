package repository

import (
	"Lab1/internal/app/ds"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Material methods
func (r *Repository) GetMaterials() ([]*ds.Material, error) {
	var materials []*ds.Material
	err := r.db.Where("status = ?", "действует").Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}

func (r *Repository) GetMaterialsWithFilter(filter string) ([]ds.Material, error) {
	var materials []ds.Material
	query := r.db.Where("status != ?", "удалён")

	if filter != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+filter+"%", "%"+filter+"%")
	}

	err := query.Find(&materials).Error
	return materials, err
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

func (r *Repository) CreateMaterial(material *ds.Material) error {
	material.Status = "действует"
	return r.db.Create(material).Error
}

func (r *Repository) UpdateMaterial(id uint, material *ds.Material) error {
	return r.db.Model(&ds.Material{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":         material.Name,
		"description":  material.Description,
		"price_per_m2": material.PricePerM2,
		"lambda":       material.Lambda,
		"thickness":    material.Thickness,
	}).Error
}

func (r *Repository) DeleteMaterial(id uint) error {
	return r.db.Model(&ds.Material{}).Where("id = ?", id).Update("status", "удалён").Error
}

func (r *Repository) UpdateMaterialImage(id uint, imageURL string) error {
	return r.db.Model(&ds.Material{}).Where("id = ?", id).Update("image_url", imageURL).Error
}

// MaterialsApplication methods
func (r *Repository) GetUserDraft(userID uint) (*ds.MaterialsApplication, error) {
	var application ds.MaterialsApplication
	err := r.db.Where("creator_id = ? AND status = ?", userID, "черновик").First(&application).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &application, err
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

func (r *Repository) GetApplications(status string, startDate, endDate *time.Time) ([]ds.MaterialsApplication, error) {
	var applications []ds.MaterialsApplication
	query := r.db.Preload("Creator").Preload("Moderator").
		Where("status != ? AND status != ?", "удалён", "черновик")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	err := query.Find(&applications).Error
	return applications, err
}

func (r *Repository) GetApplicationByID(id uint) (*ds.MaterialsApplication, error) {
	var application ds.MaterialsApplication
	err := r.db.Preload("Materials.Material").
		Preload("Creator").
		Preload("Moderator").
		Where("id = ? AND status != ?", id, "удалён").
		First(&application).Error
	return &application, err
}

func (r *Repository) UpdateApplication(id uint, application *ds.MaterialsApplication) error {
	return r.db.Model(&ds.MaterialsApplication{}).Where("id = ?", id).Updates(map[string]interface{}{
		"total_area":   application.TotalArea,
		"indoor_temp":  application.IndoorTemp,
		"outdoor_temp": application.OutdoorTemp,
	}).Error
}

func (r *Repository) SubmitApplication(id uint) error {
	now := time.Now()
	return r.db.Model(&ds.MaterialsApplication{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "сформирован",
		"submitted_at": &now,
	}).Error
}

func (r *Repository) CompleteApplication(id, moderatorID uint, totalSavings float64) error {
	now := time.Now()
	return r.db.Model(&ds.MaterialsApplication{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        "завершён",
		"completed_at":  &now,
		"moderator_id":  moderatorID,
		"total_savings": totalSavings,
	}).Error
}

func (r *Repository) RejectApplication(id, moderatorID uint) error {
	now := time.Now()
	return r.db.Model(&ds.MaterialsApplication{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "отклонён",
		"completed_at": &now,
		"moderator_id": moderatorID,
	}).Error
}

func (r *Repository) DeleteApplication(id uint) error {
	return r.db.Model(&ds.MaterialsApplication{}).Where("id = ?").Update("status", "удалён").Error
}

// ApplicationMaterial methods
func (r *Repository) AddMaterialToApplication(appID, materialID uint, area float64) error {
	var count int64
	r.db.Model(&ds.ApplicationMaterial{}).
		Where("application_id = ? AND material_id = ?", appID, materialID).
		Count(&count)

	if count > 0 {
		return r.db.Model(&ds.ApplicationMaterial{}).
			Where("application_id = ? AND material_id = ?", appID, materialID).
			Update("area", area).Error
	}

	appMaterial := ds.ApplicationMaterial{
		ApplicationID: appID,
		MaterialID:    materialID,
		Area:          area,
	}
	err := r.db.Create(&appMaterial).Error
	if err != nil {
		return err
	}

	// Пересчитываем общую площадь
	return r.RecalculateApplicationArea(appID)
}

func (r *Repository) RemoveMaterialFromApplication(appID, materialID uint) error {
	err := r.db.Where("application_id = ? AND material_id = ?", appID, materialID).
		Delete(&ds.ApplicationMaterial{}).Error
	if err != nil {
		return err
	}

	// Пересчитываем общую площадь
	return r.RecalculateApplicationArea(appID)
}

func (r *Repository) UpdateApplicationMaterial(appID, materialID uint, area float64) error {
	err := r.db.Model(&ds.ApplicationMaterial{}).
		Where("application_id = ? AND material_id = ?", appID, materialID).
		Update("area", area).Error
	if err != nil {
		return err
	}

	// Пересчитываем общую площадь
	return r.RecalculateApplicationArea(appID)
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

// User methods
func (r *Repository) CreateUser(user *ds.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	var user ds.User
	err := r.db.Where("login = ?", login).First(&user).Error
	return &user, err
}

func (r *Repository) UpdateUser(id uint, user *ds.User) error {
	return r.db.Model(&ds.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"login": user.Login,
	}).Error
}
