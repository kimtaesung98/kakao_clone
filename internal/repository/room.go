// internal/repository/room.go
package repository

import (
	"kakao-clone/internal/model"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *model.Room) error
	FindByID(id string) (*model.Room, error)
	FindDirectRoom(userA, userB uint) (*model.Room, error)
	FindByUserID(userID uint) ([]model.Room, error)
	AddMember(roomID string, userID uint) error
	RemoveMember(roomID string, userID uint) error
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) FindByID(id string) (*model.Room, error) {
	var room model.Room
	err := r.db.Preload("Members.User").First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindDirectRoom — 두 유저 사이의 1:1 방 찾기
func (r *roomRepository) FindDirectRoom(userA, userB uint) (*model.Room, error) {
	var room model.Room
	err := r.db.
		Joins("JOIN room_members rm1 ON rm1.room_id = rooms.id AND rm1.user_id = ?", userA).
		Joins("JOIN room_members rm2 ON rm2.room_id = rooms.id AND rm2.user_id = ?", userB).
		Where("rooms.type = ?", "direct").
		First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindByUserID — 유저가 속한 모든 방 조회
func (r *roomRepository) FindByUserID(userID uint) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.
		Joins("JOIN room_members ON room_members.room_id = rooms.id").
		Where("room_members.user_id = ?", userID).
		Preload("Members.User").
		Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) AddMember(roomID string, userID uint) error {
	member := model.RoomMember{RoomID: roomID, UserID: userID}
	return r.db.Create(&member).Error
}

func (r *roomRepository) RemoveMember(roomID string, userID uint) error {
	return r.db.Where(
		"room_id = ? AND user_id = ?", roomID, userID,
	).Delete(&model.RoomMember{}).Error
}
