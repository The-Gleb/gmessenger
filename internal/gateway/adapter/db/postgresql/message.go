package postgresql

import (
	"context"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/service"
	"github.com/The-Gleb/gmessenger/internal/gateway/errors"
	"github.com/The-Gleb/gmessenger/pkg/client/postgresql"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

var _ service.MessageStorage = new(messageStorage)

type messageStorage struct {
	client postgresql.Client
}

func NewMessageStorage(client postgresql.Client) *messageStorage {
	return &messageStorage{
		client: client,
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

}

func (ms *messageStorage) GetByUsers(ctx context.Context, senderID, receiverID int64, limit, offset int) ([]entity.Message, error) {

	rows, err := ms.client.Query(
		ctx,
		`SELECT * FROM messages
			WHERE (sender_id = $1 AND receiver_id = $2)
			OR (sender_id = $2 AND receiver_id = $1)
			ORDER BY id
			LIMIT $3 OFFSET $4;`,
	)

	if err != nil {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, "[messageStorage.GetByUsers]")
	}

	messages, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Message, error) {
		var msg entity.Message
		err := row.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Text, &msg.Status)
		return msg, err
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, "[messageStorage.GetByUsers]")
	}

	return messages, nil

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
