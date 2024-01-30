package v1

import (
	"context"
	"encoding/json"
	"github.com/The-Gleb/gmessenger/internal/controller/http/dto"
	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	login_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/login"
	"github.com/The-Gleb/gmessenger/internal/errors"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

const (
	loginURL = "/login"
)

type LoginUsecase interface {
	Login(ctx context.Context, usecaseDTO login_usecase.LoginDTO) (entity.Session, error)
}

type loginHandler struct {
	middlewares []func(http.Handler) http.Handler
	usecase     LoginUsecase
}

func NewLoginHandler(usecase LoginUsecase) *loginHandler {
	return &loginHandler{usecase: usecase, middlewares: make([]func(http.Handler) http.Handler, 0)}
}

func (h *loginHandler) AddToRouter(r *chi.Mux) {

	r.Route(loginURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.Login)
		// r.Route(registerURL, func(r chi.Router) {

		// })
	})
}

func (h *loginHandler) Middlewares(md ...func(http.Handler) http.Handler) *loginHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *loginHandler) Login(rw http.ResponseWriter, r *http.Request) {

	var d dto.LoginDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		slog.Error("[handler.Register]: error parsing json to dto", "error", err)
		http.Error(rw, "[handler.Register]: error parsing json to dto", http.StatusBadRequest)
		return
	}

	slog.Debug("LoginDTO", "struct", d)

	s, err := h.usecase.Login(r.Context(), login_usecase.LoginDTO{
		Login:    d.Login,
		Password: d.Password,
	})
	if err != nil {
		switch errors.Code(err) {
		case errors.ErrUCWrongLoginOrPassword:
			http.Error(rw, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	c := http.Cookie{
		Name:    "sessionToken",
		Value:   s.Token,
		Expires: s.Expiry,
	}
	slog.Debug("set cookie ", "value", c.Value)
	b, err := json.Marshal(c)
	if err != nil {

	}

	rw.Write(b)

	http.SetCookie(rw, &c)

}
