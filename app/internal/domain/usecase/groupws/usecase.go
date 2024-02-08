package groupws_usecase

import (
	"context"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service"
)

type GroupHub interface {
	RouteEvent(event entity.Event, senderClient *service.Client)
	AddClient(c *service.Client)
	RemoveClient(c *service.Client)
}

type groupWSUsecase struct {
	groupHub GroupHub
}

func NewGroupWSUsecase(gh GroupHub) *groupWSUsecase {
	return &groupWSUsecase{gh}
}

func (u *groupWSUsecase) OpenGroup(ctx context.Context, dto OpenGroupDTO) error {

	// TODO: check if group exists and if user is it`s member

	newClient := &service.Client{
		Type:         service.Group, // probably useless
		Conn:         dto.Websocket,
		Message:      make(chan entity.Event, 5),
		SenderLogin:  dto.SenderLogin,
		SessionToken: dto.SenderToken,
		GroupID:      dto.GroupID,
	}

	u.groupHub.AddClient(newClient)

	return nil

}
