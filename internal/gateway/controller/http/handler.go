package http

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Handler interface {
	AddToRouter(*chi.Mux)
	Middlewares(md ...func(http.Handler) http.Handler) *Handler
}
