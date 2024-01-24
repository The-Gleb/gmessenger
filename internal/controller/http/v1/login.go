package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/The-Gleb/gmessenger/internal/controller/http/dto"
	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	login_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/login"
	"github.com/go-chi/chi/v5"
)

const (
	loginURL = "/login"
)

type LoginUsecase interface {
	Login(ctx context.Context, usecaseDTO login_usecase.LoginDTO) (entity.Session, error)
}

type loginHandler struct {
	usecase LoginUsecase
}

func NewLoginHandler(usecase LoginUsecase) *loginHandler {
	return &loginHandler{usecase: usecase}
}

func (h *loginHandler) AddToRouter(r *chi.Mux) {
	r.Post(loginURL, h.Login)
}

func (h *loginHandler) Login(rw http.ResponseWriter, r *http.Request) {

	var d dto.LoginDTO
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		// TODO
	}

	slog.Debug("LoginDTO", "struct", d)

	s, err := h.usecase.Login(r.Context(), login_usecase.LoginDTO{
		Login:    d.Login,
		Password: d.Password,
	})
	if err != nil {
		// TODO
	}
	b, err := json.Marshal(s)
	if err != nil {
		// TODO
	}
	rw.Write(b)

}
