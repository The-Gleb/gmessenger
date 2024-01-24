package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
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
		r.Post("/", h.ShowChats)
		// r.Route(registerURL, func(r chi.Router) {

		// })
	})

}

func (h *chatsHandler) Middlewares(md ...func(http.Handler) http.Handler) *chatsHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *chatsHandler) ShowChats(rw http.ResponseWriter, r *http.Request) {

	login, ok := r.Context().Value("userLogin").(string)
	if !ok {
		http.Error(rw, "cannot get userLogin", http.StatusInternalServerError)
		return
	}

	chats, err := h.usecase.ShowChats(r.Context(), login)
	if err != nil {

	}

	b, err := json.Marshal(chats)
	if err != nil {

	}

	rw.Write([]byte(fmt.Sprintf("chats:\n%s\n. user login:%s", string(b), login)))

	// b, err := json.Marshal(s)
	// if err != nil {
	// 	slog.Error("[handler.Register]: error unmarshalling json", "error", err)
	// 	http.Error(rw, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// rw.Write(b)

}
