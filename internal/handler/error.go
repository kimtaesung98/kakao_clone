// internal/handler/error.go
package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"kakao-clone/internal/service"
)

// HandleError — AppError를 HTTP 응답으로 변환
// 모든 핸들러에서 공통으로 사용
func HandleError(c *gin.Context, err error) {
	// AppError인지 확인
	if appErr, ok := service.AsAppError(err); ok {
		// 500이면 서버 로그에 원본 에러 출력
		if appErr.Code >= 500 && appErr.Err != nil {
			log.Printf("서버 에러: %v", appErr.Err)
		}
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}

	// 알 수 없는 에러
	log.Printf("알 수 없는 에러: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
}
