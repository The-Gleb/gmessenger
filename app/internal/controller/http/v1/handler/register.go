package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/gmessenger/app/internal/controller/http/dto"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	register_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/register"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
)

const (
	registerURL = "/register"
)

type RegisterUsecase interface {
	Register(ctx context.Context, usecaseDTO register_usecase.RegisterUserDTO) (entity.Session, error)
}

type registerHandler struct {
	usecase     RegisterUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewRegisterHandler(usecase RegisterUsecase) *registerHandler {
	return &registerHandler{usecase: usecase}
}

func (h *registerHandler) AddToRouter(r *chi.Mux) {
	r.Route(registerURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Post("/", h.Register)
		// r.Route(registerURL, func(r chi.Router) {

		// })
	})

}

func (h *registerHandler) Middlewares(md ...func(http.Handler) http.Handler) *registerHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *registerHandler) Register(rw http.ResponseWriter, r *http.Request) {

	var d dto.RegisterUserDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		slog.Error("[registerHandler.Register]: error parsing json to dto", "error", err)
		http.Error(rw, "[registerHandler.Register]: error parsing json to dto", http.StatusBadRequest)
		return
	}

	if d.Login == "" || d.Password == "" || d.UserName == "" {
		http.Error(rw, "[loginHandler.Login]:", http.StatusBadRequest)
		return
	}

	slog.Debug("RegisterUserDTO", "struct", d)

	s, err := h.usecase.Register(r.Context(), register_usecase.RegisterUserDTO{
		UserName: d.UserName,
		Login:    d.Login,
		Password: d.Password,
	})
	if err != nil {
		slog.Error(err.Error())

		switch errors.Code(err) {
		case errors.ErrDBLoginAlredyExists:
			http.Error(rw, err.Error(), http.StatusConflict)
			return
		default:
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	c := http.Cookie{
		Name:    "sessionToken",
		Value:   s.Token,
		Expires: s.Expiry,
	}

	http.SetCookie(rw, &c)

	// b, err := json.Marshal(s)
	// if err != nil {
	// 	slog.Error("[handler.Register]: error unmarshalling json", "error", err)
	// 	http.Error(rw, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// rw.Write(b)

}
