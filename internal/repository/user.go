// internal/repository/user.go
package repository

import (
	"kakao-clone/internal/model"

	"gorm.io/gorm"
)

// ─────────────────────────────────────────
// UserRepository interface — DB 접근 약속
// interface로 선언하면 나중에 테스트할 때
// 가짜(Mock) 구현체로 교체 가능
// ─────────────────────────────────────────
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	FindAll() ([]model.User, error)
	Search(keyword string) ([]model.User, error)
}

// ─────────────────────────────────────────
// userRepository — 실제 GORM 구현체
// ─────────────────────────────────────────
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository — 생성자
// Node.js: class UserRepository { constructor(db) { this.db = db } }
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Search(keyword string) ([]model.User, error) {
	var users []model.User
	err := r.db.Where(
		"name LIKE ? OR email LIKE ?",
		"%"+keyword+"%",
		"%"+keyword+"%",
	).Find(&users).Error
	return users, err
}
