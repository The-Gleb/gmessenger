package entity

import "time"

const (
	NewMessage          = "new_message"
	SendMessage         = "send_message"
	UserActivity        = "user_activity"
	MessageStatus       = "message_status"
	DialogNotificatiohn = "dialog_notification"
)

type Event struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type SendDialogMessageEvent struct {
	Text string `json:"text"`
}

type NewDialogMessageEvent struct {
	ID          int64     `json:"id"`
	Status      string    `json:"status"`
	Text        string    `json:"text"`
	SenderLogin string    `json:"sender"`
	CreatedAt   time.Time `json:"time"`
}

type UpdateMessageStatusEvent struct {
	MessageID int64  `json:"message_id"`
	Status    string `json:"status"`
}

// type DialogNotificatiohnEvent struct {
// 	NewDialogMessageEvent
// }

type UserActivityEvent struct {
	UserLogin string
}
