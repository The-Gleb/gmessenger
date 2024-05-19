package v1

import (
	"context"
	"encoding/json"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/entity"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

const (
	createGroupURL = "/group/create"
)

type createGroupUsecase interface {
	CreateGroup(ctx context.Context, dto entity.CreateGroupDTO) (entity.Group, error)
}

type createGroupHandler struct {
	usecase     createGroupUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewCreateGroupHandler(usecase createGroupUsecase) *createGroupHandler {
	return &createGroupHandler{usecase: usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *createGroupHandler) AddToRouter(r *chi.Mux) {
	r.Route(createGroupURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ServeHTTP)
	})

}

func (h *createGroupHandler) Middlewares(md ...func(http.Handler) http.Handler) *createGroupHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (handler *createGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//userID, ok := r.Context().Value(v1.Key("userID")).(int64)
	//if !ok {
	//	http.Error(w, "cannot get userID from context ", http.StatusInternalServerError)
	//	return
	//}
	var dto entity.CreateGroupDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	group, err := handler.usecase.CreateGroup(r.Context(), dto)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		slog.Error("[handler.CreateGroup]: error encoding json into body", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
