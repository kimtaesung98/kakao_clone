// internal/handler/auth.go
package handler

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"

	"kakao-clone/internal/service"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// ─────────────────────────────────────────
// POST /auth/register
// ─────────────────────────────────────────
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name"     binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(input.Email, input.Password, input.Name)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "회원가입 성공",
		"user":    user,
	})
}

// ─────────────────────────────────────────
// POST /auth/login
// ─────────────────────────────────────────
func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"    binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.userService.Login(input.Email, input.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// ─────────────────────────────────────────
// POST /auth/google
// ─────────────────────────────────────────
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var input struct {
		IDToken string `json:"idToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload, err := idtoken.Validate(
		context.Background(),
		input.IDToken,
		os.Getenv("GOOGLE_CLIENT_ID"),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 Google 토큰"})
		return
	}

	email := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	googleID := payload.Subject

	user, token, err := h.userService.GoogleLogin(email, name, googleID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
