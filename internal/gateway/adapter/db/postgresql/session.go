package postgresql

import (
	"context"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/service"
	"github.com/The-Gleb/gmessenger/internal/gateway/errors"
	"github.com/The-Gleb/gmessenger/pkg/client/postgresql"
	"log/slog"

	stdErrors "errors"

	"github.com/jackc/pgx/v5"
)

var _ service.SessionStorage = new(sessionStorage)

type sessionStorage struct {
	client postgresql.Client
}

func NewSessionStorage(client postgresql.Client) *sessionStorage {
	return &sessionStorage{
		client: client,
	}
}

func (ss *sessionStorage) GetByID(ctx context.Context, ID int64) (entity.Session, error) {

	row := ss.client.QueryRow(
		ctx,
		`SELECT user_id, expiry FROM sessions WHERE session_token = $1`,
		ID,
	)

	session := entity.Session{
		ID: ID,
	}

	err := row.Scan(&session.UserID, &session.Expiry)
	if err != nil {
		slog.Error(err.Error())
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return session, errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetByToken]:")
		}
		return session, errors.NewDomainError(errors.ErrDB, "[storage.GetByToken]:")
	}

	return session, nil

}

func (ss *sessionStorage) Create(ctx context.Context, session entity.Session) (entity.Session, error) {

	row := ss.client.QueryRow(
		ctx,
		`INSERT INTO sessions (user_id, expiry) VALUES ($1, $2)
		RETURNING id;`,
		session.UserID, session.Expiry,
	)

	err := row.Scan(&session.ID)

	if err != nil {
		slog.Error(err.Error())
		return entity.Session{}, errors.NewDomainError(errors.ErrDB, "[storage.Create]:")
	}

	return session, nil

}

func (ss *sessionStorage) Delete(ctx context.Context, token string) error {
	_, err := ss.client.Exec(
		ctx,
		`DELETE FROM sessions WHERE session_token = $1`,
		token)
	if err != nil {
		slog.Error(err.Error())
		return errors.NewDomainError(errors.ErrDB, "[sessionStorage.Delete]:")
	}
	//if c.RowsAffected() == 0 {
	//	slog.Error("no rows affected")
	//	return errors.NewDomainError(errors.ErrNoDataFound, "token not found")
	//}

	return nil
}
