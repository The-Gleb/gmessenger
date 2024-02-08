package entity

import "time"

const (
	SENT      = "sent"
	DELIVERED = "delivered"
	READ      = "read"
)

type Message struct {
	ID        int64     `json:"id"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"time"`
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
