package postgresql

import (
	"context"
	stdErrors "errors"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service"
	"log/slog"

	"github.com/The-Gleb/gmessenger/app/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
	"github.com/jackc/pgx/v5"
)

var _ service.SessionStorage = new(sessionStorage)

type sessionStorage struct {
	client postgresql.Client
	sqlc   *sqlc.Queries
}

func NewSessionStorage(client postgresql.Client) *sessionStorage {
	return &sessionStorage{
		client: client,
		sqlc:   sqlc.New(client),
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
	err := ss.sqlc.DeleteSession(ctx, token)
	if err != nil {
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return errors.NewDomainError(errors.ErrNoDataFound, "token not found")
		}
	}

	return nil
}
