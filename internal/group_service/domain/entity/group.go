package entity

import "time"

type Group struct {
	ID        int64
	Name      string
	MemberIDs []int64
	CreatedAt time.Time
}

type CreateGroupDTO struct {
	Name      string
	MemberIDs []int64
}

type GroupView struct {
	ID          int64
	Name        string
	LastMessage Message
	Unread      int
}
