package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	groupmsgs_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/groupmsgs"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
)

const (
	groupMsgsURL = "/group/{groupId}"
)

type GroupMsgsUsecase interface {
	GetGroupMessages(ctx context.Context, dto groupmsgs_usecase.GetGroupMessagesDTO) ([]entity.Message, error)
}

type groupMsgsHandler struct {
	middlewares []func(http.Handler) http.Handler
	usecase     GroupMsgsUsecase
}

func NewGroupMsgsHandler(usecase GroupMsgsUsecase) *groupMsgsHandler {
	return &groupMsgsHandler{usecase: usecase, middlewares: make([]func(http.Handler) http.Handler, 0)}
}

func (h *groupMsgsHandler) AddToRouter(r *chi.Mux) {

	r.Route(groupMsgsURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ServeHTTP)
	})

	// var handler http.Handler
	// handler = h
	// for _, md := range h.middlewares {
	// 	handler = md(h)
	// }
	// r.Handle(groupMsgsURL, handler)
}

func (h *groupMsgsHandler) Middlewares(md ...func(http.Handler) http.Handler) *groupMsgsHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *groupMsgsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("groupMsgsHandler working")

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

	groupId, err := strconv.ParseInt(chi.URLParam(r, "groupId"), 10, 64)
	if err != nil {
		slog.Error("couldn`t parse to int group id from URL param")
	}

	usecaseDTO := groupmsgs_usecase.GetGroupMessagesDTO{
		UserLogin: login,
		GroupID:   groupId,
	}

	slog.Debug("groupMsgs usecase dto ", "struct", usecaseDTO)

	messages, err := h.usecase.GetGroupMessages(r.Context(), usecaseDTO)
	if err != nil {
		slog.Error(err.Error())

		switch errors.Code(err) {

		case errors.ErrNotAMember:
			http.Error(rw, err.Error(), http.StatusForbidden)
			return

		case errors.ErrGroupNotFound:
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
	slog.Debug("got messages", "messages", messages)

	rw.Write(b)

}
