package v1

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

const (
	yandexGptURL = "/glebgpt"
)

type YandexGPTClient interface {
	SendMessage(ctx context.Context, message string) (string, error)
}

type yandexGPTHandler struct {
	client      YandexGPTClient
	middlewares []func(http.Handler) http.Handler
}

func NewYandexGTPHandler(c YandexGPTClient) *yandexGPTHandler {
	return &yandexGPTHandler{client: c, middlewares: make([]func(http.Handler) http.Handler, 0)}
}

func (h *yandexGPTHandler) AddToRouter(r *chi.Mux) {

	r.Route(yandexGptURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.ServeHTTP)
	})

	// var handler http.Handler
	// handler = h
	// for _, md := range h.middlewares {
	// 	handler = md(h)
	// }
	// r.Handle(dialogMsgsURL, handler)
}

func (h *yandexGPTHandler) Middlewares(md ...func(http.Handler) http.Handler) *yandexGPTHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *yandexGPTHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("yandexGPTHandler working")
	var dto struct {
		Text string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error("yandexGPTHandler decode json error:", err)
		http.Error(rw, "decode json error", http.StatusBadRequest)
		return
	}

	response, err := h.client.SendMessage(r.Context(), dto.Text)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	dto.Text = response

	b, err := json.Marshal(dto)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(b)
	if err != nil {
		http.Error(rw, " error writing to body", http.StatusInternalServerError)
		return
	}

}
