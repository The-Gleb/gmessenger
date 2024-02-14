package v1

import (
	"context"
	"log/slog"
	"net/http"

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

	// r.Handle(dialogURL, http.HandlerFunc(h.OpenDialog))
	// r.Handle(dialogURL, h.middlewares[0](http.HandlerFunc(h.OpenDialog)))
	slog.Debug("dialog handlers middlewares", "here they are", h.middlewares)
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
	login, ok := r.Context().Value("userLogin").(string)
	if !ok {
		slog.Error("cannot get userLogin")
		client.CloseWSConnection(conn, websocket.CloseInternalServerErr)
		return
	}
	token, ok := r.Context().Value("token").(string)
	if !ok {
		slog.Error("[dialogWSHandler.ServeHTTP]: couldn't get session token from context")
		client.CloseWSConnection(conn, websocket.CloseInternalServerErr)
		return
	}

	usecaseDTO := dialogws_usecase.OpenDialogDTO{
		Websocket:     conn,
		SenderLogin:   login,
		ReceiverLogin: chi.URLParam(r, "receiverLogin"),
		SenderToken:   token,
	}

	slog.Debug("dialog usecase dto ", "struct", usecaseDTO)

	_ = h.usecase.OpenDialog(r.Context(), usecaseDTO)

}
