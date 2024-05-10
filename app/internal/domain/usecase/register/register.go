package register_usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type pasetoAuthService interface {
	NewToken(data entity.TokenData) (string, error)
}

type userService interface {
	CreateWithPassword(ctx context.Context, dto entity.RegisterUserDTO) (int64, error)
}

type registerUsecase struct {
	userService   userService
	pasetoService pasetoAuthService
}

func NewRegisterUsecase(us userService, ss pasetoAuthService) *registerUsecase {
	return &registerUsecase{us, ss}
}

func (r *registerUsecase) Register(ctx context.Context, dto entity.RegisterUserDTO) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("[usecase.Register]: %w", err)
	}
	dto.Password = string(hashedPassword)

	id, err := r.userService.CreateWithPassword(ctx, dto)
	if err != nil {
		return "", fmt.Errorf("[usecase.Register]: %w", err)
	}

	slog.Debug("user created", "ID", id)

	token, err := r.pasetoService.NewToken(entity.TokenData{
		Subject:  "sessionToken",
		Duration: 5 * time.Second,
		AdditionalClaims: entity.AdditionalClaims{
			UserID: id,
		},
		Footer: entity.Footer{MetaData: "footer"},
	})

	if err != nil {
		return "", fmt.Errorf("[usecase.Register]: %w", err)
	}

	return token, nil
}
