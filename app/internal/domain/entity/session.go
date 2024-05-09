package entity

import "time"

type Session struct {
	UserID int64
	Token  string
	Expiry time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
