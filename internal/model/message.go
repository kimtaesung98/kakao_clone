// internal/model/message.go
package model

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"not null"                 json:"content"`
	Type      string    `gorm:"default:'text'"           json:"type"` // "text" | "image"
	IsRead    bool      `gorm:"default:false"            json:"isRead"`
	RoomID    string    `gorm:"not null;index"           json:"roomId"`
	SenderID  uint      `gorm:"not null"                 json:"senderId"`
	CreatedAt time.Time `                                json:"createdAt"`

	Sender User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
