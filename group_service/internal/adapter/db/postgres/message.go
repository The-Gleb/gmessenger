package postgres

import (
	"context"
	"log/slog"

	"github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
	"github.com/The-Gleb/gmessenger/group_service/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/group_service/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/group_service/internal/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

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

func (s *messageStorage) AddMessage(ctx context.Context, message entity.Message) (entity.Message, error) {

	sqlcMessage, err := s.sqlc.AddMessage(ctx, sqlc.AddMessageParams{
		Sender: pgtype.Text{
			String: message.Sender,
			Valid:  true,
		},
		Text: pgtype.Text{
			String: message.Text,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: message.Status,
			Valid:  true,
		},
		GroupID: pgtype.Int8{
			Int64: message.GroupID,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  message.Timestamp,
			Valid: true,
		},
	})
	if err != nil {
		return entity.Message{}, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	return entity.Message{
		ID:        sqlcMessage.ID,
		GroupID:   sqlcMessage.GroupID.Int64,
		Sender:    sqlcMessage.Sender.String,
		Text:      sqlcMessage.Text.String,
		Status:    sqlcMessage.Status.String,
		Timestamp: sqlcMessage.CreatedAt.Time,
	}, nil

}

func (s *messageStorage) GetMessages(ctx context.Context, groupID int64, limit, offset int) ([]entity.Message, error) {

	sqlcMessages, err := s.sqlc.GetMessages(ctx, sqlc.GetMessagesParams{
		GroupID: pgtype.Int8{
			Int64: groupID,
			Valid: true,
		},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	messages := make([]entity.Message, len(sqlcMessages))

	for i, sqlcMsg := range sqlcMessages {
		messages[i] = entity.Message{
			ID:        sqlcMsg.ID,
			GroupID:   sqlcMsg.GroupID.Int64,
			Sender:    sqlcMsg.Sender.String,
			Text:      sqlcMsg.Text.String,
			Status:    sqlcMsg.Status.String,
			Timestamp: sqlcMsg.CreatedAt.Time,
		}
	}

	return messages, nil

}

func (s *messageStorage) UpdateMessageStatus(ctx context.Context, messageID int64, status string) (entity.Message, error) {

	m, err := s.sqlc.UpdateMessageStatus(ctx, sqlc.UpdateMessageStatusParams{
		Status: pgtype.Text{
			String: status,
			Valid:  true,
		},
		ID: messageID,
	})
	if err != nil {
		slog.Error("error updating message status in db", "error", err.Error())

		if err == pgx.ErrNoRows {
			return entity.Message{}, errors.NewDomainError(errors.ErrNoDataFound, "")
		}

		return entity.Message{}, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	return entity.Message{
		ID:        m.ID,
		Sender:    m.Sender.String,
		GroupID:   m.GroupID.Int64,
		Text:      m.Text.String,
		Status:    m.Status.String,
		Timestamp: m.CreatedAt.Time,
	}, nil

}

func (s *messageStorage) GetLastMessage(ctx context.Context, groupID int64) (entity.Message, error) {

	sqlcMessage, err := s.sqlc.GetMessages(ctx, sqlc.GetMessagesParams{
		GroupID: pgtype.Int8{
			Int64: groupID,
			Valid: true,
		},
		Limit:  1,
		Offset: 0,
	})
	if err != nil && err != pgx.ErrNoRows {
		return entity.Message{}, err
	}

	if len(sqlcMessage) == 0 {
		return entity.Message{}, errors.NewDomainError(errors.ErrNoDataFound, "there is no last message")
	}

	return entity.Message{
		ID:        sqlcMessage[0].ID,
		GroupID:   sqlcMessage[0].GroupID.Int64,
		Sender:    sqlcMessage[0].Sender.String,
		Text:      sqlcMessage[0].Text.String,
		Status:    sqlcMessage[0].Status.String,
		Timestamp: sqlcMessage[0].CreatedAt.Time,
	}, nil

}
