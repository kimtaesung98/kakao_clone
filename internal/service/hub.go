// internal/service/hub.go
package service

import "sync"

// ─────────────────────────────────────────
// Client — 연결된 유저 한 명
// ─────────────────────────────────────────
type Client struct {
	UserID uint
	RoomID string
	Send   chan Message // 이 유저에게 보낼 메시지 채널
}

// ─────────────────────────────────────────
// Hub — 전체 채팅 허브
// 모든 메시지 흐름을 channel로 관리
// ─────────────────────────────────────────
type Hub struct {
	// 룸별 클라이언트 목록
	rooms map[string]map[uint]*Client

	// channel들 — 허브로 들어오는 이벤트
	Register   chan *Client // 새 연결
	Unregister chan *Client // 연결 종료
	Broadcast  chan Message // 메시지 전송

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

// ─────────────────────────────────────────
// Run — 허브 실행
// goroutine으로 실행 → 블로킹 없이 이벤트 처리
// ─────────────────────────────────────────
func (h *Hub) Run() {
	for {
		select { // ← channel 여러 개를 동시에 대기
		case client := <-h.Register:
			h.addClient(client)

		case client := <-h.Unregister:
			h.removeClient(client)

		case msg := <-h.Broadcast:
			h.broadcastToRoom(msg)
		}
	}
}

func (h *Hub) addClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[c.RoomID] == nil {
		h.rooms[c.RoomID] = make(map[uint]*Client)
	}
	h.rooms[c.RoomID][c.UserID] = c
}

func (h *Hub) removeClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[c.RoomID]; ok {
		delete(room, c.UserID)
		close(c.Send) // 채널 닫기
		if len(room) == 0 {
			delete(h.rooms, c.RoomID)
		}
	}
}

func (h *Hub) broadcastToRoom(msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[msg.RoomID]
	if !ok {
		return
	}

	// 룸의 모든 멤버에게 전송
	for _, client := range room {
		select {
		case client.Send <- msg: // 메시지 전송
		default:
			// 채널이 가득 찼으면 스킵 (느린 클라이언트 보호)
		}
	}
}

// 온라인 유저 목록
func (h *Hub) OnlineUsers(roomID string) []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var users []uint
	for uid := range h.rooms[roomID] {
		users = append(users, uid)
	}
	return users
}
