// internal/service/room_service.go
package service

import (
	"fmt"

	"kakao-clone/internal/model"
	"kakao-clone/internal/repository"
)

type RoomService struct {
	roomRepo    repository.RoomRepository
	messageRepo repository.MessageRepository
}

func NewRoomService(
	roomRepo repository.RoomRepository,
	messageRepo repository.MessageRepository,
) *RoomService {
	return &RoomService{
		roomRepo:    roomRepo,
		messageRepo: messageRepo,
	}
}

// ─────────────────────────────────────────
// 1:1 채팅방 가져오기 (없으면 생성)
// ─────────────────────────────────────────
func (s *RoomService) GetOrCreateDirect(userA, userB uint) (*model.Room, error) {
	// 기존 방 찾기
	room, err := s.roomRepo.FindDirectRoom(userA, userB)
	if err == nil {
		return room, nil
	}

	// 없으면 생성
	if userA > userB {
		userA, userB = userB, userA
	}
	room = &model.Room{
		ID:   fmt.Sprintf("direct_%d_%d", userA, userB),
		Type: "direct",
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, ErrInternal("채팅방 생성 실패", err)
	}

	// 멤버 추가
	s.roomRepo.AddMember(room.ID, userA)
	s.roomRepo.AddMember(room.ID, userB)

	return room, nil
}

// ─────────────────────────────────────────
// 그룹 채팅방 생성
// ─────────────────────────────────────────
func (s *RoomService) CreateGroup(name string, creatorID uint, memberIDs []uint) (*model.Room, error) {
	if name == "" {
		return nil, ErrBadRequest("그룹 이름을 입력하세요")
	}
	if len(memberIDs) < 2 {
		return nil, ErrBadRequest("그룹 채팅은 2명 이상이어야 합니다")
	}

	room := &model.Room{
		ID:   fmt.Sprintf("group_%d_%d", creatorID, len(memberIDs)),
		Name: name,
		Type: "group",
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, ErrInternal("그룹 생성 실패", err)
	}

	// 생성자 + 멤버 추가
	s.roomRepo.AddMember(room.ID, creatorID)
	for _, id := range memberIDs {
		s.roomRepo.AddMember(room.ID, id)
	}

	return room, nil
}

// ─────────────────────────────────────────
// 내 채팅 목록
// ─────────────────────────────────────────
type ChatListItem struct {
	Room        model.Room     `json:"room"`
	LastMessage *model.Message `json:"lastMessage"`
	UnreadCount int64          `json:"unreadCount"`
}

func (s *RoomService) GetChatList(userID uint) ([]ChatListItem, error) {
	rooms, err := s.roomRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrInternal("채팅 목록 조회 실패", err)
	}

	result := make([]ChatListItem, 0)
	for _, room := range rooms {
		// 마지막 메시지
		messages, _ := s.messageRepo.FindByRoomID(room.ID, 1)
		var lastMsg *model.Message
		if len(messages) > 0 {
			lastMsg = &messages[0]
		}

		// 안 읽은 수
		unread, _ := s.messageRepo.UnreadCount(room.ID, userID)

		result = append(result, ChatListItem{
			Room:        room,
			LastMessage: lastMsg,
			UnreadCount: unread,
		})
	}

	return result, nil
}
