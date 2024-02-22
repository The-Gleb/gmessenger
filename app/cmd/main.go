package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	storage "github.com/The-Gleb/gmessenger/app/internal/adapter/db/postgresql"
	"github.com/The-Gleb/gmessenger/app/internal/config"
	handlers "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/handler"
	middlewares "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/pkg/proto/go/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/The-Gleb/gmessenger/app/internal/domain/service"
	auth_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/auth"
	chats_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/chats"
	dialogmsgs_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/dialogmsgs"
	dialogws_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/dialogws.go"
	groupmsgs_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/groupmsgs"
	groupws_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/groupws"
	login_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/login"
	register_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/register"
	"github.com/The-Gleb/gmessenger/app/internal/logger"
	db "github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustBuild("config")
	logger.Initialize(cfg.LogLevel)

	slog.Info("config build", "config", cfg)

	dbClient, err := db.NewClient(context.Background(), 3, cfg.DB)
	if err != nil {
		panic(err)
	}

	groupServerAddr := fmt.Sprintf("%s:%d", cfg.GroupServerHost, cfg.GroupServerPort)
	// TODO: make secure
	conn, err := grpc.Dial(groupServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	slog.Debug("chech grpc conn to group server", "conn.State", conn.GetState())
	groupClient := group.NewGroupClient(conn)

	userStorage := storage.NewUserStorage(dbClient)
	sessionStorage := storage.NewSessionStorage(dbClient)
	messageStorage := storage.NewMessageStorage(dbClient)

	userService := service.NewUserService(userStorage)
	sessionService := service.NewSessionService(sessionStorage)
	messageService := service.NewMessageService(messageStorage)
	dialogWSService := service.NewDialogService(messageStorage)
	groupHub := service.NewGroupHub(groupClient)

	loginUsecase := login_usecase.NewLoginUsecase(userService, sessionService)
	registerUsecase := register_usecase.NewRegisterUsecase(userService, sessionService)
	authUsecase := auth_usecase.NewAuthUsecase(sessionService)
	chatsUsecase := chats_usecase.NewChatsUsecase(userService, groupClient, messageService)
	dialogWSUsecase := dialogws_usecase.NewDialogWSUsecase(dialogWSService)
	dialogMsgsUsecase := dialogmsgs_usecase.NewDialogMsgsUsecase(messageService)
	groupWSUsecase := groupws_usecase.NewGroupWSUsecase(groupHub)
	groupMsgsUsecase := groupmsgs_usecase.NewGroupMsgsUsecase(groupClient)

	authMiddleWare := middlewares.NewAuthMiddleware(authUsecase)
	loginHandler := handlers.NewLoginHandler(loginUsecase)
	registerHandler := handlers.NewRegisterHandler(registerUsecase)
	chatsHandler := handlers.NewChatsHandler(chatsUsecase)
	dialogWSHandler := handlers.NewDialogWSHandler(dialogWSUsecase)
	dialogMsgsHandler := handlers.NewDialogMsgsHandler(dialogMsgsUsecase)
	groupWSHandler := handlers.NewGroupWSHandler(groupWSUsecase)
	groupMsgsHandler := handlers.NewGroupMsgsHandler(groupMsgsUsecase)

	r := chi.NewRouter()

	th := &templateHandler{fileName: "index.html"}
	r.Get("/", th.ServeHTTP)
	// h := func(w http.ResponseWriter, r *http.Request) {
	// 	templ := template.Must(template.ParseFiles("./cmd/index.html"))

	// 	err = templ.Execute(w, nil)
	// 	if err != nil {
	// 		slog.Error(err.Error())
	// 	}
	// }
	// http.Handle("/", http.HandlerFunc(h))
	// err = http.ListenAndServe(":8081", nil)
	// if err != nil {
	// 	slog.Error(err.Error())
	// }

	loginHandler.AddToRouter(r)
	registerHandler.AddToRouter(r)
	chatsHandler.Middlewares(authMiddleWare.Http).AddToRouter(r)
	dialogMsgsHandler.Middlewares(authMiddleWare.Http).AddToRouter(r)
	dialogWSHandler.Middlewares(authMiddleWare.Websocket).AddToRouter(r)
	groupMsgsHandler.Middlewares(authMiddleWare.Http).AddToRouter(r)
	groupWSHandler.Middlewares(authMiddleWare.Websocket).AddToRouter(r)

	s := http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ServerShutdownSignal := make(chan os.Signal, 1)
		signal.Notify(ServerShutdownSignal, syscall.SIGINT)
		<-ServerShutdownSignal
		err = s.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
		slog.Info("server shutdown")
	}()

	slog.Info("starting server")
	if err := s.ListenAndServe(); err != nil {
		slog.Error("server error", "error", err)
	}

}

type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug("requst pointer", "pointer", r)
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles("./app/cmd/index.html"))
	})

	err := t.templ.Execute(w, nil)
	if err != nil {
		slog.Error(err.Error())
	}
}
