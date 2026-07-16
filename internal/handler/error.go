// internal/handler/error.go
package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"kakao-clone/internal/service"
)

func handleError(c *gin.Context, err error) {
	if appErr, ok := service.AsAppError(err); ok {
		if appErr.Code >= 500 && appErr.Err != nil {
			log.Printf("서버 에러: %v", appErr.Err)
		}
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}
	log.Printf("알 수 없는 에러: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
}
