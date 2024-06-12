package main

import (
	"context"
	"fmt"
	storage "github.com/The-Gleb/gmessenger/internal/gateway/adapter/db/postgresql"
	"github.com/The-Gleb/gmessenger/internal/gateway/adapter/yandexgpt"
	"github.com/The-Gleb/gmessenger/internal/gateway/config"
	handlers "github.com/The-Gleb/gmessenger/internal/gateway/controller/http/v1/handler"
	middlewares "github.com/The-Gleb/gmessenger/internal/gateway/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/service"
	auth_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/auth"
	chats_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/chats"
	"github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/create_group"
	dialogmsgs_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/dialogmsgs"
	dialogws_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/dialogws.go"
	groupmsgs_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/groupmsgs"
	groupws_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/groupws"
	login_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/login"
	oauth_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/oauth"
	register_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/register"
	username_usecase "github.com/The-Gleb/gmessenger/internal/gateway/domain/usecase/username"
	"github.com/The-Gleb/gmessenger/internal/logger"
	"github.com/The-Gleb/gmessenger/pkg/client/postgresql"
	"github.com/The-Gleb/gmessenger/pkg/proto/group"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
)

var (
	catalogID = "b1glna78gcervi580uvp"
	keyId     = "ajeig6h7umsuu4q5ljte"
	apiKey    = "AQVN0oGif1ue8TRaz65Zacw4LcQDX67iH6rrlVbo"
)

func main() {
	cfg := config.MustBuild("config")
	logger.Initialize(cfg.LogLevel)

	slog.Info("config build", "config", cfg)

	dbClient, err := postgresql.NewClient(context.Background(), 3, cfg.DB)
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

	yandexGptClient := yandexgpt.NewYandexGPTClient(catalogID, apiKey)

	userService := service.NewUserService(userStorage)
	sessionService := service.NewSessionService(sessionStorage, time.Duration(24)*time.Hour)
	messageService := service.NewMessageService(messageStorage)
	dialogWSService := service.NewDialogService(messageStorage)
	pasetoAuthService, err := service.NewPaseto(make([]byte, 32), time.Duration(24)*time.Hour)
	groupHub := service.NewGroupHub(groupClient)
	otpService := service.NewOtpService(30 * time.Hour)

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
	createGroupUsecase := create_group.NewCreateGroupUsecase(groupClient)

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
	otpHandler := handlers.NewOtpHandler(otpService)
	creatGroupHandler := handlers.NewCreateGroupHandler(createGroupUsecase)
	yandexGPTHandler := handlers.NewYandexGTPHandler(yandexGptClient)

	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "app/gateway/templates"))
	FileServer(r, "/static/", filesDir)

	th := &templateHandler{fileName: "index.html"}
	r.Get("/", th.ServeHTTP)

	loginHandler.Middlewares(corsMiddleware.AllowCors).AddToRouter(r)
	registerHandler.Middlewares(corsMiddleware.AllowCors).AddToRouter(r)
	chatsHandler.Middlewares(corsMiddleware.AllowCors, authMiddleWare.Http).AddToRouter(r)
	dialogMsgsHandler.Middlewares(corsMiddleware.AllowCors, authMiddleWare.Http).AddToRouter(r)
	dialogWSHandler.Middlewares(authMiddleWare.Websocket).AddToRouter(r)
	groupMsgsHandler.Middlewares(corsMiddleware.AllowCors, authMiddleWare.Http).AddToRouter(r)
	groupWSHandler.Middlewares(authMiddleWare.Websocket, corsMiddleware.AllowCors).AddToRouter(r)
	setUsernameHandler.Middlewares(corsMiddleware.AllowCors, authMiddleWare.Http).AddToRouter(r)
	userInfoHandler.Middlewares(corsMiddleware.AllowCors, authMiddleWare.Http).AddToRouter(r)
	otpHandler.Middlewares(corsMiddleware.AllowCors, authMiddleWare.Http).AddToRouter(r)
	creatGroupHandler.Middlewares(corsMiddleware.AllowCors, authMiddleWare.Http).AddToRouter(r)
	yandexGPTHandler.AddToRouter(r)
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
	// filesDir := http.Dir(filepath.Join(workDir, "app/gateway/templates"))

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
		t.templ = template.Must(template.ParseFiles(workDir + "/app/gateway/templates/login/login.html"))
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
