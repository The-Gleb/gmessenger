package entity

type Chat struct {
	Type        string
	ID          int64
	Name        string
	LastMessage string
	Unread      int
}

type ChatView struct {
	Type        string
	ID          int64
	Name        string
	LastMessage string
	Unread      int
}
