package service

import (
	"context"
	"fmt"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
)

type SessionStorage interface {
	Create(ctx context.Context, session entity.Session) (entity.Session, error)
	Delete(ctx context.Context, token string) error
}

type SessionService struct {
	repo SessionStorage
	ttl  time.Duration
}

func NewSessionService(s SessionStorage) *SessionService {
	return &SessionService{repo: s}
}

//func (ss *SessionService) GetByToken(ctx context.Context, token string) (entity.Session, error) {
//
//	session, err := ss.repo.GetByToken(ctx, token)
//	if err != nil {
//		return entity.Session{}, err
//	}
//	if session.IsExpired() {
//		return session, errors.NewDomainError(errors.ErrSessionExpired, "[service.GetByToken]:")
//	}
//	return session, nil
//}

func (ss *SessionService) Create(ctx context.Context, userID int64) (entity.Session, error) {

	session, err := ss.repo.Create(ctx, entity.Session{
		UserID: userID,
		Expiry: time.Now().Add(ss.ttl),
	})
	if err != nil {
		return session, fmt.Errorf("[Create]: %w", err)
	}

	return session, nil
}

func (ss *SessionService) Delete(ctx context.Context, token string) error {
	err := ss.repo.Delete(ctx, token)
	if err != nil {
		return fmt.Errorf("[Delete]: %w", err)
	}
	return nil
}
