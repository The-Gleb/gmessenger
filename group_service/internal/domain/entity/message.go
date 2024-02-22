package entity

import "time"

const (
	SENT      = "SENT"
	DELIVERED = "DELIVERED"
	READ      = "READ"
)

type Message struct {
	ID        int64     `json:"id"`
	Sender    string    `json:"sender"`
	GroupID   int64     `json:"group_id"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"time"`
}
