package register_usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type sessionService interface {
	Create(ctx context.Context, session entity.Session) error
}

type userService interface {
	CreateWithPassword(ctx context.Context, dto entity.RegisterUserDTO) (int64, error)
}

type registerUsecase struct {
	userService    userService
	sessionService sessionService
}

func NewRegisterUsecase(us userService, ss sessionService) *registerUsecase {
	return &registerUsecase{us, ss}
}

func (r *registerUsecase) Register(ctx context.Context, dto entity.RegisterUserDTO) (entity.Session, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.Session{}, fmt.Errorf("[usecase.Register]: %w", err)
	}
	dto.Password = string(hashedPassword)

	id, err := r.userService.CreateWithPassword(ctx, dto)
	if err != nil {
		return entity.Session{}, fmt.Errorf("[usecase.Register]: %w", err)
	}

	slog.Debug("user created", "ID", id)

	// TODO: JWT
	newSession := entity.Session{
		UserID: id,
		Token:  uuid.NewString(),
		Expiry: time.Now().Add(24 * time.Hour), // TODO config
	}

	err = r.sessionService.Create(ctx, newSession)

	if err != nil {
		return entity.Session{}, fmt.Errorf("[usecase.Register]: %w", err)
	}

	return newSession, nil
}
