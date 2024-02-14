package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/go-chi/chi/v5"
)

const (
	chatsURL = "/chats"
)

type ChatsUsecase interface {
	ShowChats(ctx context.Context, login string) ([]entity.Chat, error)
}

type chatsHandler struct {
	usecase     ChatsUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewChatsHandler(usecase ChatsUsecase) *chatsHandler {
	return &chatsHandler{usecase: usecase}
}

func (h *chatsHandler) AddToRouter(r *chi.Mux) {
	r.Route(chatsURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ShowChats)
	})

}

func (h *chatsHandler) Middlewares(md ...func(http.Handler) http.Handler) *chatsHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *chatsHandler) ShowChats(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("in chat handler")

	login, ok := r.Context().Value("userLogin").(string)
	if !ok {
		slog.Error("cannot get userLogin")
		http.Error(rw, "cannot get userLogin", http.StatusInternalServerError)
		return
	}

	chats, err := h.usecase.ShowChats(r.Context(), login)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO:
		return
	}

	b, err := json.Marshal(chats)
	if err != nil {
		http.Error(rw, "error parsing to json", http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(b)
	if err != nil {
		http.Error(rw, "error writing to body", http.StatusInternalServerError)
		return
	}

}
