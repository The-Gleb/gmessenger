package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/gmessenger/internal/controller/http/dto"
	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	register_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/register"
	"github.com/The-Gleb/gmessenger/internal/errors"
	"github.com/go-chi/chi/v5"
)

const (
	registerURL = "/register"
)

type RegisterUsecase interface {
	Register(ctx context.Context, usecaseDTO register_usecase.RegisterUserDTO) (entity.Session, error)
}

type registerHandler struct {
	usecase RegisterUsecase
}

func NewRegisterHandler(usecase RegisterUsecase) *registerHandler {
	return &registerHandler{usecase: usecase}
}

func (h *registerHandler) AddToRouter(r *chi.Mux) {
	r.Post(registerURL, h.Register)
}

func (h *registerHandler) Register(rw http.ResponseWriter, r *http.Request) {

	var d dto.RegisterUserDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		slog.Error("[handler.Register]: error parsing json to dto", "error", err)
		http.Error(rw, "[handler.Register]: error parsing json to dto", http.StatusInternalServerError)
		return
	}

	slog.Debug("RegisterUserDTO", "struct", d)

	s, err := h.usecase.Register(r.Context(), register_usecase.RegisterUserDTO{
		UserName: d.UserName,
		Login:    d.Login,
		Password: d.Password,
	})
	if err != nil {
		switch errors.Code(err) {
		case errors.ErrDBLoginAlredyExists:
			http.Error(rw, err.Error(), http.StatusConflict)
			return
		default:
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	b, err := json.Marshal(s)
	if err != nil {
		slog.Error("[handler.Register]: error unmarshalling json", "error", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Write(b)

}
