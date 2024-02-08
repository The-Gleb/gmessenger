package service

import (
	"context"
	"fmt"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/google/uuid"
)

type SessionStorage interface {
	Create(ctx context.Context, session entity.Session) error
	GetByToken(ctx context.Context, token string) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type SessionService struct {
	repo SessionStorage
}

func NewSessionService(s SessionStorage) *SessionService {
	return &SessionService{repo: s}
}

func (ss *SessionService) GetLoginByToken(ctx context.Context, token string) (entity.Session, error) {

	session, err := ss.repo.GetByToken(ctx, token)
	if err != nil {
		return entity.Session{}, err
	}
	if session.IsExpired() {
		return session, errors.NewDomainError(errors.ErrSessionExpired, "[service.GetLoginByToken]:")
	}
	return session, nil
}

func (ss *SessionService) Create(ctx context.Context, session entity.Session) error {

	for {
		err := ss.repo.Create(ctx, session)
		if errors.Code(err) == errors.ErrNotUniqueToken { //TODO
			session.Token = uuid.NewString()
			continue
		}
		if err != nil {
			return fmt.Errorf("[Create]: %w", err)
		}
		break
	}
	return nil
}

func (ss *SessionService) Delete(ctx context.Context, token string) error {
	err := ss.repo.Delete(ctx, token)
	if err != nil {
		return fmt.Errorf("[Delete]: %w", err)
	}
	return nil
}
