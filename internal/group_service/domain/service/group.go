package service

import (
	"context"

	"github.com/The-Gleb/gmessenger/group_service/internal/domain/entity"
)

type GroupStorage interface {
	Create(ctx context.Context, group entity.GroupCreate) (entity.Group, error)
	IsMember(ctx context.Context, userLogin string, groupID int64) (bool, error)
	GetMembers(ctx context.Context, groupID int64) ([]string, error)
	GetGroups(ctx context.Context, userLogin string, limit, offset int) ([]entity.Group, error)
	Exists(ctx context.Context, groupID int64) (bool, error)
}

type groupService struct {
	storage GroupStorage
}

func NewGroupService(gs GroupStorage) *groupService {
	return &groupService{gs}
}

func (gs *groupService) Create(ctx context.Context, group entity.GroupCreate) (entity.Group, error) {
	return gs.storage.Create(ctx, group)
}
func (gs *groupService) IsMember(ctx context.Context, userLogin string, groupID int64) (bool, error) {
	return gs.storage.IsMember(ctx, userLogin, groupID)
}
func (gs *groupService) GetMembers(ctx context.Context, groupID int64) ([]string, error) {
	return gs.storage.GetMembers(ctx, groupID)
}
func (gs *groupService) GetGroups(ctx context.Context, userLogin string, limit, offset int) ([]entity.Group, error) {
	return gs.storage.GetGroups(ctx, userLogin, limit, offset)
}
func (gs *groupService) Exists(ctx context.Context, groupID int64) (bool, error) {
	return gs.storage.Exists(ctx, groupID)
}
