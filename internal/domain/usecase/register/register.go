package register_usecase

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
	GetAllUsernames(ctx context.Context) ([]string, error)
}

type registerUsecase struct {
	userService    UserService
	sessionService SessionService
}

func NewRegisterUsecase(us UserService, ss SessionService) *registerUsecase {
	return &registerUsecase{us, ss}
}

func (r *registerUsecase) Register(ctx context.Context, registerUserDTO RegisterUserDTO) (entity.Session, error) {

	// err := json.NewDecoder(body).Decode(&createUserDTO)
	// if err != nil {
	// 	return "", time.Now(), errors.New("Register: couldn`t umarshall json")
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUserDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.Session{}, fmt.Errorf("[Register]: %w", err)
	}
	registerUserDTO.Password = string(hashedPassword)

	user, err := r.userService.Create(ctx, entity.User{
		UserName: registerUserDTO.UserName,
		Login:    registerUserDTO.Login,
		Password: registerUserDTO.Password,
	})
	if err != nil {
		return entity.Session{}, fmt.Errorf("[Register]: %w", err)
	}

	newSession := entity.Session{
		UserLogin: user.Login,
		Token:     uuid.NewString(),
		Expiry:    time.Now().Add(24 * time.Hour), // TODO config
	}

	err = r.sessionService.Create(ctx, newSession)

	if err != nil {
		return entity.Session{}, fmt.Errorf("[Register]: %w", err)
	}

	return newSession, nil
}
