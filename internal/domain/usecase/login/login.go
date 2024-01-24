package login_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SessionService interface {
	Create(ctx context.Context, session entity.Session) error
	GetLoginByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type UserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByLogin(ctx context.Context, login string) (entity.User, error)
	GetPassword(ctx context.Context, login string) (string, error)
}

type loginUsecase struct {
	userService    UserService
	sessionService SessionService
}

func NewLoginUsecase(us UserService, ss SessionService) *loginUsecase {
	return &loginUsecase{us, ss}
}

func (uc *loginUsecase) Login(ctx context.Context, loginDTO LoginDTO) (entity.Session, error) {

	password, err := uc.userService.GetPassword(ctx, loginDTO.Login)
	if err != nil {
		return entity.Session{}, fmt.Errorf("[usecase.Login]: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(loginDTO.Password))
	if err != nil {
		return entity.Session{}, err // TODO
	}

	newSession := entity.Session{
		UserLogin: loginDTO.Login,
		Token:     uuid.NewString(),
		Expiry:    time.Now().Add(24 * time.Hour), // TODO config
	}

	err = uc.sessionService.Create(ctx, newSession)

	if err != nil {
		return entity.Session{}, fmt.Errorf("[Login]: %w", err)
	}

	return newSession, nil
}
