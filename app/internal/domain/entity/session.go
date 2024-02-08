package entity

import "time"

type Session struct {
	UserLogin string
	Token     string
	Expiry    time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
