package entity

type Chat struct {
	Type         string
	ReceiverID   string
	ReceiverName string
	LastMessage  Message
	Unread       int64
}
