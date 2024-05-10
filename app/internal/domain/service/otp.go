package service

import (
	"github.com/google/uuid"
	"sync"
)

type otpService struct {
	otps map[string]int64
	m    *sync.Mutex
}

func NewOtpService() *otpService {
	return &otpService{
		otps: make(map[string]int64),
	}
}

func (s *otpService) GenerateOtp(userID int64) string {
	s.m.Lock()
	defer s.m.Unlock()
	otp := uuid.NewString()
	if _, ok := s.otps[otp]; !ok {
		s.otps[otp] = userID
	}
	return otp
}
