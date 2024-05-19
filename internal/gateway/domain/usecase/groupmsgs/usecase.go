package groupmsgs_usecase

import (
	"context"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/gateway/errors"
	"github.com/The-Gleb/gmessenger/pkg/proto/group"
	"log/slog"
)

type groupMsgsUsecase struct {
	groupClient group.GroupClient
}

func NewGroupMsgsUsecase(gc group.GroupClient) *groupMsgsUsecase {
	return &groupMsgsUsecase{
		groupClient: gc,
	}
}

func (u *groupMsgsUsecase) GetGroupMessages(ctx context.Context, dto GetGroupMessagesDTO) ([]entity.Message, error) {

	isMemberResp, err := u.groupClient.CheckMember(ctx, &group.CheckMemberRequest{
		UserId:  dto.UserId,
		GroupId: dto.GroupID,
	})

	if err != nil {
		slog.Error(err.Error())
		return []entity.Message{}, err // TODO:
	}

	if !isMemberResp.GetIsMember() {
		slog.Error("client is not a member of this chat", "userLogin", dto.UserId, "group ID", dto.GroupID)
		return []entity.Message{}, errors.NewDomainError(errors.ErrNotAMember, "")
	}

	getMessagesResp, err := u.groupClient.GetMessages(ctx, &group.GetMessagesRequest{
		GroupId: dto.GroupID,
	})

	if err != nil {
		slog.Error(err.Error())
		return []entity.Message{}, err // TODO
	}

	grpcMessages := getMessagesResp.GetMessages()

	messages := make([]entity.Message, len(grpcMessages))

	for i, grpcMessage := range grpcMessages {
		messages[i] = entity.Message{
			ID: grpcMessage.GetId(),
			//SenderID:    grpcMessage.GetSenderID(), // TODO: regenerate proto
			Text:      grpcMessage.GetText(),
			Status:    grpcMessage.GetStatus().String(),
			Timestamp: grpcMessage.Timestamp.AsTime(),
		}
	}

	return messages, err

}
