package service

import (
	"fmt"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/vk-rv/pvx"
	"time"
)

type PasetoAuth struct {
	pasetoKey    *pvx.SymKey
	symmetricKey []byte
}

const keySize = 32

var ErrInvalidSize = fmt.Errorf("bad key size: it must be %d bytes", keySize)

func NewPaseto(key []byte,
	tokenTTL time.Duration) (*PasetoAuth, error) {

	if len(key) != keySize {
		return nil, ErrInvalidSize
	}

	pasetoKey := pvx.NewSymmetricKey(key, pvx.Version4)

	return &PasetoAuth{
		symmetricKey: key,
		pasetoKey:    pasetoKey,
	}, nil
}

func (pa *PasetoAuth) NewToken(data entity.TokenData) (string, error) {

	serviceClaims := &entity.ServiceClaims{}

	iss := time.Now()

	serviceClaims.IssuedAt = &iss
	serviceClaims.Expiration = &data.Expiration
	serviceClaims.Subject = data.Subject

	serviceClaims.AdditionalClaims = data.AdditionalClaims
	serviceClaims.Footer = data.Footer

	pv4 := pvx.NewPV4Local()

	authToken, err := pv4.Encrypt(pa.pasetoKey, serviceClaims,
		pvx.WithFooter(serviceClaims.Footer))
	if err != nil {
		return "", err
	}

	return authToken, nil

}

func (pa *PasetoAuth) VerifyToken(token string) (*entity.ServiceClaims, error) {
	pv4 := pvx.NewPV4Local()
	tk := pv4.Decrypt(token, pa.pasetoKey)

	f := entity.Footer{}
	sc := entity.ServiceClaims{
		Footer: f,
	}

	err := tk.Scan(&sc, &f)
	if err != nil {
		return &sc, err
	}

	return &sc, nil
}
