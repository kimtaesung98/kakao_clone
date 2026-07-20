// internal/handler/websocket.go
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"kakao-clone/internal/model"
	"kakao-clone/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

type Hub struct {
	rooms      map[string]map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *BroadcastMsg
	messageSvc *service.MessageService
}

type BroadcastMsg struct {
	RoomID  string
	Message model.Message
}

func NewHub(messageSvc *service.MessageService) *Hub {
	return &Hub{
		rooms:      make(map[string]map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMsg),
		messageSvc: messageSvc,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			log.Printf("유저 %d 연결됨", client.UserID)

		case client := <-h.Unregister:
			for roomID, room := range h.rooms {
				if _, ok := room[client.UserID]; ok {
					delete(room, client.UserID)
					close(client.Send)
					if len(room) == 0 {
						delete(h.rooms, roomID)
					}
				}
			}
			log.Printf("유저 %d 연결 끊김", client.UserID)

		case msg := <-h.Broadcast:
			if room, ok := h.rooms[msg.RoomID]; ok {
				data, _ := json.Marshal(msg.Message)
				for _, client := range room {
					select {
					case client.Send <- data:
					default:
						close(client.Send)
						delete(room, client.UserID)
					}
				}
			}
		}
	}
}

func (h *Hub) JoinRoom(roomID string, client *Client) {
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[uint]*Client)
	}
	h.rooms[roomID][client.UserID] = client
	log.Printf("유저 %d → 룸 %s 입장", client.UserID, roomID)
}

func (h *Hub) OnlineUsers(roomID string) []uint {
	room, ok := h.rooms[roomID]
	if !ok {
		return []uint{}
	}
	users := make([]uint, 0, len(room))
	for uid := range room {
		users = append(users, uid)
	}
	return users
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var input struct {
			Type    string `json:"type"`
			RoomID  string `json:"roomId"`
			Content string `json:"content"`
		}

		if err := json.Unmarshal(rawMsg, &input); err != nil {
			continue
		}

		switch input.Type {
		case "join":
			c.Hub.JoinRoom(input.RoomID, c)
			onlineData, _ := json.Marshal(gin.H{
				"type":        "online_users",
				"roomId":      input.RoomID,
				"onlineUsers": c.Hub.OnlineUsers(input.RoomID),
			})
			c.Send <- onlineData

		case "message":
			msg := model.Message{
				Content:   input.Content,
				Type:      "text",
				RoomID:    input.RoomID,
				SenderID:  c.UserID,
				CreatedAt: time.Now(),
			}
			if err := c.Hub.messageSvc.SaveMessage(&msg); err != nil {
				log.Println("메시지 저장 실패:", err)
				continue
			}
			c.Hub.Broadcast <- &BroadcastMsg{
				RoomID:  input.RoomID,
				Message: msg,
			}
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		msg, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

type WebSocketHandler struct {
	hub *Hub
}

func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		auth := c.GetHeader("Authorization")
		tokenStr = strings.TrimPrefix(auth, "Bearer ")
	}

	userID, err := validateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 토큰"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket 업그레이드 실패:", err)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func validateToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return uint(claims["userId"].(float64)), nil
}
