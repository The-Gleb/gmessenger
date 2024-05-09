package groupws_usecase

import (
	"context"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
)

type GroupHub interface {
	RouteEvent(event entity.Event, senderClient *client.Client)
	AddClient(c *client.Client)
	RemoveClient(c *client.Client)
}

type groupWSUsecase struct {
	groupHub GroupHub
}

func NewGroupWSUsecase(gh GroupHub) *groupWSUsecase {
	return &groupWSUsecase{gh}
}

func (u *groupWSUsecase) OpenGroup(ctx context.Context, dto OpenGroupDTO) error {

	// TODO: check if group exists and if user is it`s member

	newClient := &client.Client{
		Type:         client.Group, // probably useless
		Conn:         dto.Websocket,
		Message:      make(chan entity.Event, 5),
		SenderID:     dto.SenderID,
		SessionToken: dto.SenderToken,
		GroupID:      dto.GroupID,
	}

	u.groupHub.AddClient(newClient)

	return nil

}
