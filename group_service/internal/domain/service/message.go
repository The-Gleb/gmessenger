package service

import (
	"context"

	"github.com/The-Gleb/gmessenger/group_service/internal/domain/entity"
)

type MessageStorage interface {
	AddMessage(ctx context.Context, message entity.Message) (entity.Message, error)
	GetMessages(ctx context.Context, groupID int64, limit, offset int) ([]entity.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID int64, status string) (entity.Message, error)
	GetLastMessage(ctx context.Context, groupID int64) (entity.Message, error)
}

type messageService struct {
	storage MessageStorage
}

func NewMessageService(s MessageStorage) *messageService {
	return &messageService{storage: s}
}

func (ms *messageService) AddMessage(ctx context.Context, message entity.Message) (entity.Message, error) {
	return ms.storage.AddMessage(ctx, message)
}
func (ms *messageService) GetMessages(ctx context.Context, groupID int64, limit, offset int) ([]entity.Message, error) {
	return ms.storage.GetMessages(ctx, groupID, limit, offset)
}
func (ms *messageService) UpdateMessageStatus(ctx context.Context, messageID int64, status string) (entity.Message, error) {
	return ms.storage.UpdateMessageStatus(ctx, messageID, status)
}
func (ms *messageService) GetLastMessage(ctx context.Context, groupID int64) (entity.Message, error) {
	return ms.storage.GetLastMessage(ctx, groupID)
}
