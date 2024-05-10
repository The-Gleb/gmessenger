package oauth_usecase

import (
	"context"
	"fmt"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"time"
)

type PasetoAuthService interface {
	NewToken(data entity.TokenData) (string, error)
}

type userService interface {
	GetOrCreateByEmail(ctx context.Context, email string) (int64, bool, error)
}

type oauthUsecase struct {
	userService   userService
	pasetoService PasetoAuthService
}

func NewOAuthUsecase(us userService, ss PasetoAuthService) *oauthUsecase {
	return &oauthUsecase{us, ss}
}

func (uc *oauthUsecase) OAuth(ctx context.Context, email string) (string, bool, error) {

	id, isNew, err := uc.userService.GetOrCreateByEmail(ctx, email)
	if err != nil {
		return "", false, err
	}

	token, err := uc.pasetoService.NewToken(entity.TokenData{
		Subject:  "sessionToken",
		Duration: 5 * time.Second,
		AdditionalClaims: entity.AdditionalClaims{
			UserID: id,
		},
		Footer: entity.Footer{MetaData: "footer"},
	})

	if err != nil {
		return "", false, fmt.Errorf("[usecase.OAuth]: %w", err)
	}

	return token, isNew, nil
}
