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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (ss *sessionStorage) GetByToken(ctx context.Context, token string) (entity.Session, error) {

	row := ss.client.QueryRow(
		ctx,
		`SELECT user_id, expiry FROM sessions WHERE session_token = $1`,
		token,
	)

	session := entity.Session{
		Token: token,
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

func (ss *sessionStorage) Create(ctx context.Context, session entity.Session) error {

	_, err := ss.client.Exec(
		ctx,
		`INSERT INTO sessions (user_id, session_token, expiry) VALUES ($1, $2, $3)`,
		session.UserID, session.Token, session.Expiry,
	)
	if err != nil {
		slog.Error(err.Error())
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.NewDomainError(errors.ErrNotUniqueToken, "[storage.Create]")
		}
		return errors.NewDomainError(errors.ErrDB, "[storage.Create]:")
	}

	return nil

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
