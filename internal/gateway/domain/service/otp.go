package service

import (
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/google/uuid"
	"log/slog"
	"sync"
	"time"
)

type otpService struct {
	otps map[string]entity.OTP
	ttl  time.Duration
	mu   sync.Mutex
}

func NewOtpService(ttl time.Duration) *otpService {
	return &otpService{
		otps: make(map[string]entity.OTP),
		ttl:  ttl,
	}
}

func (s *otpService) GenerateOtp(userID, sessionID int64) string {
	s.mu.Lock()
	defer s.mu.Unlock()
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
		Expiry: time.Now().Add(s.ttl),
	}

	slog.Debug("otp is created", "token", s.otps[token])

	return token
}

func (s *otpService) VerifyOtp(token string) (entity.OTPData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	otp, ok := s.otps[token]
	if !ok {
		return entity.OTPData{}, false
	}

	delete(s.otps, token)
	if otp.Expiry.Before(time.Now()) {
		slog.Debug("otp is expired", "otp", otp, "timeNow", time.Now())
		return entity.OTPData{}, false
	}
	return otp.Data, true
}
