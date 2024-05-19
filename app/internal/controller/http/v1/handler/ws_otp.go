package v1

import (
	"encoding/json"
	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

const (
	otpURL = "/otp"
)

type otpService interface {
	GenerateOtp(userID, sessionID int64) string
}

type otpHandler struct {
	otpService  otpService
	middlewares []func(http.Handler) http.Handler
}

func NewOtpHandler(otpService otpService) *otpHandler {
	return &otpHandler{otpService: otpService,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *otpHandler) AddToRouter(r *chi.Mux) {
	r.Route(otpURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ServeHTTP)
	})

}

func (h *otpHandler) Middlewares(md ...func(http.Handler) http.Handler) *otpHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (handler *otpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(v1.Key("userID")).(int64)
	if !ok {
		http.Error(w, "cannot get userID from context ", http.StatusInternalServerError)
		return
	}
	sessionID, ok := r.Context().Value(v1.Key("sessionID")).(int64)
	if !ok {
		http.Error(w, "cannot get sessionID from context ", http.StatusInternalServerError)
		return
	}

	otp := handler.otpService.GenerateOtp(userID, sessionID)

	err := json.NewEncoder(w).Encode(struct {
		Token string `json:"ws_otp"`
	}{otp})

	if err != nil {
		slog.Error("[handler.Otp]: error encoding json into body", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
