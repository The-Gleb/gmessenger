package v1

import (
	"context"
	"net/http"

	"github.com/The-Gleb/gmessenger/internal/errors"
)

type AuthUsecase interface {
	Auth(ctx context.Context, token string) (string, error)
}

type authMiddleWare struct {
	usecase AuthUsecase
}

func NewAuthMiddleware(usecase AuthUsecase) *authMiddleWare {
	return &authMiddleWare{usecase}
}

func (m *authMiddleWare) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("sessionToken")
		if err != nil {
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}

		userLogin, err := m.usecase.Auth(r.Context(), c.Value)
		if err != nil {
			http.Error(w, string(errors.ErrNotAuthenticated), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userLogin", userLogin)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
