package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	oauthURL         = "/auth/{provider}"
	oauthCallBackURL = "/auth/{provider}/callback"
	oauthLogoutURL   = "/logout/{provider}"

	key    = "some_key"
	MaxAge = 86 * 30
	IsProd = false
)

type OAuthUsecase interface {
	OAuth(ctx context.Context, email string) (string, bool, error)
}

type oauthHandler struct {
	usecase     OAuthUsecase
	middlewares []func(http.Handler) http.Handler
}

func NewOAuthHandler(usecase OAuthUsecase) *oauthHandler {
	googleClientID := "1037112513364-ik8j0vf49fvv9rrpabpl54g5kehb89r5.apps.googleusercontent.com"
	googleClientSecret := "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"

	// yandexClientID := ""
	// yandexClientSecret := ""

	//key := securecookie.GenerateRandomKey(32)
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:8081/auth/google/callback"),
	)

	return &oauthHandler{usecase: usecase,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (h *oauthHandler) AddToRouter(r *chi.Mux) {
	// r.Get(oauthRegisterURL, r.ServeHTTP)
	// r.Get(oauthRegisterCallBackURL, oauthCallback)
	// r.Get(oauthLogoutURL, oauthLogout)
	r.Route(oauthURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.ServeHTTP)

	})

	r.Route(oauthCallBackURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.oauthCallback)

	})

	r.Route(oauthLogoutURL, func(r chi.Router) {
		r.Use(h.middlewares...)
		r.Get("/", h.oauthLogout)

	})

}

func (h *oauthHandler) Middlewares(md ...func(http.Handler) http.Handler) *oauthHandler {
	h.middlewares = append(h.middlewares, md...)
	return h
}

func (h *oauthHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	slog.Info("oauth handler working!", "provider", provider)

	if gothUser, err := gothic.CompleteUserAuth(rw, r); err == nil {
		slog.Info("complete auth", "user", gothUser)

		//t, _ := template.New("foo").Parse(userTemplate)
		//t.Execute(rw, gothUser)
	} else {
		slog.Debug("begin auth")
		gothic.BeginAuthHandler(rw, r)
	}

}

func (h *oauthHandler) oauthCallback(rw http.ResponseWriter, r *http.Request) {
	slog.Debug("oauth callback handler")
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(rw, r)
	if err != nil {
		slog.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("complete auth in callback handler", "email", user.Email, "userID", user.UserID)
	// t, _ := template.New("foo").Parse(userTemplate)
	// t.Execute(rw, user)

	token, _, err := h.usecase.OAuth(r.Context(), user.Email)
	if err != nil {
		slog.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(struct {
		Token string `json:"token"`
	}{token})

	if err != nil {
		slog.Error("[handler.Register]: error encoding json into body", "error", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *oauthHandler) oauthLogout(rw http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	err := gothic.Logout(rw, r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Location", "/static/login/login.html")
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
