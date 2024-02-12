package entity

type Chat struct {
	Type          string
	GroupID       int64
	ReceiverLogin string
	Name          string
	LastMessage   Message
	Unread        int64
}
