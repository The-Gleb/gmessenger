package create_group

import (
	"context"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/pkg/proto/group"
	"log/slog"
)

type createGroupUsecase struct {
	groupClient group.GroupClient
}

func NewCreateGroupUsecase(c group.GroupClient) *createGroupUsecase {
	return &createGroupUsecase{c}
}

func (uc createGroupUsecase) CreateGroup(ctx context.Context, dto entity.CreateGroupDTO) (entity.Group, error) {
	response, err := uc.groupClient.Create(ctx, &group.CreateRequest{
		Name:      dto.Name,
		MemberIds: dto.MemberIDs,
	})

	if err != nil {
		slog.Error(err.Error())
		return entity.Group{}, err
	}

	return entity.Group{
		ID:        response.Id,
		Name:      dto.Name,
		MemberIDs: dto.MemberIDs,
	}, err

}
