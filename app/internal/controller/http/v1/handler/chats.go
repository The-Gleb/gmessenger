package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/go-chi/chi/v5"
)

const (
	chatsURL = "/chats"
)

type ChatsUsecase interface {
	ShowChats(ctx context.Context, userID int64) ([]entity.Chat, error)
}

type chatsHandler struct {
	usecase     ChatsUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewChatsHandler(usecase ChatsUsecase) *chatsHandler {
	return &chatsHandler{
		usecase:     usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
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

var testChats []entity.Chat = []entity.Chat{
	{
		Type:         "dialog",
		ReceiverID:   "2",
		ReceiverName: "John",
		Unread:       10,
		LastMessage: entity.Message{
			ID:         1,
			SenderID:   11,
			SenderName: "Gleb",
			ReceiverID: 12,
			Text:       "What`s up?",
			Timestamp:  time.Now(),
			Status:     entity.DELIVERED,
		},
	},
	{
		Type:         "group",
		ReceiverID:   "123",
		ReceiverName: "PUNK VEGANS",
		Unread:       0,
		LastMessage: entity.Message{
			ID:         125,
			SenderID:   123,
			SenderName: "Anya Larina",
			ReceiverID: 1,
			Text:       "Soe meat?",
			Timestamp:  time.Now(),
			Status:     entity.READ,
		},
	}, {
		Type:         "group",
		ReceiverID:   "123",
		ReceiverName: "Some chat",
	},
}

func (h *chatsHandler) ShowChats(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("in chat handler")

	userID, ok := r.Context().Value(v1.Key("userID")).(int64)
	if !ok {
		http.Error(rw, "cannot get userID from context ", http.StatusInternalServerError)
		return
	}

	chats, err := h.usecase.ShowChats(r.Context(), userID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO:
		return
	}

	err = json.NewEncoder(rw).Encode(chats)
	if err != nil {
		http.Error(rw, "error encoding to json", http.StatusInternalServerError)
		return
	}

}
