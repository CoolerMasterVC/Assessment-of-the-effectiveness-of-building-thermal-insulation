// internal/app/ds/application_material.go
package ds

import "time"

type ApplicationMaterial struct {
	ApplicationID uint      `gorm:"primaryKey" json:"-"`
	MaterialID    uint      `gorm:"primaryKey" json:"material_id"`
	Area          float64   `gorm:"not null" json:"area"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Application MaterialsApplication `gorm:"foreignKey:ApplicationID" json:"-"`
	Material    Material             `gorm:"foreignKey:MaterialID" json:"material"`
}
