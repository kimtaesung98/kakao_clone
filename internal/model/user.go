// internal/model/user.go
package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null"     json:"email"`
	Password  *string   `gorm:"default:null"             json:"-"`
	Name      string    `gorm:"not null"                 json:"name"`
	Avatar    string    `gorm:"default:''"               json:"avatar"`    // 프로필 이미지
	StatusMsg string    `gorm:"default:''"               json:"statusMsg"` // 상태메시지
	Provider  string    `gorm:"default:'local'"          json:"provider"`
	GoogleID  string    `gorm:"default:''"               json:"googleId,omitempty"`
	CreatedAt time.Time `                                json:"createdAt"`
	UpdatedAt time.Time `                                json:"updatedAt"`
}
