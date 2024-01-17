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
	GetLoginByToken(ctx context.Context, token string) (entity.Session, error)
	Create(ctx context.Context, Session entity.Session) error
}

type UserService interface {
	Get(ctx context.Context, login string) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
}

type loginUsecase struct {
	userService    UserService
	sessionService SessionService
}

func NewLoginUsecase(us UserService, ss SessionService) *loginUsecase {
	return &loginUsecase{us, ss}
}

func (uc *loginUsecase) Register(ctx context.Context, loginDTO LoginDTO) (entity.Session, error) {

	user, err := uc.userService.Get(ctx, loginDTO.Login)
	if err != nil {
		return entity.Session{}, fmt.Errorf("[Register]: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Password))
	if err != nil {
		return entity.Session{}, err // TODO
	}

	newSession := entity.Session{
		UserName: user.Login,
		Token:    uuid.NewString(),
		Expiry:   time.Now().Add(24 * time.Hour), // TODO config
	}

	err = uc.sessionService.Create(ctx, newSession)

	if err != nil {
		return entity.Session{}, fmt.Errorf("[Login]: %w", err)
	}

	return newSession, nil
}
