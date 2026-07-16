// internal/db/db.go
package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"kakao-clone/internal/model"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}

	// 모든 모델 마이그레이션
	err = database.AutoMigrate(
		&model.User{},
		&model.Room{},
		&model.RoomMember{},
		&model.Message{},
		&model.Friend{},
	)
	if err != nil {
		log.Fatal("마이그레이션 실패:", err)
	}

	DB = database
	log.Println("PostgreSQL 연결 + 마이그레이션 완료")
}
