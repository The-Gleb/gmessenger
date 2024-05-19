package entity

import "time"

type Session struct {
	ID     int64
	UserID int64
	Expiry time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
