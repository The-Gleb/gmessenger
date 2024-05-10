package v1

import (
	"context"
	"encoding/json"
	middlewares "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/go-chi/chi/v5"
)

const (
	setUsernameURL = "/set_username"
)

type username interface {
	SetUsername(ctx context.Context, dto entity.SetUsernameDTO) error
}

type setUsernameHandler struct {
	usecase     username
	middlewares []func(http.Handler) http.Handler
}

func NewSetUsernameHandler(usecase username) *setUsernameHandler {
	return &setUsernameHandler{usecase: usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *setUsernameHandler) AddToRouter(r *chi.Mux) {
	r.Route(setUsernameURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.ServeHTTP)
		// r.Route(registerURL, func(r chi.Router) {

		// })
	})

}

func (h *setUsernameHandler) Middlewares(md ...func(http.Handler) http.Handler) *setUsernameHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *setUsernameHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	var dto entity.SetUsernameDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("[setUsernameHandler]: error parsing json to dto", "error", err)
		http.Error(rw, "[setUsernameHandler]: error parsing json to dto", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middlewares.Key("userID")).(int64)
	if !ok {
		slog.Error("user id not found in context", "userID", userID)
		http.Error(rw, "[setUsernameHandler]: user id not found in context", http.StatusInternalServerError)
		return
	}
	dto.UserID = userID

	if dto.Username == "" || dto.UserID == 0 {
		slog.Error("invalid request body", "body", dto)
		http.Error(rw, "[setUsernameHandler]:", http.StatusBadRequest)
		return
	}

	err = h.usecase.SetUsername(r.Context(), dto)
	if err != nil {
		slog.Error("[setUsernameHandler]: error setting username", "error", err)
		http.Error(rw, "[setUsernameHandler]: error setting username", http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, chatsURL, http.StatusFound)

}
