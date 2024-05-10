package login_usecase

import (
	"context"
	"fmt"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type sessionService interface {
	Create(ctx context.Context, userID int64) (entity.Session, error)
}

type pasetoAuthService interface {
	NewToken(data entity.TokenData) (string, error)
}

type userService interface {
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type loginUsecase struct {
	userService    userService
	sessionService sessionService
	pasetoService  pasetoAuthService
}

func NewLoginUsecase(us userService, ps pasetoAuthService, ss sessionService) *loginUsecase {
	return &loginUsecase{
		userService:    us,
		sessionService: ss,
		pasetoService:  ps,
	}
}

func (uc *loginUsecase) Login(ctx context.Context, dto entity.LoginDTO) (string, error) {

	user, err := uc.userService.GetByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Code(err) == errors.ErrNoDataFound {
			return "", errors.NewDomainError(errors.ErrUCWrongLoginOrPassword, "[usecase.Login]:")
		}
		return "", fmt.Errorf("[usecase.Login]: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return "", errors.NewDomainError(errors.ErrUCWrongLoginOrPassword, "[usecase.Login]:")
	}

	session, err := uc.sessionService.Create(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("[Login]: %w", err)
	}

	token, err := uc.pasetoService.NewToken(entity.TokenData{
		Subject:    "sessionToken",
		Expiration: session.Expiry,
		AdditionalClaims: entity.AdditionalClaims{
			UserID:    user.ID,
			SessionID: session.ID,
		},
		Footer: entity.Footer{MetaData: "footer"},
	})
	if err != nil {
		return "", fmt.Errorf("[Login]: %w", err)
	}

	return token, nil
}
