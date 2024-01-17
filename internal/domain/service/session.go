package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	"github.com/google/uuid"
)

type SessionStorage interface {
	Create(ctx context.Context, session entity.Session) error
	GetLoginByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type SessionService struct {
	repo SessionStorage
}

func NewSessionService(s SessionStorage) *SessionService {
	return &SessionService{repo: s}
}

func (ss *SessionService) GetLoginByToken(ctx context.Context, token string) (entity.Session, error) {

	session, err := ss.repo.GetLoginByToken(ctx, token)
	if err != nil {
		return entity.Session{}, err
	}
	if session.IsExpired() {
		return session, errors.New("session token is expired") // TODO
	}
	return session, nil
}

func (ss *SessionService) Create(ctx context.Context, session entity.Session) error {

	for {
		err := ss.repo.Create(ctx, session)
		if errors.Is(err, errors.New("")) { //TODO
			session.Token = uuid.NewString()
			continue
		}
		if err != nil {
			return fmt.Errorf("[Register]: %w", err)
		}
		break
	}
	return nil
}
