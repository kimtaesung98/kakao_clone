// main.go
package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"kakao-clone/internal/db"
	"kakao-clone/internal/handler"
	"kakao-clone/internal/middleware"
	"kakao-clone/internal/repository"
	"kakao-clone/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env 파일 없음")
	}

	db.Connect()

	// ─────────────────────────────────────────
	// 의존성 주입 — 아래에서 위로 생성
	// Repository → Service → Handler
	// ─────────────────────────────────────────
	userRepo := repository.NewUserRepository(db.DB)
	roomRepo := repository.NewRoomRepository(db.DB)
	messageRepo := repository.NewMessageRepository(db.DB)
	friendRepo := repository.NewFriendRepository(db.DB)

	userSvc := service.NewUserService(userRepo)
	roomSvc := service.NewRoomService(roomRepo, messageRepo)
	messageSvc := service.NewMessageService(messageRepo)
	friendSvc := service.NewFriendService(friendRepo, userRepo)

	authHandler := handler.NewAuthHandler(userSvc)
	userHandler := handler.NewUserHandler(userSvc)
	chatHandler := handler.NewChatHandler(roomSvc, messageSvc)
	friendHandler := handler.NewFriendHandler(friendSvc)

	// ─────────────────────────────────────────
	// Gin 설정
	// ─────────────────────────────────────────
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "kakao-clone"})
	})

	// 인증
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/google", authHandler.GoogleLogin)
	}

	// 인증 필요 라우트
	api := r.Group("/")
	api.Use(middleware.AuthRequired())
	{
		// 유저
		api.GET("/users/search", userHandler.Search)
		api.PATCH("/users/profile", userHandler.UpdateProfile)

		// 채팅
		api.GET("/chats", chatHandler.GetChatList)
		api.POST("/chats/direct", chatHandler.StartDirect)
		api.POST("/chats/group", chatHandler.CreateGroup)
		api.GET("/chats/:roomId/messages", chatHandler.GetMessages)
		api.PATCH("/chats/:roomId/read", chatHandler.MarkAsRead)

		// 친구
		api.GET("/friends", friendHandler.GetFriends)
		api.POST("/friends/:id", friendHandler.Add)
		api.PATCH("/friends/:id/accept", friendHandler.Accept)
		api.DELETE("/friends/:id", friendHandler.Delete)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("카카오 클론 서버 시작: http://localhost:%s", port)
	r.Run(":" + port)
}
