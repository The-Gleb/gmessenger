package dialogws_usecase

import (
	"context"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/domain/service"
)

type MessageService interface {
	Create(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetByID(ctx context.Context, id int64) (entity.Message, error)
	GetByUsers(ctx context.Context, senderLogin, receiverLogin string) ([]entity.Message, error)
	UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) // id to string
}

type DialogService interface {
	AddClient(c *service.Client)
}

type dialogWSUsecase struct {
	messageService MessageService
	dialogService  DialogService
}

func NewDialogWSUsecase(ms MessageService, ds DialogService) *dialogWSUsecase {
	return &dialogWSUsecase{
		messageService: ms,
		dialogService:  ds,
	}
}

func (u *dialogWSUsecase) OpenDialog(ctx context.Context, dto OpenDialogDTO) error {

	newClient := &service.Client{
		Conn:          dto.Websocket,
		Message:       make(chan entity.Event, 5),
		SenderLogin:   dto.SenderLogin,
		ReceiverLogin: dto.ReceiverLogin,
		SessionToken:  dto.SenderToken,
	}

	u.dialogService.AddClient(newClient)

	return nil

}
