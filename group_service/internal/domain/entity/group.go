package entity

import "time"

type Group struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	// MembersLogins []string
}

type GroupCreate struct {
	Name          string
	CreatedAt     time.Time
	MembersLogins []string
}

type GroupView struct {
	ID          int64
	Name        string
	LastMessage Message
	Unread      int
}
