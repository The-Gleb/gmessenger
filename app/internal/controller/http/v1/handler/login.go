package v1

import (
	"context"
	"encoding/json"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

const (
	loginURL = "/login"
)

type LoginUsecase interface {
	Login(ctx context.Context, dto entity.LoginDTO) (entity.Session, error)
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
	})
}

func (h *loginHandler) Middlewares(md ...func(http.Handler) http.Handler) *loginHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *loginHandler) Login(rw http.ResponseWriter, r *http.Request) {

	var dto entity.LoginDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("[loginHandler.Login]: error parsing json to dto", "error", err)
		http.Error(rw, "[loginHandler.Login]: error parsing json to dto", http.StatusBadRequest)
		return
	}
	if dto.Email == "" || dto.Password == "" {
		http.Error(rw, "[loginHandler.Login]:", http.StatusBadRequest)
		return
	}

	slog.Debug("LoginDTO", "struct", dto)

	s, err := h.usecase.Login(r.Context(), dto)
	if err != nil {
		slog.Error(err.Error())

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

	http.SetCookie(rw, &c)

	rw.WriteHeader(http.StatusOK)

	// b, err := json.Marshal(c)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	http.Error(rw, "[loginHandler.Login]: error parsing to json", http.StatusInternalServerError)
	// 	return
	// }

	// _, err = rw.Write(b)
	// if err != nil {
	// 	http.Error(rw, " error writing to body", http.StatusInternalServerError)
	// 	return
	// }

}
