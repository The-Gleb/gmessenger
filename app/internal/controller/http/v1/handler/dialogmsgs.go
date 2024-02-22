package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	dialogmsgs_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/dialogmsgs"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
)

const (
	dialogMsgsURL = "/dialog/{receiverLogin}"
)

type DialogMsgsUsecase interface {
	GetDialogMessages(ctx context.Context, dto dialogmsgs_usecase.GetDialogMessagesDTO) ([]entity.Message, error)
}

type dialogMsgsHandler struct {
	middlewares []func(http.Handler) http.Handler
	usecase     DialogMsgsUsecase
}

func NewDialogMsgsHandler(usecase DialogMsgsUsecase) *dialogMsgsHandler {
	return &dialogMsgsHandler{usecase: usecase, middlewares: make([]func(http.Handler) http.Handler, 0)}
}

func (h *dialogMsgsHandler) AddToRouter(r *chi.Mux) {

	r.Route(dialogMsgsURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ServeHTTP)
	})

	// var handler http.Handler
	// handler = h
	// for _, md := range h.middlewares {
	// 	handler = md(h)
	// }
	// r.Handle(dialogMsgsURL, handler)
}

func (h *dialogMsgsHandler) Middlewares(md ...func(http.Handler) http.Handler) *dialogMsgsHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *dialogMsgsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("dialogMsgsHandler working")

	login, ok := r.Context().Value(v1.Key("userLogin")).(string)
	if !ok {
		slog.Error("cannot get userLogin")
		http.Error(rw, "cannot get userLogin", http.StatusInternalServerError)
		return
	}
	// token, ok := r.Context().Value("token").(string)
	// if !ok {
	// 	slog.Error("[handler.OpenChat]: couldn't get session token from context")
	// 	http.Error(rw, "[handler.OpenChat]: couldn't get session token from context", http.StatusInternalServerError)
	// }

	usecaseDTO := dialogmsgs_usecase.GetDialogMessagesDTO{
		SenderLogin:   login,
		ReceiverLogin: chi.URLParam(r, "receiverLogin"),
	}

	slog.Debug("dialogMsgs usecase dto ", "struct", usecaseDTO)

	messages, err := h.usecase.GetDialogMessages(r.Context(), usecaseDTO)
	if err != nil {
		slog.Error(err.Error())

		switch errors.Code(err) {
		case errors.ErrReceiverNotFound:
			http.Error(rw, err.Error(), http.StatusNotFound)

		default:
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	b, err := json.Marshal(messages)
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
