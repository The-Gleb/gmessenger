package entity

type Chat struct {
	Type         string  `json:"type"`
	ReceiverID   string  `json:"receiver_id"`
	ReceiverName string  `json:"receiver_name"`
	LastMessage  Message `json:"last_message"`
	Unread       int64   `json:"unread"`
}
