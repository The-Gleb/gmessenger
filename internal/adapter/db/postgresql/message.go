package postgresql

import (
	"context"
	"log/slog"

	"github.com/The-Gleb/gmessenger/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/pkg/client/postgresql"
	"github.com/jackc/pgx/v5/pgtype"
)

type messageStorage struct {
	sqlc *sqlc.Queries
}

func NewMessageStorage(client postgresql.Client) *messageStorage {
	return &messageStorage{
		sqlc: sqlc.New(client),
	}
}

func (ms *messageStorage) GetByID(ctx context.Context, msgID int64) (entity.Message, error) {

	m, err := ms.sqlc.GetMessageByID(ctx, msgID)
	if err != nil {

	}

	return entity.Message{
		ID:        m.ID,
		Sender:    m.Sender.String,
		Receiver:  m.Receiver.String,
		Text:      m.Text.String,
		Status:    m.Status.String,
		Timestamp: m.CreatedAt.Time,
	}, nil

}

func (ms *messageStorage) Create(ctx context.Context, msg entity.Message) (entity.Message, error) {

	m, err := ms.sqlc.CreateMessage(ctx, sqlc.CreateMessageParams{
		Sender: pgtype.Text{
			String: msg.Sender,
			Valid:  true,
		},
		Receiver: pgtype.Text{
			String: msg.Receiver,
			Valid:  true,
		},
		Text: pgtype.Text{
			String: msg.Text,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: msg.Status,
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  msg.Timestamp,
			Valid: true,
		},
	})
	if err != nil {
		slog.Error("couldn`t add new message to db", "error", err)
	}

	return entity.Message{
		ID:        m.ID,
		Sender:    m.Sender.String,
		Receiver:  m.Receiver.String,
		Text:      m.Text.String,
		Status:    m.Status.String,
		Timestamp: m.CreatedAt.Time,
	}, nil

}

func (ms *messageStorage) GetByUsers(ctx context.Context, sender, receiver string) ([]entity.Message, error) {
	sqlcMessages, err := ms.sqlc.GetMessagesByUsers(ctx, sqlc.GetMessagesByUsersParams{
		Sender: pgtype.Text{
			String: sender,
			Valid:  true,
		},
		Receiver: pgtype.Text{
			String: receiver,
			Valid:  true,
		},
	})
	if err != nil {

	}
	messages := make([]entity.Message, len(sqlcMessages))

	for i, m := range sqlcMessages {
		messages[i] = entity.Message{
			ID:        m.ID,
			Sender:    m.Sender.String,
			Receiver:  m.Receiver.String,
			Text:      m.Text.String,
			Status:    m.Status.String,
			Timestamp: m.CreatedAt.Time,
		}
	}

	return messages, nil
}

func (ms *messageStorage) UpdateStatus(ctx context.Context, msgID int64, status string) (entity.Message, error) {

	m, err := ms.sqlc.UpdateMessageStatus(ctx, sqlc.UpdateMessageStatusParams{
		Status: pgtype.Text{
			String: status,
			Valid:  true,
		},
		ID: msgID,
	})
	if err != nil {

	}

	return entity.Message{
		ID:        m.ID,
		Sender:    m.Sender.String,
		Receiver:  m.Receiver.String,
		Text:      m.Text.String,
		Status:    m.Status.String,
		Timestamp: m.CreatedAt.Time,
	}, nil
}
