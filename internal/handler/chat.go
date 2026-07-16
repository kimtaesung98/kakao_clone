// internal/handler/chat.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"kakao-clone/internal/service"
)

type ChatHandler struct {
	roomService    *service.RoomService
	messageService *service.MessageService
}

func NewChatHandler(
	roomService *service.RoomService,
	messageService *service.MessageService,
) *ChatHandler {
	return &ChatHandler{
		roomService:    roomService,
		messageService: messageService,
	}
}

// GET /chats — 내 채팅 목록
func (h *ChatHandler) GetChatList(c *gin.Context) {
	myID := c.MustGet("userId").(uint)
	list, err := h.roomService.GetChatList(myID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

// POST /chats/direct — 1:1 채팅방 시작
func (h *ChatHandler) StartDirect(c *gin.Context) {
	var input struct {
		TargetID uint `json:"targetId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	myID := c.MustGet("userId").(uint)
	room, err := h.roomService.GetOrCreateDirect(myID, input.TargetID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, room)
}

// POST /chats/group — 그룹 채팅방 생성
func (h *ChatHandler) CreateGroup(c *gin.Context) {
	var input struct {
		Name      string `json:"name"      binding:"required"`
		MemberIDs []uint `json:"memberIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	myID := c.MustGet("userId").(uint)
	room, err := h.roomService.CreateGroup(input.Name, myID, input.MemberIDs)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, room)
}

// GET /chats/:roomId/messages — 메시지 조회
func (h *ChatHandler) GetMessages(c *gin.Context) {
	roomID := c.Param("roomId")
	limit := 50
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil {
			limit = n
		}
	}

	messages, err := h.messageService.GetMessages(roomID, limit)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, messages)
}

// PATCH /chats/:roomId/read — 읽음 처리
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	roomID := c.Param("roomId")
	myID := c.MustGet("userId").(uint)

	if err := h.messageService.MarkAsRead(roomID, myID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "읽음 처리 완료"})
}
