// internal/app/ds/application.go
package ds

import "time"

type MaterialsApplication struct {
	ID           uint       `gorm:"primaryKey"`
	Status       string     `gorm:"default:'черновик';check:status IN ('черновик', 'удалён', 'сформирован', 'завершён', 'отклонён')"`
	CreatorID    uint       `gorm:"not null"`
	TotalArea    float64    `gorm:"not null;default:0"` // общая площадь утепления
	IndoorTemp   float64    `gorm:"not null;default:22"`
	OutdoorTemp  float64    `gorm:"not null;default:-15"`
	CreatedAt    time.Time  `gorm:"not null;default:current_timestamp"`
	SubmittedAt  *time.Time `gorm:"default:null"`
	CompletedAt  *time.Time `gorm:"default:null"`
	ModeratorID  *uint      `gorm:"default:null"`
	TotalSavings float64    `gorm:"default:0"` // расчетная экономия

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`
}
