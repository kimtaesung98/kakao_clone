// internal/service/room.go
package service

import (
	"fmt"
	"time"
)

// ─────────────────────────────────────────
// Message — 채팅방에서 주고받는 메시지
// ─────────────────────────────────────────
type Message struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	SenderID  uint      `json:"senderId"`
	RoomID    string    `json:"roomId"`
	CreatedAt time.Time `json:"createdAt"`
}

// ─────────────────────────────────────────
// Room interface — 채팅방의 약속
// 1:1이든 그룹이든 이 메서드를 구현하면 Room
// ─────────────────────────────────────────
type Room interface {
	GetID() string      // 방 ID 반환
	GetType() string    // "direct" or "group"
	Join(userID uint)   // 입장
	Leave(userID uint)  // 퇴장
	GetMembers() []uint // 멤버 목록
	Send(msg Message)   // 메시지 전송
}

// ─────────────────────────────────────────
// DirectRoom — 1:1 채팅방
// Room interface 구현
// ─────────────────────────────────────────
type DirectRoom struct {
	ID      string
	UserIDs [2]uint // 반드시 2명
}

func NewDirectRoom(userA, userB uint) *DirectRoom {
	// 방 ID = 두 유저 ID 조합 (작은 ID가 앞)
	// 항상 같은 방 ID가 나오도록
	if userA > userB {
		userA, userB = userB, userA
	}
	return &DirectRoom{
		ID:      fmt.Sprintf("direct_%d_%d", userA, userB),
		UserIDs: [2]uint{userA, userB},
	}
}

func (r *DirectRoom) GetID() string      { return r.ID }
func (r *DirectRoom) GetType() string    { return "direct" }
func (r *DirectRoom) GetMembers() []uint { return r.UserIDs[:] }
func (r *DirectRoom) Join(userID uint)   { /* 1:1은 입장 개념 없음 */ }
func (r *DirectRoom) Leave(userID uint)  { /* 1:1은 퇴장 개념 없음 */ }
func (r *DirectRoom) Send(msg Message) {
	// 실제 전송 로직은 RoomManager가 처리
}

// ─────────────────────────────────────────
// GroupRoom — 그룹 채팅방
// Room interface 구현
// ─────────────────────────────────────────
type GroupRoom struct {
	ID      string
	Name    string
	Members []uint // 여러 명
}

func NewGroupRoom(id, name string, creatorID uint) *GroupRoom {
	return &GroupRoom{
		ID:      id,
		Name:    name,
		Members: []uint{creatorID},
	}
}

func (r *GroupRoom) GetID() string      { return r.ID }
func (r *GroupRoom) GetType() string    { return "group" }
func (r *GroupRoom) GetMembers() []uint { return r.Members }

func (r *GroupRoom) Join(userID uint) {
	// 중복 체크
	for _, id := range r.Members {
		if id == userID {
			return
		}
	}
	r.Members = append(r.Members, userID)
}

func (r *GroupRoom) Leave(userID uint) {
	for i, id := range r.Members {
		if id == userID {
			// 슬라이스에서 제거
			r.Members = append(r.Members[:i], r.Members[i+1:]...)
			return
		}
	}
}

func (r *GroupRoom) Send(msg Message) {
	// 실제 전송 로직은 RoomManager가 처리
}

// ─────────────────────────────────────────
// RoomManager — 모든 채팅방 관리
// interface 덕분에 1:1과 그룹을 하나로 관리
// ─────────────────────────────────────────
type RoomManager struct {
	rooms map[string]Room // Room interface 타입으로 저장
	// → DirectRoom도, GroupRoom도 Room interface니까 같이 저장 가능
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]Room),
	}
}

func (m *RoomManager) AddRoom(room Room) {
	m.rooms[room.GetID()] = room
}

func (m *RoomManager) GetRoom(id string) (Room, bool) {
	room, ok := m.rooms[id]
	return room, ok
}

// GetOrCreateDirect — 1:1 방 가져오기 (없으면 생성)
func (m *RoomManager) GetOrCreateDirect(userA, userB uint) Room {
	room := NewDirectRoom(userA, userB)
	if existing, ok := m.rooms[room.ID]; ok {
		return existing
	}
	m.rooms[room.ID] = room
	return room
}
