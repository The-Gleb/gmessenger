package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/gmessenger/app/internal/errors"
)

type Key string

type AuthUsecase interface {
	Auth(ctx context.Context, token string) (string, error)
}

type authMiddleWare struct {
	usecase AuthUsecase
}

func NewAuthMiddleware(usecase AuthUsecase) *authMiddleWare {
	return &authMiddleWare{usecase}
}

func (m *authMiddleWare) Http(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("auth middleware working")
		c, err := r.Cookie("sessionToken")

		if err != nil {
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}
		slog.Debug("Cookie is", "cookie", c.Value)

		userLogin, err := m.usecase.Auth(r.Context(), c.Value)
		if err != nil {
			slog.Debug("got userlogig from auth usecase", "login", userLogin)
			slog.Error(err.Error())
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), Key("userLogin"), userLogin)
		ctx = context.WithValue(ctx, Key("token"), c.Value)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *authMiddleWare) Websocket(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("websocket auth middleware working")

		token := r.URL.Query().Get("token")
		if token == "" {
			slog.Error("there is no token in a query")
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}
		slog.Debug("got token", "token", token)

		userLogin, err := m.usecase.Auth(r.Context(), token)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), Key("userLogin"), userLogin)
		ctx = context.WithValue(ctx, Key("token"), token)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
