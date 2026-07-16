// internal/service/message_service.go
package service

import (
	"kakao-clone/internal/model"
	"kakao-clone/internal/repository"
)

type MessageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{messageRepo: repo}
}

func (s *MessageService) GetMessages(roomID string, limit int) ([]model.Message, error) {
	messages, err := s.messageRepo.FindByRoomID(roomID, limit)
	if err != nil {
		return nil, ErrInternal("메시지 조회 실패", err)
	}
	return messages, nil
}

func (s *MessageService) MarkAsRead(roomID string, userID uint) error {
	return s.messageRepo.MarkAsRead(roomID, userID)
}

func (s *MessageService) SaveMessage(msg *model.Message) error {
	return s.messageRepo.Create(msg)
}
