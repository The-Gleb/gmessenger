package entity

import "time"

const (
	SENT      = "SENT"
	DELIVERED = "DELIVERED"
	READ      = "READ"
)

type Message struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	GroupID    int64     `json:"group_id"`
	Text       string    `json:"text"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"time"`
}

type CreateMessageDTO struct {
	SenderID int64
	GroupID  int64
	Text     string
}

type GetGroupMessagesDTO struct {
	GroupID int64
	Limit   int
	Offset  int
}
