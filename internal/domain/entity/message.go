package entity

import "time"

type Message struct {
	ID        int64
	Sender    string
	Receiver  string
	Text      string
	Status    string
	Timestamp time.Time
}

func (msg *Message) SetSent() {
	msg.Status = "SENT"
}

func (msg *Message) SetDelivered() {
	msg.Status = "DELIVERED"
}

func (msg *Message) SetRead() {
	msg.Status = "READ"
}
