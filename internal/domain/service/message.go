package service

import (
	"context"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
)

type MessageStorage interface {
	Create(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetByID(ctx context.Context, id int64) (entity.Message, error)
	GetByUsers(ctx context.Context, senderLogin, receiverLogin string) ([]entity.Message, error)
	UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) // id to string
}

type MessageService struct {
	repo MessageStorage
}

func NewMessageService(s MessageStorage) *MessageService {
	return &MessageService{repo: s}
}

func (ss *MessageService) GetByID(ctx context.Context, msgID int64) (entity.Message, error) {
	return ss.repo.GetByID(ctx, msgID)
}

func (ms *MessageService) Create(ctx context.Context, msg entity.Message) (entity.Message, error) {
	return ms.repo.Create(ctx, msg)
}

func (ms *MessageService) GetByUsers(ctx context.Context, sender, receiver string) ([]entity.Message, error) {
	return ms.repo.GetByUsers(ctx, sender, receiver)
}

func (ms *MessageService) UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) {
	return ms.repo.UpdateStatus(ctx, msgID, status)
}
