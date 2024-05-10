package entity

const GROUP_CHAT = "GROUP_CHAT"
const PERSONAL_CHAT = "PERSONAL_CHAT"

type Chat struct {
	Type         string  `json:"type"`
	ReceiverID   string  `json:"receiver_id"`
	ReceiverName string  `json:"receiver_name"`
	LastMessage  Message `json:"last_message,omitempty"`
	Unread       int64   `json:"unread"`
}
