// internal/model/room.go
package model

import "time"

// Room — 채팅방 (1:1 · 그룹 통합)
type Room struct {
	ID        string    `gorm:"primaryKey"      json:"id"`
	Name      string    `gorm:"default:''"      json:"name"` // 그룹방 이름
	Type      string    `gorm:"not null"        json:"type"` // "direct" | "group"
	CreatedAt time.Time `                       json:"createdAt"`
	UpdatedAt time.Time `                       json:"updatedAt"`

	// 연관관계
	Members  []RoomMember `gorm:"foreignKey:RoomID" json:"members,omitempty"`
	Messages []Message    `gorm:"foreignKey:RoomID" json:"messages,omitempty"`
}

// RoomMember — 채팅방 멤버 (다대다 중간 테이블)
type RoomMember struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID   string    `gorm:"not null;index"           json:"roomId"`
	UserID   uint      `gorm:"not null;index"           json:"userId"`
	JoinedAt time.Time `                                json:"joinedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
