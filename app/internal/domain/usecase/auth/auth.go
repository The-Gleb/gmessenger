package auth_usecase

import (
	"context"
	"fmt"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
)

type PasetoAuthService interface {
	VerifyToken(token string) (*entity.ServiceClaims, error)
}

type authUsecase struct {
	authService PasetoAuthService
}

func NewAuthUsecase(ss PasetoAuthService) *authUsecase {
	return &authUsecase{ss}
}

func (uc *authUsecase) Auth(ctx context.Context, token string) (entity.AdditionalClaims, error) {

	serviceClaims, err := uc.authService.VerifyToken(token)
	if err != nil {
		return entity.AdditionalClaims{}, fmt.Errorf("[Auth]: %w", err)
	}

	//if serviceClaims.Expiration.Before(time.Now()) {
	//
	//}

	return serviceClaims.AdditionalClaims, nil
}
