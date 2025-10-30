// internal/app/ds/material.go
package ds

import "time"

type Material struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;size:255" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"default:'действует';check:status IN ('действует', 'удалён')" json:"status"`
	ImageURL    string    `gorm:"size:500" json:"image_url"`
	PricePerM2  float64   `gorm:"not null;default:0" json:"price_per_m2"`
	Lambda      float64   `gorm:"not null;default:0" json:"lambda"`
	Thickness   float64   `gorm:"not null;default:0" json:"thickness"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
