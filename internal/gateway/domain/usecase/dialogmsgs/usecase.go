package dialogmsgs_usecase

import (
	"context"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/service/client"
	"log/slog"
)

type MessageService interface {
	Create(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetByID(ctx context.Context, id int64) (entity.Message, error)
	GetByUsers(ctx context.Context, senderID, receiverID int64, limit, offset int) ([]entity.Message, error)
	UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) // id to string
}

type DialogService interface {
	AddClient(c *client.Client)
}

type dialogMsgsUsecase struct {
	messageService MessageService
}

func NewDialogMsgsUsecase(ms MessageService) *dialogMsgsUsecase {
	return &dialogMsgsUsecase{
		messageService: ms,
	}
}

func (u *dialogMsgsUsecase) GetDialogMessages(ctx context.Context, dto GetDialogMessagesDTO) ([]entity.Message, error) {

	// TODO: check if receiver exists

	messages, err := u.messageService.GetByUsers(ctx, dto.SenderID, dto.ReceiverID, 100, 0)
	slog.Debug("messages by user", "messages", messages)
	if err != nil {
		slog.Error(err.Error()) // TODO:
		return []entity.Message{}, err
	}

	return messages, nil

}
