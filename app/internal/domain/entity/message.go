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
	ReceiverID int64     `json:"receiver_id"`
	Text       string    `json:"text"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"time"`
}

// move to service
func (msg *Message) SetSent() {
	msg.Status = "SENT"
}

func (msg *Message) SetDelivered() {
	msg.Status = "DELIVERED"
}

func (msg *Message) SetRead() {
	msg.Status = "READ"
}
