package entity

import (
	"github.com/vk-rv/pvx"
	"time"
)

type TokenData struct {
	Subject    string
	Expiration time.Time
	AdditionalClaims
	Footer
}

type AdditionalClaims struct {
	UserID    int64 `json:"user_id"`
	SessionID int64 `json:"session_id"`
}

type Footer struct {
	MetaData string `json:"meta_data"`
}

type ServiceClaims struct {
	pvx.RegisteredClaims
	AdditionalClaims
	Footer
}
