package http

import "github.com/go-chi/chi"

type Handler interface {
	AddToRouter(*chi.Mux)
}
