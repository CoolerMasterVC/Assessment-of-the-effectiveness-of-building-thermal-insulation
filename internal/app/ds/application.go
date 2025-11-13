package ds

import "time"

type MaterialsApplication struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Status       string     `gorm:"default:'черновик';check:status IN ('черновик', 'удалён', 'сформирован', 'завершён', 'отклонён')" json:"status"`
	CreatorID    uint       `gorm:"not null" json:"creator_id"`
	TotalArea    float64    `gorm:"not null;default:0" json:"total_area"`
	IndoorTemp   float64    `gorm:"not null;default:22" json:"indoor_temp"`
	OutdoorTemp  float64    `gorm:"not null;default:-15" json:"outdoor_temp"`
	TotalSavings float64    `gorm:"default:0" json:"total_savings"` // ДОЛЖНО БЫТЬ ЭТО ПОЛЕ
	CreatedAt    time.Time  `gorm:"not null;default:current_timestamp" json:"created_at"`
	SubmittedAt  *time.Time `gorm:"default:null" json:"submitted_at"`
	CompletedAt  *time.Time `gorm:"default:null" json:"completed_at"`
	ModeratorID  *uint      `gorm:"default:null" json:"moderator_id"`

	Creator   User                  `gorm:"foreignKey:CreatorID" json:"creator"`
	Moderator User                  `gorm:"foreignKey:ModeratorID" json:"moderator"`
	Materials []ApplicationMaterial `gorm:"foreignKey:ApplicationID" json:"materials"`
}
