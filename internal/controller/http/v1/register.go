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
	r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		// TODO
	}

}
