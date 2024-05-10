package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/go-chi/chi/v5"
)

const (
	userInfoURL = "/me"
)

type UserInfoUsecase interface {
	GetUserInfoByID(ctx context.Context, userID int64) (entity.UserInfo, error)
}

type userInfoHandler struct {
	usecase     UserInfoUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewUserInfoHandler(usecase UserInfoUsecase) *userInfoHandler {
	return &userInfoHandler{usecase: usecase}
}

func (h *userInfoHandler) AddToRouter(r *chi.Mux) {
	r.Route(userInfoURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ServeHTTP)
	})

}

func (h *userInfoHandler) Middlewares(md ...func(http.Handler) http.Handler) *userInfoHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *userInfoHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("in chat handler")

	userID, ok := r.Context().Value(v1.Key("userID")).(int64)
	if !ok {
		http.Error(rw, "cannot get userID from context ", http.StatusInternalServerError)
		return
	}

	userInfo, err := h.usecase.GetUserInfoByID(r.Context(), userID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO:
		return
	}

	err = json.NewEncoder(rw).Encode(userInfo)
	if err != nil {
		http.Error(rw, "error encoding to json", http.StatusInternalServerError)
		return
	}

}
