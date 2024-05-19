package register_usecase

import (
	"context"
	"fmt"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type sessionService interface {
	Create(ctx context.Context, userID int64) (entity.Session, error)
}

type pasetoAuthService interface {
	NewToken(data entity.TokenData) (string, error)
}

type userService interface {
	CreateWithPassword(ctx context.Context, dto entity.RegisterUserDTO) (int64, error)
}

type registerUsecase struct {
	userService    userService
	pasetoService  pasetoAuthService
	sessionService sessionService
}

func NewRegisterUsecase(us userService, ps pasetoAuthService, ss sessionService) *registerUsecase {
	return &registerUsecase{
		userService:    us,
		pasetoService:  ps,
		sessionService: ss,
	}
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

	session, err := r.sessionService.Create(ctx, id)
	if err != nil {
		return "", fmt.Errorf("[Login]: %w", err)
	}

	token, err := r.pasetoService.NewToken(entity.TokenData{
		Subject:    "sessionToken",
		Expiration: session.Expiry,
		AdditionalClaims: entity.AdditionalClaims{
			UserID:    id,
			SessionID: session.ID,
		},
		Footer: entity.Footer{MetaData: "footer"},
	})

	if err != nil {
		return "", fmt.Errorf("[usecase.Register]: %w", err)
	}

	return token, nil
}
