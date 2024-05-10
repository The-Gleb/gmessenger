package postgresql

import (
	"context"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service"
	"github.com/jackc/pgx/v5/pgtype"
	"log/slog"

	stdErrors "errors"

	"github.com/The-Gleb/gmessenger/app/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ service.UserStorage = new(userStorage)

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

func (us userStorage) GetOrCreateByEmail(ctx context.Context, email string) (int64, bool, error) {
	tx, err := us.client.Begin(ctx)
	if err != nil {
		slog.Error(err.Error())
		return 0, false, errors.NewDomainError(errors.ErrDB, "[storage.GetOrCreateByEmail]")
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	row := tx.QueryRow(
		ctx,
		`SELECT (id) FROM users WHERE email = $1`,
		email,
	)

	var id int64
	err = row.Scan(&id)
	if err == nil {
		return id, false, nil
	}
	if !stdErrors.Is(err, pgx.ErrNoRows) {
		slog.Error(err.Error())
		return 0, false, errors.NewDomainError(errors.ErrDB, "[storage.GetOrCreateByEmail]")
	}

	row = tx.QueryRow(
		ctx,
		`
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING id;
		`,
		email,
	)

	err = row.Scan(&id)
	if err != nil {
		slog.Error(err.Error())
		return 0, false, errors.NewDomainError(errors.ErrDB, "[storage.GetOrCreateByEmail]")
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error(err.Error())
		return 0, false, errors.NewDomainError(errors.ErrDB, "[storage.GetOrCreateByEmail]")
	}

	return id, true, nil

}

func (us userStorage) SetUsername(ctx context.Context, dto entity.SetUsernameDTO) error {
	_, err := us.client.Exec(
		ctx,
		`UPDATE users SET username = $1 WHERE id = $2`,
		dto.Username, dto.UserID,
	)
	if err != nil {
		slog.Error(err.Error())
		return errors.NewDomainError(errors.ErrDB, "[storage.SetUsername]")
	}
	return nil
}

func (us userStorage) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	user, err := us.sqlc.GetUser(ctx, login)
	switch err {
	default:

		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.GetByLogin]")

	case pgx.ErrNoRows:
		return entity.User{}, errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetByLogin]: user not found")

	case nil:
		return entity.User{
			Username: user.Username.String,
			Email:    user.Login,
			Password: user.Password.String,
		}, nil
	}

}

func (us userStorage) CreateWithPassword(ctx context.Context, dto entity.RegisterUserDTO) (int64, error) {

	tx, err := us.client.Begin(ctx)
	if err != nil {
		slog.Error(err.Error())
		return 0, errors.NewDomainError(errors.ErrDB, "[storage.CreateWithPassword]")
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	row := tx.QueryRow(
		ctx,
		`INSERT INTO users (email, username)
		VALUES ($1, $2)
		RETURNING id;`,
		dto.Email, dto.Username,
	)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		slog.Error(err.Error())
		var pgErr *pgconn.PgError
		if stdErrors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return 0, errors.NewDomainError(errors.ErrUserExists, "[storage.CreateWithPassword]")
			}
		}
		return 0, errors.NewDomainError(errors.ErrDB, "[storage.CreateWithPassword]")
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO user_password (user_id, password) VALUES ($1, $2);`,
		id, dto.Password,
	)
	if err != nil {
		slog.Error(err.Error())
		return 0, errors.NewDomainError(errors.ErrDB, "[storage.CreateWithPassword]")
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error(err.Error())
		return 0, errors.NewDomainError(errors.ErrDB, "[storage.CreateWithPassword]")
	}

	return id, nil

}

func (us userStorage) GetUserInfoByID(ctx context.Context, ID int64) (entity.UserInfo, error) {

	row := us.client.QueryRow(
		ctx,
		`SELECT username, email FROM users WHERE id = $1;`,
		ID,
	)

	userInfo := entity.UserInfo{
		ID: ID,
	}
	err := row.Scan(&userInfo.Username, &userInfo.Email)
	if err != nil {
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return entity.UserInfo{}, errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetUserInfoByID]")
		}
		return entity.UserInfo{}, errors.NewDomainError(errors.ErrDB, "[storage.GetUserInfoByID]")
	}

	return userInfo, nil
}

func (us userStorage) GetByEmail(ctx context.Context, email string) (entity.User, error) {

	tx, err := us.client.Begin(ctx)
	if err != nil {
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.GetPasswordByEmail]")
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	row := tx.QueryRow(
		ctx,
		`SELECT id FROM users WHERE email = $1;`,
		email,
	)

	user := entity.User{
		Email: email,
	}
	err = row.Scan(&user.ID)
	if err != nil {
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetPasswordByEmail]")
		}
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.GetPasswordByEmail]")
	}

	row = tx.QueryRow(
		ctx,
		`SELECT password FROM user_password WHERE user_id = $1;`,
		user.ID,
	)

	err = row.Scan(&user.Password)
	if err != nil {
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errors.NewDomainError(errors.ErrNoDataFound, "[storage.GetPasswordByEmail]")
		}
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.GetPasswordByEmail]")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.User{}, errors.NewDomainError(errors.ErrDB, "[storage.GetPasswordByEmail]")
	}

	return user, nil

}

func (us userStorage) GetAllUsersView(ctx context.Context) ([]entity.UserView, error) {
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
	// TODO:
	return nil, nil
}

func (us userStorage) GetChatsView(ctx context.Context, userID int64) ([]entity.Chat, error) {
	// TODO: refactor

	tx, err := us.client.Begin(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, "[storage.GetChatsView]")
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	rows, err := tx.Query(
		ctx,
		`
			SELECT users.id, users.username, messages.id, messages.sender_id, messages.receiver_id
			FROM users LEFT JOIN messages
			ON
			    (sender_id = users.id AND receiver_id = $1)
				OR (receiver_id = users.id AND sender_id = $1);
		`, userID,
	)

	chats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Chat, error) {
		var chat entity.Chat
		var messageID pgtype.Int8
		var receiverID pgtype.Int8
		var senderID pgtype.Int8

		err := row.Scan(&chat.ReceiverID, &chat.ReceiverName,
			&messageID, &receiverID,
			&senderID)

		chat.Type = entity.PERSONAL_CHAT
		chat.LastMessage.ID = messageID.Int64
		chat.LastMessage.ReceiverID = receiverID.Int64
		chat.LastMessage.SenderID = senderID.Int64

		return chat, err
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, "[storage.GetChatsView]")
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, "[storage.GetChatsView]")
	}

	return chats, nil

	//dbTx, err := us.client.Begin(ctx)
	//if err != nil {
	//	return []entity.Chat{}, err
	//}
	//defer dbTx.Rollback(ctx) //nolint:all
	//
	//sqlcTx := us.sqlc.WithTx(dbTx)
	//
	//sqlcUsers, err := sqlcTx.GetAllUsers(ctx)
	//if err != nil {
	//	slog.Error(err.Error())
	//	return []entity.Chat{}, errors.NewDomainError(errors.ErrDB, "")
	//}
	//
	//chatsView := make([]entity.Chat, len(sqlcUsers))
	//for i, sqlcUser := range sqlcUsers {
	//	chatsView[i].Type = client.Dialog
	//	chatsView[i].ReceiverLogin = sqlcUser.Login
	//	chatsView[i].Name = sqlcUser.Username.String
	//	sqlcLastMsg, err := sqlcTx.GetLastMessage(ctx, sqlc.GetLastMessageParams{
	//		Sender: sqlcUser.Username,
	//		Receiver: pgtype.Text{
	//			String: userLogin,
	//			Valid:  true,
	//		},
	//		Limit:  1,
	//		Offset: 0,
	//	})
	//	if err != nil {
	//		if err == pgx.ErrNoRows {
	//			chatsView[i].Unread = 0
	//			continue
	//		}
	//		slog.Error(err.Error())
	//		return []entity.Chat{}, errors.NewDomainError(errors.ErrDB, "")
	//	}
	//	chatsView[i].LastMessage = entity.Message{
	//		ID:        sqlcLastMsg.ID,
	//		Sender:    sqlcLastMsg.Sender.String,
	//		Text:      sqlcLastMsg.Text.String,
	//		Status:    sqlcLastMsg.Status.String,
	//		Timestamp: sqlcLastMsg.CreatedAt.Time,
	//	}
	//
	//	unreadNumber, err := sqlcTx.GetUnreadNumber(ctx, sqlc.GetUnreadNumberParams{
	//		Sender: sqlcUser.Username,
	//		Receiver: pgtype.Text{
	//			String: userLogin,
	//			Valid:  true,
	//		},
	//	})
	//	if err != nil {
	//		slog.Error(err.Error())
	//		return []entity.Chat{}, errors.NewDomainError(errors.ErrDB, "")
	//	}
	//	chatsView[i].Unread = unreadNumber
	//}
	//
	//err = dbTx.Commit(ctx)
	//if err != nil {
	//	slog.Error(err.Error())
	//	return []entity.Chat{}, errors.NewDomainError(errors.ErrDB, "")
	//}

	return nil, nil

}
