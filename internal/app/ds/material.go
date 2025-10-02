// internal/app/ds/material.go
package ds

type Material struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null;size:255"`
	Description string  `gorm:"type:text"`
	Status      string  `gorm:"default:'действует';check:status IN ('действует', 'удалён')"`
	ImageURL    string  `gorm:"size:500"`
	PricePerM2  float64 `gorm:"not null;default:0"`
	Lambda      float64 `gorm:"not null;default:0"` // коэффициент теплопроводности
	Thickness   float64 `gorm:"not null;default:0"` // толщина материала
}
