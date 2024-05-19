package entity

const GROUP = "GROUP"
const DIALOG = "DIALOG"

type Chat struct {
	Type         string  `json:"type"`
	ReceiverID   int64   `json:"receiver_id"`
	ReceiverName string  `json:"receiver_name"`
	LastMessage  Message `json:"last_message,omitempty"`
	Unread       int64   `json:"unread"`
}
