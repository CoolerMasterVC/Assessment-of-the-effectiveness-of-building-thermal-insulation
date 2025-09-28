// internal/app/ds/application_material.go
package ds

type ApplicationMaterial struct {
	ApplicationID uint    `gorm:"primaryKey"`
	MaterialID    uint    `gorm:"primaryKey"`
	Area          float64 `gorm:"not null"` // площадь данного материала

	Application Application `gorm:"foreignKey:ApplicationID"`
	Material    Material    `gorm:"foreignKey:MaterialID"`
}
