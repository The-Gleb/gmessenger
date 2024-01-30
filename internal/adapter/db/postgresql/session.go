package postgresql

import (
	"context"
	stdErrors "errors"
	"log/slog"

	"github.com/The-Gleb/gmessenger/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/errors"
	"github.com/The-Gleb/gmessenger/pkg/client/postgresql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type sessionStorage struct {
	sqlc *sqlc.Queries
}

func NewSessionStorage(client postgresql.Client) *sessionStorage {
	return &sessionStorage{
		sqlc: sqlc.New(client),
	}
}

func (ss *sessionStorage) GetByToken(ctx context.Context, token string) (entity.Session, error) {

	s, err := ss.sqlc.GetSessionByToken(ctx, token)
	switch err {
	default:
		slog.Error(err.Error())
		return entity.Session{}, errors.NewDomainError(errors.ErrDB, "[storage.GetByToken]: ")
	case pgx.ErrNoRows:
		slog.Error(err.Error())
		return entity.Session{}, errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetByToken]: session not found")

	case nil:
		return entity.Session{
			Token:     s.Token,
			UserLogin: s.UserLogin.String,
			Expiry:    s.Expiry.Time,
		}, nil
	}

}

func (ss *sessionStorage) Create(ctx context.Context, session entity.Session) error {

	err := ss.sqlc.CreateSession(ctx, sqlc.CreateSessionParams{
		Token: session.Token,
		UserLogin: pgtype.Text{
			String: session.UserLogin,
			Valid:  true,
		},
		Expiry: pgtype.Timestamp{
			Time:  session.Expiry,
			Valid: true,
		},
	})
	if err != nil {
		slog.Error(err.Error())
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.NewDomainError(errors.ErrNotUniqueToken, "[storage.Create]: not unique token")
		}
		return errors.NewDomainError(errors.ErrDB, "[storage.Create]: ")
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
