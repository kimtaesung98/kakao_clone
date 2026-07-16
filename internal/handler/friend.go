// internal/handler/friend.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"kakao-clone/internal/service"
)

type FriendHandler struct {
	friendService *service.FriendService
}

func NewFriendHandler(friendService *service.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

// GET /friends — 친구 목록
func (h *FriendHandler) GetFriends(c *gin.Context) {
	myID := c.MustGet("userId").(uint)
	friends, err := h.friendService.GetFriends(myID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

// POST /friends/:id — 친구 추가
func (h *FriendHandler) Add(c *gin.Context) {
	friendID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
		return
	}

	myID := c.MustGet("userId").(uint)
	if err := h.friendService.Add(myID, uint(friendID)); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "친구 추가 완료"})
}

// PATCH /friends/:id/accept — 친구 수락
func (h *FriendHandler) Accept(c *gin.Context) {
	friendID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
		return
	}

	myID := c.MustGet("userId").(uint)
	if err := h.friendService.Accept(myID, uint(friendID)); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "친구 수락 완료"})
}

// DELETE /friends/:id — 친구 삭제
func (h *FriendHandler) Delete(c *gin.Context) {
	friendID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
		return
	}

	myID := c.MustGet("userId").(uint)
	if err := h.friendService.Delete(myID, uint(friendID)); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "친구 삭제 완료"})
}
