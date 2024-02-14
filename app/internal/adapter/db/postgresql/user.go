package postgresql

import (
	"context"
	"log/slog"

	stdErrors "errors"

	"github.com/The-Gleb/gmessenger/app/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
	"github.com/jackc/pgerrcode"
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
	client Client
	sqlc   *sqlc.Queries
}

func NewUserStorage(client postgresql.Client) *userStorage {
	return &userStorage{
		client: client,
		sqlc:   sqlc.New(client),
	}
}

func (us *userStorage) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	user, err := us.sqlc.GetUser(ctx, login)
	switch err {
	default:
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.GetByLogin]")

	case pgx.ErrNoRows:
		return entity.User{}, errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetByLogin]: user not found")

	case nil:
		return entity.User{
			UserName: user.Username.String,
			Login:    user.Login,
			Password: user.Password.String,
		}, nil
	}

}

func (us *userStorage) Create(ctx context.Context, user entity.User) (entity.User, error) {

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
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return entity.User{}, errors.NewDomainError(errors.ErrDBLoginAlredyExists, "[storage.Create]:")
		}

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

func (us *userStorage) GetAllUsernames(ctx context.Context) ([]string, error) {
	sqlcNames, err := us.sqlc.GetAllUsernames(ctx)
	switch err {
	// case pgx.ErrNoRows:
	// 	return "", errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetPassworc]: user not found")
	// case nil:
	// 	return password.String, nil
	// default:
	// 	return "", errors.NewDomainError(errors.ErrDB, "[storage.GetByLogin]")
	}
	slog.Debug("got names from db", "struct", sqlcNames)
	usernames := make([]string, len(sqlcNames))

	for i, sqlcName := range sqlcNames {
		usernames[i] = sqlcName.String
	}

	return usernames, nil
}

func (us *userStorage) GetChatsView(ctx context.Context, userLogin string) ([]entity.Chat, error) {
	dbTx, err := us.client.Begin(ctx)
	if err != nil {
		return []entity.Chat{}, err
	}
	defer dbTx.Rollback(ctx) //nolint:all

	sqlcTx := us.sqlc.WithTx(dbTx)

	sqlcUsers, err := sqlcTx.GetAllUsers(ctx)
	if err != nil {
		slog.Error(err.Error())
		return []entity.Chat{}, errors.NewDomainError(errors.ErrDB, "")
	}

	chatsView := make([]entity.Chat, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		chatsView[i].Type = client.Dialog
		chatsView[i].ReceiverLogin = sqlcUser.Login
		chatsView[i].Name = sqlcUser.Username.String
		sqlcLastMsg, err := sqlcTx.GetLastMessage(ctx, sqlc.GetLastMessageParams{
			Sender: sqlcUser.Username,
			Receiver: pgtype.Text{
				String: userLogin,
				Valid:  true,
			},
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			if err == pgx.ErrNoRows {
				chatsView[i].Unread = 0
				continue
			}
			slog.Error(err.Error())
			return []entity.Chat{}, errors.NewDomainError(errors.ErrDB, "")
		}
		chatsView[i].LastMessage = entity.Message{
			ID:        sqlcLastMsg.ID,
			Sender:    sqlcLastMsg.Sender.String,
			Text:      sqlcLastMsg.Text.String,
			Status:    sqlcLastMsg.Status.String,
			Timestamp: sqlcLastMsg.CreatedAt.Time,
		}

		unreadNumber, err := sqlcTx.GetUnreadNumber(ctx, sqlc.GetUnreadNumberParams{
			Sender: sqlcUser.Username,
			Receiver: pgtype.Text{
				String: userLogin,
				Valid:  true,
			},
		})
		if err != nil {
			slog.Error(err.Error())
			return []entity.Chat{}, errors.NewDomainError(errors.ErrDB, "")
		}
		chatsView[i].Unread = unreadNumber
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		slog.Error(err.Error())
		return []entity.Chat{}, errors.NewDomainError(errors.ErrDB, "")
	}

	return chatsView, nil

}
