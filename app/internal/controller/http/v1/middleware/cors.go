package v1

import (
	"log/slog"
	"net/http"
)

type corsMiddleWare struct {
}

func NewCorsMiddleware(usecase AuthUsecase, otpService otpService) *corsMiddleWare {
	return &corsMiddleWare{}
}

func (m *corsMiddleWare) AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("cors middleware working")

		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
