// internal/repository/message.go
package repository

import (
	"kakao-clone/internal/model"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(msg *model.Message) error
	FindByRoomID(roomID string, limit int) ([]model.Message, error)
	MarkAsRead(roomID string, userID uint) error
	UnreadCount(roomID string, userID uint) (int64, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *model.Message) error {
	return r.db.Create(msg).Error
}

func (r *messageRepository) FindByRoomID(roomID string, limit int) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.
		Where("room_id = ?", roomID).
		Preload("Sender").
		Order("created_at asc").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) MarkAsRead(roomID string, userID uint) error {
	return r.db.Model(&model.Message{}).
		Where("room_id = ? AND sender_id != ? AND is_read = false", roomID, userID).
		Update("is_read", true).Error
}

func (r *messageRepository) UnreadCount(roomID string, userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Message{}).
		Where("room_id = ? AND sender_id != ? AND is_read = false", roomID, userID).
		Count(&count).Error
	return count, err
}
