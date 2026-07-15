// internal/service/user.go
package service

import (
	"fmt"

	"gorm.io/gorm"

	"kakao-clone/internal/db"
	"kakao-clone/internal/model"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// GetByID — 유저 조회
// 에러 종류에 따라 다른 AppError 반환
func (s *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	result := db.DB.First(&user, id)

	if result.Error != nil {
		// 종류 구분
		if result.Error == gorm.ErrRecordNotFound {
			// 404
			return nil, ErrNotFound("유저를 찾을 수 없습니다")
		}
		// 500 — 원본 에러 보존
		return nil, ErrInternal("유저 조회 실패", result.Error)
	}

	return &user, nil
}

// GetByEmail — 이메일로 조회
func (s *UserService) GetByEmail(email string) (*model.User, error) {
	var user model.User
	result := db.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound("등록되지 않은 이메일입니다")
		}
		return nil, ErrInternal("유저 조회 실패", result.Error)
	}

	return &user, nil
}

// AddFriend — 친구 추가
// 에러 wrapping 예시
func (s *UserService) AddFriend(userID, friendID uint) error {
	// 본인 확인
	if userID == friendID {
		return ErrBadRequest("자기 자신을 친구 추가할 수 없습니다")
	}

	// 친구 존재 확인
	if _, err := s.GetByID(friendID); err != nil {
		// %w 로 wrapping — 원본 에러 컨텍스트 추가
		return fmt.Errorf("친구 추가 실패: %w", err)
	}

	// 친구 관계 저장 (모델은 나중에)
	return nil
}
