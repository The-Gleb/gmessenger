package v1

import (
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/controller/http/dto"
	"github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/handler/mocks"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	login_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/login"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_loginHandler_Login(t *testing.T) {

	validLoginReqBody, err := json.Marshal(dto.LoginDTO{
		Login:    "login1",
		Password: "password1",
	})
	require.NoError(t, err)

	invalidLoginReqBody, err := json.Marshal(dto.LoginDTO{
		Login: "login1",
	})
	require.NoError(t, err)

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockLoginUsecase := mocks.NewMockLoginUsecase(ctrl)
	loginHandler := NewLoginHandler(mockLoginUsecase)
	loginHandler.AddToRouter(r)
	server := httptest.NewServer(r)

	tests := []struct {
		name    string
		reqBody json.RawMessage
		code    int
		prepare func()
	}{
		{
			name:    "positive",
			reqBody: validLoginReqBody,
			code:    200,
			prepare: func() {
				mockLoginUsecase.
					EXPECT().
					Login(gomock.Any(), gomock.Eq(login_usecase.LoginDTO{
						Login:    "login1",
						Password: "password1",
					})).
					Return(entity.Session{
						UserLogin: "login1",
						Token:     "123",
						Expiry:    time.Now().Add(time.Hour),
					}, nil)
			},
		},
		{
			name:    "negative, body with no password",
			reqBody: invalidLoginReqBody,
			code:    400,
			prepare: func() {
			},
		},
		{
			name:    "negative, invalid body",
			reqBody: []byte("domasldkfjdf"),
			code:    400,
			prepare: func() {
			},
		},
		{
			name:    "negative, wrong login/password",
			reqBody: validLoginReqBody,
			code:    401,
			prepare: func() {
				mockLoginUsecase.
					EXPECT().
					Login(gomock.Any(), gomock.Eq(login_usecase.LoginDTO{
						Login:    "login1",
						Password: "password1",
					})).
					Return(entity.Session{}, errors.NewDomainError(errors.ErrUCWrongLoginOrPassword, ""))
			},
		},
		{
			name:    "negative, some db err",
			reqBody: validLoginReqBody,
			code:    500,
			prepare: func() {
				mockLoginUsecase.
					EXPECT().
					Login(gomock.Any(), gomock.Eq(login_usecase.LoginDTO{
						Login:    "login1",
						Password: "password1",
					})).
					Return(entity.Session{}, errors.NewDomainError(errors.ErrDB, ""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.prepare()

			resp, _ := testRequest(t, "", server, "POST", "/login", tt.reqBody)
			defer resp.Body.Close()

			require.Equal(t, tt.code, resp.StatusCode)

			if tt.code != 200 {
				return
			}

			cookies := resp.Cookies()
			require.NotEqual(t, "0", len(cookies))
			require.Equal(t, "sessionToken", cookies[0].Name)
			require.NotEmpty(t, cookies[0].Value)
			slog.Info("received cookie", "cookie", cookies[0])

		})
	}
}
