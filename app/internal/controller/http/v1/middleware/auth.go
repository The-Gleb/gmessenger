package v1

import (
	"context"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"log/slog"
	"net/http"
	"strings"

	"github.com/The-Gleb/gmessenger/app/internal/errors"
)

type Key string

type AuthUsecase interface {
	Auth(ctx context.Context, token string) (entity.AdditionalClaims, error)
}

type otpService interface {
	VerifyOtp(token string) (entity.OTPData, bool)
}

type authMiddleWare struct {
	pasetoUsecase AuthUsecase
	otpSvc        otpService
}

func NewAuthMiddleware(usecase AuthUsecase, otpService otpService) *authMiddleWare {
	return &authMiddleWare{
		usecase,
		otpService,
	}
}

func (m *authMiddleWare) Http(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("auth middleware working")

		//c, err := r.Cookie("sessionToken")
		//if err != nil {
		//	slog.Error("sessionToken cookie not found")
		//	http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
		//	return
		//}

		bearerToken := r.Header.Get("Authorization")
		reqToken := strings.Split(bearerToken, " ")[1]

		slog.Debug("Token is", "token", reqToken)

		userData, err := m.pasetoUsecase.Auth(r.Context(), reqToken)
		if err != nil {
			slog.Error(err.Error())
			//http.Redirect(w, r, "/static/login/login.html", http.StatusMovedPermanently)
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}
		slog.Debug("got userID from auth usecase", "ID", userData.UserID)

		ctx := context.WithValue(r.Context(), Key("userID"), userData.UserID)
		ctx = context.WithValue(r.Context(), Key("sessionID"), userData.SessionID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *authMiddleWare) Websocket(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("websocket auth middleware working")

		token := r.URL.Query().Get("otp")
		if token == "" {
			slog.Error("there is no token in a query")
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}
		slog.Debug("got token", "token", token)

		data, ok := m.otpSvc.VerifyOtp(token)
		if !ok {
			slog.Error("otp verification failed")
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), Key("userID"), data.UserID)
		ctx = context.WithValue(r.Context(), Key("sessionID"), data.SessionID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
