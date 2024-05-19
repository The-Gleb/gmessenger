package postgresql

import (
	"context"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"log/slog"

	// stdErrors "errors"
	"github.com/The-Gleb/gmessenger/app/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
)

var _ service.MessageStorage = new(messageStorage)

type messageStorage struct {
	client postgresql.Client
	sqlc   *sqlc.Queries
}

func NewMessageStorage(client postgresql.Client) *messageStorage {
	return &messageStorage{
		client: client,
		sqlc:   sqlc.New(client),
	}
}

func (ms *messageStorage) GetByID(ctx context.Context, msgID int64) (entity.Message, error) {

	//m, err := ms.sqlc.GetMessageByID(ctx, msgID)
	//if err != nil {
	//	slog.Error(err.Error())
	//
	//	if err == pgx.ErrNoRows {
	//		return entity.Message{}, errors.NewDomainError(errors.ErrNoDataFound, "")
	//	}
	//
	//	return entity.Message{}, errors.NewDomainError(errors.ErrDB, "")
	//}
	//
	//return entity.Message{
	//	ID:        m.ID,
	//	SenderID:    m.Sender.Int,
	//	ReceiverID:  m.Receiver.String,
	//	Text:      m.Text.String,
	//	Status:    m.Status.String,
	//	Timestamp: m.CreatedAt.Time,
	//}, nil

	return entity.Message{}, nil

}

func (ms *messageStorage) Create(ctx context.Context, msg entity.Message) (entity.Message, error) {

	row := ms.client.QueryRow(
		ctx,
		`INSERT INTO messages (sender_id, receiver_id, text, status, created_at) VALUES ($1, $2, $3, $4,$5)
		RETURNING id;`,
		msg.SenderID, msg.ReceiverID, msg.Text, msg.Status, msg.Timestamp,
	)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		slog.Error(err.Error())
		return entity.Message{}, errors.NewDomainError(errors.ErrDB, "[messageStorage.Create]")
	}
	msg.ID = id
	return msg, nil

	//m, err := ms.sqlc.CreateMessage(ctx, sqlc.CreateMessageParams{
	//	Sender: pgtype.Text{
	//		String: msg.Sender,
	//		Valid:  true,
	//	},
	//	Receiver: pgtype.Text{
	//		String: msg.Receiver,
	//		Valid:  true,
	//	},
	//	Text: pgtype.Text{
	//		String: msg.Text,
	//		Valid:  true,
	//	},
	//	Status: pgtype.Text{
	//		String: msg.Status,
	//		Valid:  true,
	//	},
	//	CreatedAt: pgtype.Timestamp{
	//		Time:  time.Now(),
	//		Valid: true,
	//	},
	//})
	//if err != nil {
	//	slog.Error("couldn`t add new message to db", "error", err)
	//	return entity.Message{}, errors.NewDomainError(errors.ErrDB, "couldn`t add new message to db")
	//}
	//
	//return entity.Message{
	//	ID:        m.ID,
	//	Sender:    m.Sender.String,
	//	Receiver:  m.Receiver.String,
	//	Text:      m.Text.String,
	//	Status:    m.Status.String,
	//	Timestamp: m.CreatedAt.Time,
	//}, nil
	return entity.Message{}, nil

}

func (ms *messageStorage) GetByUsers(ctx context.Context, senderID, receiverID int64, limit, offset int) ([]entity.Message, error) {
	//sqlcMessages, err := ms.sqlc.GetMessagesByUsers(ctx, sqlc.GetMessagesByUsersParams{
	//	Sender: pgtype.Text{
	//		String: sender,
	//		Valid:  true,
	//	},
	//	Receiver: pgtype.Text{
	//		String: receiver,
	//		Valid:  true,
	//	},
	//	Limit:  int32(limit),
	//	Offset: int32(offset),
	//})
	//if err != nil {
	//	slog.Error("error getting messages from db", "error", err.Error())
	//	return []entity.Message{}, errors.NewDomainError(errors.ErrDB, "")
	//}
	//messages := make([]entity.Message, len(sqlcMessages))
	//
	//for i, m := range sqlcMessages {
	//	messages[i] = entity.Message{
	//		ID:        m.ID,
	//		Sender:    m.Sender.String,
	//		Receiver:  m.Receiver.String,
	//		Text:      m.Text.String,
	//		Status:    m.Status.String,
	//		Timestamp: m.CreatedAt.Time,
	//	}
	//}
	//
	//return messages, nil
	return nil, nil
}

func (ms *messageStorage) UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) {

	//m, err := ms.sqlc.UpdateMessageStatus(ctx, sqlc.UpdateMessageStatusParams{
	//	Status: pgtype.Text{
	//		String: status,
	//		Valid:  true,
	//	},
	//	ID: msgID,
	//})
	//if err != nil {
	//	slog.Error("error updating message status in db", "error", err.Error())
	//
	//	if err == pgx.ErrNoRows {
	//		return entity.Message{}, errors.NewDomainError(errors.ErrNoDataFound, "")
	//	}
	//
	//	return entity.Message{}, errors.NewDomainError(errors.ErrDB, "")
	//}
	//
	//return entity.Message{
	//	ID:        m.ID,
	//	Sender:    m.Sender.String,
	//	Receiver:  m.Receiver.String,
	//	Text:      m.Text.String,
	//	Status:    m.Status.String,
	//	Timestamp: m.CreatedAt.Time,
	//}, nil
	return entity.Message{}, nil
}
