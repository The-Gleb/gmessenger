package oauth_usecase

import (
	"context"
	"fmt"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
)

type sessionService interface {
	Create(ctx context.Context, userID int64) (entity.Session, error)
}

type PasetoAuthService interface {
	NewToken(data entity.TokenData) (string, error)
}

type userService interface {
	GetOrCreateByEmail(ctx context.Context, email string) (int64, bool, error)
}

type oauthUsecase struct {
	userService    userService
	pasetoService  PasetoAuthService
	sessionService sessionService
}

func NewOAuthUsecase(us userService, ps PasetoAuthService, ss sessionService) *oauthUsecase {
	return &oauthUsecase{
		userService:    us,
		pasetoService:  ps,
		sessionService: ss,
	}
}

func (uc *oauthUsecase) OAuth(ctx context.Context, email string) (string, bool, error) {

	id, isNew, err := uc.userService.GetOrCreateByEmail(ctx, email)
	if err != nil {
		return "", false, err
	}

	session, err := uc.sessionService.Create(ctx, id)
	if err != nil {
		return "", false, fmt.Errorf("[Login]: %w", err)
	}

	token, err := uc.pasetoService.NewToken(entity.TokenData{
		Subject:    "sessionToken",
		Expiration: session.Expiry,
		AdditionalClaims: entity.AdditionalClaims{
			UserID:    id,
			SessionID: session.ID,
		},
		Footer: entity.Footer{MetaData: "footer"},
	})

	if err != nil {
		return "", false, fmt.Errorf("[usecase.OAuth]: %w", err)
	}

	return token, isNew, nil
}
