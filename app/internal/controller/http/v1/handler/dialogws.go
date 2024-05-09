package v1

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	dialogws_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/dialogws.go"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

const (
	dialogWSURL = "/dialog/{receiverLogin}/ws"
)

type DialogWSUsecase interface {
	OpenDialog(ctx context.Context, dto dialogws_usecase.OpenDialogDTO) error
}

type dialogWSHandler struct {
	middlewares []func(http.Handler) http.Handler
	usecase     DialogWSUsecase
}

func NewDialogWSHandler(usecase DialogWSUsecase) *dialogWSHandler {
	return &dialogWSHandler{usecase: usecase, middlewares: make([]func(http.Handler) http.Handler, 0)}
}

func (h *dialogWSHandler) AddToRouter(r *chi.Mux) {

	// r.Route(dialogURL, func(r chi.Router) {

	// 	r.Use(h.middlewares...)

	// })

	var handler http.Handler
	handler = h
	for _, md := range h.middlewares {
		handler = md(h)
	}

	r.Handle(dialogWSURL, handler)
}

func (h *dialogWSHandler) Middlewares(md ...func(http.Handler) http.Handler) *dialogWSHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *dialogWSHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("im working")

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, "[dialogWSHandler.ServeHTTP]: couldn't get session token from context", http.StatusInternalServerError)
		return
	}
	userID, ok := r.Context().Value(v1.Key("userID")).(int64)
	if !ok {
		slog.Error("cannot get userLogin")
		client.CloseWSConnection(conn, websocket.CloseInternalServerErr)
		return
	}
	token, ok := r.Context().Value(v1.Key("token")).(string)
	if !ok {
		slog.Error("[dialogWSHandler.ServeHTTP]: couldn't get session token from context")
		client.CloseWSConnection(conn, websocket.CloseInternalServerErr)
		return
	}

	receiverID, err := strconv.ParseInt(chi.URLParam(r, "receiverID"), 10, 64)
	if err != nil {
		slog.Error("[dialogWSHandler.ServeHTTP]: couldn't get receiverID", "error", err)
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	usecaseDTO := dialogws_usecase.OpenDialogDTO{
		Websocket:   conn,
		SenderID:    userID,
		ReceiverID:  receiverID,
		SenderToken: token,
	}

	slog.Debug("dialog usecase dto ", "struct", usecaseDTO)

	_ = h.usecase.OpenDialog(r.Context(), usecaseDTO)

}
