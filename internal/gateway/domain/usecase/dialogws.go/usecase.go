package dialogws_usecase

import (
	"context"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/service/client"
)

// type MessageService interface {
// 	Create(ctx context.Context, msg entity.Message) (entity.Message, error)
// 	GetByID(ctx context.Context, id int64) (entity.Message, error)
// 	GetByUsers(ctx context.Context, senderLogin, receiverLogin string) ([]entity.Message, error)
// 	UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) // id to string
// }

type DialogHub interface {
	RouteEvent(event entity.Event, senderClient *client.Client)
	AddClient(c *client.Client)
	RemoveClient(c *client.Client)
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

	newClient := &client.Client{
		Type:       client.Dialog, // probably useless
		Conn:       dto.Websocket,
		Message:    make(chan entity.Event, 5),
		SenderID:   dto.SenderID,
		ReceiverID: dto.ReceiverID,
		SessionID:  dto.SessionID,
	}

	u.dialogHub.AddClient(newClient)

	return nil

}
