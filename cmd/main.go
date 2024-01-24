package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	storage "github.com/The-Gleb/gmessenger/internal/adapter/db/postgresql"
	"github.com/The-Gleb/gmessenger/internal/config"
	handlers "github.com/The-Gleb/gmessenger/internal/controller/http/v1/handler"
	middlewares "github.com/The-Gleb/gmessenger/internal/controller/http/v1/middleware"

	"github.com/The-Gleb/gmessenger/internal/domain/service"
	auth_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/auth"
	chats_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/chats"
	login_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/login"
	register_usecase "github.com/The-Gleb/gmessenger/internal/domain/usecase/register"
	"github.com/The-Gleb/gmessenger/internal/logger"
	db "github.com/The-Gleb/gmessenger/pkg/client/postgresql"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.BuildConfig("config")
	logger.Initialize(cfg.LogLevel)

	slog.Info("config build", "config", cfg)

	dbClient, err := db.NewClient(context.Background(), 3, cfg.DB)
	if err != nil {
		panic(err)
	}

	userStorage := storage.NewUserStorage(dbClient)
	userService := service.NewUserService(userStorage)

	sessionStorage := storage.NewSessionStorage(dbClient)
	sessionService := service.NewSessionService(sessionStorage)

	loginUsecase := login_usecase.NewLoginUsecase(userService, sessionService)
	registerUsecase := register_usecase.NewRegisterUsecase(userService, sessionService)
	authUsecase := auth_usecase.NewAuthUsecase(sessionService)
	chatsUsecase := chats_usecase.NewChatsUsecase(userService, sessionService)

	authMiddleWare := middlewares.NewAuthMiddleware(authUsecase)
	loginHandler := handlers.NewLoginHandler(loginUsecase)
	registerHandler := handlers.NewRegisterHandler(registerUsecase)
	chatsHandler := handlers.NewChatsHandler(chatsUsecase)

	r := chi.NewRouter()

	loginHandler.AddToRouter(r)
	registerHandler.AddToRouter(r)
	chatsHandler.Middlewares(authMiddleWare.Do).AddToRouter(r)

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
		s.Shutdown(context.Background())
		slog.Info("server shutdown")
	}()

	slog.Info("starting server")
	if err := s.ListenAndServe(); err != nil {
		slog.Error("server error", "error", err)
	}

}

// type templateHandler struct {
// 	once     sync.Once
// 	fileName string
// 	templ    *template.Template
// }

// func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	t.once.Do(func() {
// 		t.templ = template.Must(template.ParseFiles((filepath.Join("templates", t.fileName))))
// 	})
// 	t.templ.Execute(w, r)
// }
// func main() {

// 	var addr = flag.String("a", ":8080", "run address")

// 	flag.Parse()

// 	r := newRoom()

// 	http.Handle("/", &templateHandler{fileName: "chat.html"})
// 	http.Handle("/room", r)

// 	go r.run()

// 	if err := http.ListenAndServe(*addr, nil); err != nil {
// 		log.Fatal(err)
// 	}
// }
