package login_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type sessionService interface {
	Create(ctx context.Context, session entity.Session) error
}

type userService interface {
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type loginUsecase struct {
	userService    userService
	sessionService sessionService
}

func NewLoginUsecase(us userService, ss sessionService) *loginUsecase {
	return &loginUsecase{us, ss}
}

func (uc *loginUsecase) Login(ctx context.Context, dto entity.LoginDTO) (entity.Session, error) {

	user, err := uc.userService.GetByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Code(err) == errors.ErrNoDataFound {
			return entity.Session{}, errors.NewDomainError(errors.ErrUCWrongLoginOrPassword, "[usecase.Login]:")
		}
		return entity.Session{}, fmt.Errorf("[usecase.Login]: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return entity.Session{}, errors.NewDomainError(errors.ErrUCWrongLoginOrPassword, "[usecase.Login]:")
	}

	// TODO: make JWT
	newSession := entity.Session{
		UserID: user.ID,
		Token:  uuid.NewString(),
		Expiry: time.Now().Add(24 * time.Hour), // TODO config
	}

	err = uc.sessionService.Create(ctx, newSession)

	if err != nil {
		return entity.Session{}, fmt.Errorf("[Login]: %w", err)
	}

	return newSession, nil
}
