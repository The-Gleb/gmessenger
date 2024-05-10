package v1

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	groupws_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/groupws"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

const (
	groupWSURL = "/group/{groupId}/ws"
)

type GroupWSUsecase interface {
	OpenGroup(ctx context.Context, dto groupws_usecase.OpenGroupDTO) error
}

type groupWSHandler struct {
	middlewares []func(http.Handler) http.Handler
	usecase     GroupWSUsecase
}

func NewGroupWSHandler(usecase GroupWSUsecase) *groupWSHandler {
	return &groupWSHandler{usecase: usecase, middlewares: make([]func(http.Handler) http.Handler, 0)}
}

func (h *groupWSHandler) AddToRouter(r *chi.Mux) {

	var handler http.Handler
	handler = h
	for _, md := range h.middlewares {
		handler = md(h)
	}

	r.Handle(groupWSURL, handler)

}

func (h *groupWSHandler) Middlewares(md ...func(http.Handler) http.Handler) *groupWSHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *groupWSHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("im group ws handler and i`m working")

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, "cannot get userID", http.StatusInternalServerError)
		return
	}
	userID, ok := r.Context().Value(v1.Key("userID")).(int64)
	if !ok {
		slog.Error("cannot get userID")
		client.CloseWSConnection(conn, websocket.CloseInternalServerErr)
		return
	}
	sessionID, ok := r.Context().Value(v1.Key("session")).(int64)
	if !ok {
		http.Error(rw, "cannot get userID from context ", http.StatusInternalServerError)
		return
	}

	groupID, err := strconv.ParseInt(chi.URLParam(r, "groupId"), 10, 64)
	if err != nil {
		slog.Error("couldn`t parse to int group id from URL param")
		client.CloseWSConnection(conn, websocket.CloseInternalServerErr)
		return
	}

	usecaseDTO := groupws_usecase.OpenGroupDTO{
		Websocket: conn,
		SenderID:  userID,
		GroupID:   groupID,
		SessionID: sessionID,
	}

	slog.Debug("group usecase dto ", "struct", usecaseDTO)

	_ = h.usecase.OpenGroup(r.Context(), usecaseDTO)

}
