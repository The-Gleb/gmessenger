package service

import (
	"context"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
)

type UserStorage interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByLogin(ctx context.Context, login string) (entity.User, error)
	GetPassword(ctx context.Context, login string) (string, error)
	GetAllUsernames(ctx context.Context) ([]string, error)
}

type UserService struct {
	repo UserStorage
}

func NewUserService(s UserStorage) *UserService {
	return &UserService{repo: s}
}

func (us *UserService) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	return us.repo.GetByLogin(ctx, login)
}

func (us *UserService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	return us.repo.Create(ctx, user)
}

func (us *UserService) GetPassword(ctx context.Context, login string) (string, error) {
	return us.repo.GetPassword(ctx, login)
}

func (us *UserService) GetAllUsernames(ctx context.Context) ([]string, error) {
	return us.repo.GetAllUsernames(ctx)
}
