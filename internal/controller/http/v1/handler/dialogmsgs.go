package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	dialogmsgs_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/dialogmsgs"
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

	// r.Route(dialogURL, func(r chi.Router) {

	// 	r.Use(h.middlewares...)

	// })

	var handler http.Handler
	handler = h
	for _, md := range h.middlewares {
		handler = md(h)
	}

	r.Handle(dialogMsgsURL, handler)

	// r.Handle(dialogURL, http.HandlerFunc(h.OpenDialog))
	// r.Handle(dialogURL, h.middlewares[0](http.HandlerFunc(h.OpenDialog)))
}

func (h *dialogMsgsHandler) Middlewares(md ...func(http.Handler) http.Handler) *dialogMsgsHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *dialogMsgsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("dialogMsgsHandler working")

	login, ok := r.Context().Value("userLogin").(string)
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
		slog.Error(err.Error()) // TODO
	}

	b, err := json.Marshal(messages)
	if err != nil {
		slog.Error(err.Error()) // TODO
	}
	slog.Debug("got messages", "messages", messages)

	rw.Write(b)

}
