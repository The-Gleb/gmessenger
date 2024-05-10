package entity

import "time"

type OTP struct {
	Token  string
	Data   OTPData
	Expiry time.Time
}

type OTPData struct {
	UserID    int64
	SessionID int64
}
