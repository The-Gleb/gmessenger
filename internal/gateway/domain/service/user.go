package service

import (
	"context"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
)

type UserStorage interface {
	GetOrCreateByEmail(ctx context.Context, email string) (int64, bool, error)
	SetUsername(ctx context.Context, dto entity.SetUsernameDTO) error
	CreateWithPassword(ctx context.Context, dto entity.RegisterUserDTO) (int64, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetAllUsersView(ctx context.Context) ([]entity.UserView, error)
	GetUserInfoByID(ctx context.Context, id int64) (entity.UserInfo, error)

	GetChatsView(ctx context.Context, userID int64) ([]entity.Chat, error)
}

type UserService struct {
	repo UserStorage
}

func NewUserService(s UserStorage) *UserService {
	return &UserService{repo: s}
}

func (us UserService) GetUserInfoByID(ctx context.Context, ID int64) (entity.UserInfo, error) {
	return us.repo.GetUserInfoByID(ctx, ID)
}

func (us *UserService) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	return us.repo.GetByEmail(ctx, email)
}

func (us *UserService) CreateWithPassword(ctx context.Context, dto entity.RegisterUserDTO) (int64, error) {
	return us.repo.CreateWithPassword(ctx, dto)
}

func (us *UserService) GetOrCreateByEmail(ctx context.Context, email string) (int64, bool, error) {
	return us.repo.GetOrCreateByEmail(ctx, email)
}

func (us *UserService) SetUsername(ctx context.Context, dto entity.SetUsernameDTO) error {
	return us.repo.SetUsername(ctx, dto)
}

func (us *UserService) GetAllUsersView(ctx context.Context) ([]entity.UserView, error) {
	return us.repo.GetAllUsersView(ctx)
}

func (us *UserService) GetChatsView(ctx context.Context, userID int64) ([]entity.Chat, error) {
	return us.repo.GetChatsView(ctx, userID)
}
