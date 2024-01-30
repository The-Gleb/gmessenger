package dto

type OpenChatDTO struct {
	ChatType      string `json:"type"`
	ChatID        string `json:"id"`
	SenderLogin   string `json:"sender_login"`
	ReceiverLogin string `json:"receiver_login,omitempty"`
}
