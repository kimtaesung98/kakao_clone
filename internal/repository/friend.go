// internal/repository/friend.go
package repository

import (
	"kakao-clone/internal/model"

	"gorm.io/gorm"
)

type FriendRepository interface {
	Add(userID, friendID uint) error
	Accept(userID, friendID uint) error
	FindByUserID(userID uint) ([]model.Friend, error)
	IsFriend(userID, friendID uint) bool
	Delete(userID, friendID uint) error
}

type friendRepository struct {
	db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) FriendRepository {
	return &friendRepository{db: db}
}

func (r *friendRepository) Add(userID, friendID uint) error {
	friend := model.Friend{
		UserID:   userID,
		FriendID: friendID,
		Status:   "pending",
	}
	return r.db.Create(&friend).Error
}

func (r *friendRepository) Accept(userID, friendID uint) error {
	return r.db.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ?", friendID, userID).
		Update("status", "accepted").Error
}

func (r *friendRepository) FindByUserID(userID uint) ([]model.Friend, error) {
	var friends []model.Friend
	err := r.db.
		Where("user_id = ? AND status = ?", userID, "accepted").
		Preload("FriendUser").
		Find(&friends).Error
	return friends, err
}

func (r *friendRepository) IsFriend(userID, friendID uint) bool {
	var count int64
	r.db.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ? AND status = ?",
			userID, friendID, "accepted").
		Count(&count)
	return count > 0
}

func (r *friendRepository) Delete(userID, friendID uint) error {
	return r.db.
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Delete(&model.Friend{}).Error
}
