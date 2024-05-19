package username_usecase

import (
	"context"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
)

type userService interface {
	SetUsername(ctx context.Context, dto entity.SetUsernameDTO) error
}

type usernameUsecase struct {
	userService userService
}

func NewUsernameUsecase(us userService) *usernameUsecase {
	return &usernameUsecase{us}
}

func (r *usernameUsecase) SetUsername(ctx context.Context, dto entity.SetUsernameDTO) error {

	return r.userService.SetUsername(ctx, dto)

}
