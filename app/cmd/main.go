package main

import (
	"context"
	"fmt"
	oauth_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/oauth"
	username_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/username"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	storage "github.com/The-Gleb/gmessenger/app/internal/adapter/db/postgresql"
	"github.com/The-Gleb/gmessenger/app/internal/config"
	handlers "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/handler"
	middlewares "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	db "github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
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

	// db "github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
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
	pasetoAuthService, err := service.NewPaseto(make([]byte, 32), time.Duration(1)*time.Hour)
	groupHub := service.NewGroupHub(groupClient)
	otpService := service.NewOtpService()

	loginUsecase := login_usecase.NewLoginUsecase(userService, pasetoAuthService, sessionService)
	registerUsecase := register_usecase.NewRegisterUsecase(userService, pasetoAuthService, sessionService)
	authUsecase := auth_usecase.NewAuthUsecase(pasetoAuthService)
	chatsUsecase := chats_usecase.NewChatsUsecase(userService, groupClient, messageService)
	dialogWSUsecase := dialogws_usecase.NewDialogWSUsecase(dialogWSService)
	dialogMsgsUsecase := dialogmsgs_usecase.NewDialogMsgsUsecase(messageService)
	groupWSUsecase := groupws_usecase.NewGroupWSUsecase(groupHub)
	groupMsgsUsecase := groupmsgs_usecase.NewGroupMsgsUsecase(groupClient)
	setUsernameUsecase := username_usecase.NewUsernameUsecase(userService)
	oauthUsecase := oauth_usecase.NewOAuthUsecase(userService, pasetoAuthService, sessionService)

	authMiddleWare := middlewares.NewAuthMiddleware(authUsecase, otpService)
	corsMiddleware := middlewares.NewCorsMiddleware()

	loginHandler := handlers.NewLoginHandler(loginUsecase)
	registerHandler := handlers.NewRegisterHandler(registerUsecase)
	chatsHandler := handlers.NewChatsHandler(chatsUsecase)
	dialogWSHandler := handlers.NewDialogWSHandler(dialogWSUsecase)
	dialogMsgsHandler := handlers.NewDialogMsgsHandler(dialogMsgsUsecase)
	groupWSHandler := handlers.NewGroupWSHandler(groupWSUsecase)
	groupMsgsHandler := handlers.NewGroupMsgsHandler(groupMsgsUsecase)
	setUsernameHandler := handlers.NewSetUsernameHandler(setUsernameUsecase)
	oauthHandler := handlers.NewOAuthHandler(oauthUsecase)
	userInfoHandler := handlers.NewUserInfoHandler(userService)

	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "app/cmd/templates"))
	FileServer(r, "/static/", filesDir)

	th := &templateHandler{fileName: "index.html"}
	r.Get("/", th.ServeHTTP)

	loginHandler.Middlewares(corsMiddleware.AllowCors).AddToRouter(r)
	registerHandler.Middlewares(corsMiddleware.AllowCors).AddToRouter(r)
	chatsHandler.Middlewares(authMiddleWare.Http, corsMiddleware.AllowCors).AddToRouter(r)
	dialogMsgsHandler.Middlewares(authMiddleWare.Http, corsMiddleware.AllowCors).AddToRouter(r)
	dialogWSHandler.Middlewares(authMiddleWare.Websocket, corsMiddleware.AllowCors).AddToRouter(r)
	groupMsgsHandler.Middlewares(authMiddleWare.Http, corsMiddleware.AllowCors).AddToRouter(r)
	groupWSHandler.Middlewares(authMiddleWare.Websocket, corsMiddleware.AllowCors).AddToRouter(r)
	setUsernameHandler.Middlewares(authMiddleWare.Http, corsMiddleware.AllowCors).AddToRouter(r)
	userInfoHandler.Middlewares(authMiddleWare.Http, corsMiddleware.AllowCors).AddToRouter(r)
	oauthHandler.AddToRouter(r)

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

	slog.Debug("root handler working")
	// filesDir := http.Dir(filepath.Join(workDir, "app/cmd/templates"))

	c, err := r.Cookie("sessionToken")
	if err == nil && c.Value != "" {
		slog.Debug("sessionToken found", "token", c)
		http.Redirect(w, r, "/chats", http.StatusMovedPermanently)
		return
	}

	slog.Debug("sessionToken not found, redirecting to login")

	http.Redirect(w, r, "/static/login/login.html", http.StatusFound)
	workDir, _ := os.Getwd()
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(workDir + "/app/cmd/templates/login/login.html"))
	})

	err = t.templ.Execute(w, nil)
	if err != nil {
		slog.Error(err.Error())
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
