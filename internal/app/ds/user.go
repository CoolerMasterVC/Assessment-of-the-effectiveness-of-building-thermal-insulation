// internal/app/ds/user.go
package ds

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Login       string    `gorm:"unique;not null;size:150" json:"login"`
	Password    string    `gorm:"not null;size:128" json:"-"`
	IsModerator bool      `gorm:"default:false" json:"is_moderator"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
