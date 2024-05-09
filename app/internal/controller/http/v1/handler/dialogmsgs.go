package v1

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	dialogmsgs_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/dialogmsgs"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
)

const (
	dialogMsgsURL = "/dialog/{receiverID}"
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

var testMessages = []entity.Message{
	{
		ID:         1,
		SenderID:   1,
		SenderName: "Gleb",
		ReceiverID: 2,
		Text:       "Hello John!",
		Timestamp:  time.Now(),
	},
	{
		ID:         2,
		SenderID:   1,
		SenderName: "Gleb",
		ReceiverID: 2,
		Text:       "How are you!",
		Timestamp:  time.Now(),
	},
	{
		ID:         3,
		SenderID:   2,
		SenderName: "John",
		ReceiverID: 1,
		Text:       "Hello Gleb!",
		Timestamp:  time.Now(),
	}, {
		ID:         4,
		SenderID:   2,
		SenderName: "John",
		ReceiverID: 1,
		Text:       "Im good!",
		Timestamp:  time.Now(),
	},
}

func (h *dialogMsgsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("dialogMsgsHandler working")

	userID, ok := r.Context().Value(v1.Key("userID")).(int64)
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

	receiverID, err := strconv.ParseInt(chi.URLParam(r, "receiverID"), 10, 64)
	if err != nil {
		slog.Error("[dialogMsgsHandler.ServeHTTP]: couldn't get receiverID", "error", err)
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	usecaseDTO := dialogmsgs_usecase.GetDialogMessagesDTO{
		SenderID:   userID,
		ReceiverID: receiverID,
	}

	slog.Debug("dialogMsgs usecase dto ", "struct", usecaseDTO)

	_, err = h.usecase.GetDialogMessages(r.Context(), usecaseDTO)
	if err != nil {
		slog.Error(err.Error())

		switch errors.Code(err) {
		case errors.ErrReceiverNotFound:
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	workDir, _ := os.Getwd()

	templ := template.Must(template.ParseFiles(workDir + "/app/cmd/templates/dialog/dialog.html"))

	data := struct {
		Messages []entity.Message
		UserID   int64
	}{testMessages, userID}

	err = templ.Execute(rw, data)

	//b, err := json.Marshal(messages)
	//if err != nil {
	//	slog.Error(err.Error())
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//_, err = rw.Write(b)
	//if err != nil {
	//	http.Error(rw, " error writing to body", http.StatusInternalServerError)
	//	return
	//}

}
