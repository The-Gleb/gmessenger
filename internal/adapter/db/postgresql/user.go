package postgresql

import (
	"context"

	"github.com/The-Gleb/gmessenger/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/errors"
	"github.com/The-Gleb/gmessenger/pkg/client/postgresql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type userStorage struct {
	// client Client
	sqlc *sqlc.Queries
}

func NewUserStorage(client postgresql.Client) *userStorage {
	return &userStorage{
		// client: client,
		sqlc: sqlc.New(client),
	}
}

func (us *userStorage) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	user, err := us.sqlc.GetUser(ctx, login)
	switch err {
	case pgx.ErrNoRows:
		return entity.User{}, errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetByLogin]: user not found")
	default:
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.GetByLogin]")
	case nil:
		return entity.User{
			UserName: user.Username.String,
			Login:    user.Login,
			Password: user.Password.String,
		}, nil
	}

}

func (us *userStorage) Create(ctx context.Context, user entity.User) (entity.User, error) {

	_, err := us.sqlc.GetUser(ctx, user.Login)
	switch err {
	case pgx.ErrNoRows:
		break
	case nil:
		return entity.User{}, errors.NewDomainError(errors.ErrDBLoginAlredyExists, "[storage.Create]:")
	default:
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.Create]")
	}

	params := sqlc.CreateUserParams{
		Username: pgtype.Text{
			String: user.UserName,
			Valid:  true,
		},
		Login: user.Login,
		Password: pgtype.Text{
			String: user.Password,
			Valid:  true,
		},
	}
	sqlcUser, err := us.sqlc.CreateUser(ctx, params)
	if err != nil {
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.Create]")
	}

	return entity.User{
		UserName: sqlcUser.Username.String,
		Login:    sqlcUser.Login,
		Password: sqlcUser.Password.String,
	}, nil
}

func (us *userStorage) GetPassword(ctx context.Context, login string) (string, error) {
	password, err := us.sqlc.GetPassword(ctx, login)
	switch err {
	case pgx.ErrNoRows:
		return "", errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetPassworc]: user not found")
	case nil:
		return password.String, nil
	default:
		return "", errors.NewDomainError(errors.ErrDB, "[storage.GetByLogin]")
	}

}
