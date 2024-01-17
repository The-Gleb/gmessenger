package service

import (
	"context"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
)

type UserStorage interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Get(ctx context.Context, login string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
}

type UserService struct {
	repo UserStorage
}

func NewUserService(s UserStorage) *UserService {
	return &UserService{repo: s}
}

func (us *UserService) Get(ctx context.Context, login string) (entity.User, error) {
	return us.repo.Get(ctx, login)
}

func (us *UserService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	return us.repo.Create(ctx, user)
}

func (us *UserService) GetAll(ctx context.Context) ([]entity.User, error) {
	return us.repo.GetAll(ctx)
}
