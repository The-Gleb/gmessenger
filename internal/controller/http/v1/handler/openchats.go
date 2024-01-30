package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/gmessenger/internal/controller/http/dto"
	openchat_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/openchat.go"
	"github.com/go-chi/chi/v5"
)

const (
	openChatURL = "/chats/{chatID}"
)

type OpenChatUsecase interface {
	OpenChat(ctx context.Context, dto openchat_usecase.OpenChatDTO) error
}

type openChatHandler struct {
	middlewares []func(http.Handler) http.Handler
	usecase     OpenChatUsecase
}

func NewOpenChatHandler(usecase OpenChatUsecase) *openChatHandler {
	return &openChatHandler{usecase: usecase, middlewares: make([]func(http.Handler) http.Handler, 0)}
}

func (h *openChatHandler) AddToRouter(r *chi.Mux) {

	r.Route(openChatURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.OpenChat)
		// r.Route(registerURL, func(r chi.Router) {

		// })
	})
}

func (h *openChatHandler) Middlewares(md ...func(http.Handler) http.Handler) *openChatHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *openChatHandler) OpenChat(rw http.ResponseWriter, r *http.Request) {

	var d dto.OpenChatDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		slog.Error("[handler.OpenChat]: error parsing json to dto", "error", err)
		http.Error(rw, "[handler.OpenChat]: error parsing json to dto", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	token, ok := r.Context().Value("token").(string)
	if !ok {
		http.Error(rw, "[handler.OpenChat]: couldn't get session token from context", http.StatusInternalServerError)
	}

	usecaseDTO := openchat_usecase.OpenChatDTO{
		ChatType:    d.ChatType,
		ChatID:      d.ChatID,
		SenderLogin: d.SenderLogin,
		Websocket:   conn,
	}
	if usecaseDTO.ChatType == "personal" {
		usecaseDTO.ReceiverLogin = d.ReceiverLogin
		usecaseDTO.SenderToken = token
	}

	err = h.usecase.OpenChat(r.Context(), usecaseDTO)
	if err != nil {

	}

}
