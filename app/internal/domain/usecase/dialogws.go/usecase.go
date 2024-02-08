package dialogws_usecase

import (
	"context"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service"
)

// type MessageService interface {
// 	Create(ctx context.Context, msg entity.Message) (entity.Message, error)
// 	GetByID(ctx context.Context, id int64) (entity.Message, error)
// 	GetByUsers(ctx context.Context, senderLogin, receiverLogin string) ([]entity.Message, error)
// 	UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) // id to string
// }

type DialogHub interface {
	RouteEvent(event entity.Event, senderClient *service.Client)
	AddClient(c *service.Client)
	RemoveClient(c *service.Client)
}

type dialogWSUsecase struct {
	// messageService MessageService
	dialogHub DialogHub
}

func NewDialogWSUsecase(dh DialogHub) *dialogWSUsecase {
	return &dialogWSUsecase{
		dialogHub: dh,
	}
}

func (u *dialogWSUsecase) OpenDialog(ctx context.Context, dto OpenDialogDTO) error {

	// TODO: check if receiver exists

	newClient := &service.Client{
		Type:          service.Dialog, // probably useless
		Conn:          dto.Websocket,
		Message:       make(chan entity.Event, 5),
		SenderLogin:   dto.SenderLogin,
		ReceiverLogin: dto.ReceiverLogin,
		SessionToken:  dto.SenderToken,
	}

	u.dialogHub.AddClient(newClient)

	return nil

}
