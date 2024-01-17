package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/The-Gleb/gmessenger/internal/controller/http/dto"
	"github.com/The-Gleb/gmessenger/internal/domain/entity"
	register_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/register"
	"github.com/go-chi/chi/v5"
)

const (
	loginURL = "/login"
)

type LoginUsecase interface {
	login(ctx context.Context, usecaseDTO register_usecase.RegisterUserDTO) (entity.Session, error)
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
	r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		// TODO
	}

}
