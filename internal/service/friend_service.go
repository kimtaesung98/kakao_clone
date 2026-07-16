// internal/service/friend_service.go
package service

import (
	"kakao-clone/internal/model"
	"kakao-clone/internal/repository"
)

type FriendService struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
}

func NewFriendService(
	friendRepo repository.FriendRepository,
	userRepo repository.UserRepository,
) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

// 친구 추가 요청
func (s *FriendService) Add(userID, friendID uint) error {
	if userID == friendID {
		return ErrBadRequest("자기 자신을 친구 추가할 수 없습니다")
	}

	// 상대방 존재 확인
	if _, err := s.userRepo.FindByID(friendID); err != nil {
		return ErrNotFound("유저를 찾을 수 없습니다")
	}

	// 이미 친구인지 확인
	if s.friendRepo.IsFriend(userID, friendID) {
		return ErrBadRequest("이미 친구입니다")
	}

	return s.friendRepo.Add(userID, friendID)
}

// 친구 요청 수락
func (s *FriendService) Accept(userID, friendID uint) error {
	return s.friendRepo.Accept(userID, friendID)
}

// 친구 목록
func (s *FriendService) GetFriends(userID uint) ([]model.Friend, error) {
	friends, err := s.friendRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrInternal("친구 목록 조회 실패", err)
	}
	return friends, nil
}

// 친구 삭제
func (s *FriendService) Delete(userID, friendID uint) error {
	return s.friendRepo.Delete(userID, friendID)
}
