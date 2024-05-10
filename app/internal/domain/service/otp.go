package service

import (
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/google/uuid"
	"sync"
	"time"
)

type otpService struct {
	otps map[string]entity.OTP
	TTL  time.Duration
	m    *sync.Mutex
}

func NewOtpService() *otpService {
	return &otpService{
		otps: make(map[string]entity.OTP),
	}
}

func (s *otpService) GenerateOtp(userID, sessionID int64) string {
	s.m.Lock()
	defer s.m.Unlock()
	token := uuid.NewString()
	for _, ok := s.otps[token]; ok; _, ok = s.otps[token] {
		token = uuid.NewString()
	}
	s.otps[token] = entity.OTP{
		Token: token,
		Data: entity.OTPData{
			UserID:    userID,
			SessionID: sessionID,
		},
		Expiry: time.Now().Add(s.TTL),
	}
	return token
}

func (s *otpService) VerifyOtp(token string) (entity.OTPData, bool) {
	s.m.Lock()
	defer s.m.Unlock()
	otp, ok := s.otps[token]
	if !ok {
		return entity.OTPData{}, false
	}
	delete(s.otps, token)
	if otp.Expiry.Before(time.Now()) {
		return entity.OTPData{}, false
	}
	return otp.Data, true
}
