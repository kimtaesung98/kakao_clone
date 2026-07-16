// internal/handler/user.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kakao-clone/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GET /users/search?q=키워드
func (h *UserHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "검색어를 입력하세요"})
		return
	}

	myID := c.MustGet("userId").(uint)
	users, err := h.userService.Search(keyword, myID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// PATCH /users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var input struct {
		Name      string `json:"name"`
		StatusMsg string `json:"statusMsg"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	myID := c.MustGet("userId").(uint)
	user, err := h.userService.UpdateProfile(myID, input.Name, input.StatusMsg)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
