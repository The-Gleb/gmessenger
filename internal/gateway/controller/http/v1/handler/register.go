package v1

import (
	"context"
	"encoding/json"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/gateway/errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	registerURL = "/register"
)

type RegisterUsecase interface {
	Register(ctx context.Context, dto entity.RegisterUserDTO) (string, error)
}

type registerHandler struct {
	usecase     RegisterUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewRegisterHandler(usecase RegisterUsecase) *registerHandler {
	return &registerHandler{usecase: usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *registerHandler) AddToRouter(r *chi.Mux) {
	r.Route(registerURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.Register)
		// r.Route(registerURL, func(r chi.Router) {

		// })
	})

}

func (h *registerHandler) Middlewares(md ...func(http.Handler) http.Handler) *registerHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *registerHandler) Register(rw http.ResponseWriter, r *http.Request) {

	var dto entity.RegisterUserDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("[setUsernameHandler.Register]: error parsing json to dto", "error", err)
		http.Error(rw, "[setUsernameHandler.Register]: error parsing json to dto", http.StatusBadRequest)
		return
	}

	if dto.Email == "" || dto.Password == "" || dto.Username == "" {
		http.Error(rw, "[loginHandler.Login]:", http.StatusBadRequest)
		return
	}

	slog.Debug("RegisterUserDTO", "struct", dto)

	// TODO: verify email

	// for now assume that email is verified

	token, err := h.usecase.Register(r.Context(), dto)
	if err != nil {
		slog.Error(err.Error())

		switch errors.Code(err) {
		case errors.ErrDBLoginAlredyExists, errors.ErrUserExists:
			http.Error(rw, err.Error(), http.StatusConflict)
			return
		default:
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = json.NewEncoder(rw).Encode(struct {
		Token string `json:"token"`
	}{token})

	if err != nil {
		slog.Error("[handler.Register]: error encoding json into body", "error", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	//rw.WriteHeader(http.StatusOK)

}
