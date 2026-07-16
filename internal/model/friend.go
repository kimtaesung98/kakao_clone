// internal/model/friend.go
package model

import "time"

type Friend struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index"           json:"userId"`
	FriendID  uint      `gorm:"not null;index"           json:"friendId"`
	Status    string    `gorm:"default:'pending'"        json:"status"` // "pending" | "accepted"
	CreatedAt time.Time `                                json:"createdAt"`

	FriendUser User `gorm:"foreignKey:FriendID" json:"friendUser,omitempty"`
}
