package v1

import (
	"encoding/json"
	v1 "github.com/The-Gleb/gmessenger/app/gateway/controller/http/v1/handler"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	register_usecase "github.com/The-Gleb/gmessenger/app/gateway/domain/usecase/register"
	"github.com/The-Gleb/gmessenger/app/internal/controller/http/dto"
	"github.com/The-Gleb/gmessenger/app/internal/controller/http/v1_test/mocks"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_registerHandler_Register(t *testing.T) {
	validRegisterReqBody, err := json.Marshal(dto.RegisterUserDTO{
		Username: "user1",
		Email:    "user1@gmail.com",
		Password: "password1",
	})
	require.NoError(t, err)

	invalidRegisterReqBody, err := json.Marshal(dto.RegisterUserDTO{
		Username: "",
		Email:    "",
		Password: "password1",
	})
	require.NoError(t, err)

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)

	mockRegisterUsecase := mocks.NewMockRegisterUsecase(ctrl)
	registerHandler := v1.NewRegisterHandler(mockRegisterUsecase)
	registerHandler.AddToRouter(r)
	server := httptest.NewServer(r)

	tests := []struct {
		name    string
		reqBody json.RawMessage
		code    int
		prepare func()
	}{
		{
			name:    "positive",
			reqBody: validRegisterReqBody,
			code:    200,
			prepare: func() {
				mockRegisterUsecase.
					EXPECT().
					Register(gomock.Any(), gomock.Eq(entity.RegisterUserDTO{
						Username: "",
						Email:    "",
						Password: "password1",
					})).
					Return(entity.Session{
						UserID: 12,
						Token:  "123",
						Expiry: time.Now().Add(time.Hour),
					}, nil)
			},
		},
		{
			name:    "negative, body with no username",
			reqBody: invalidRegisterReqBody,
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
			reqBody: validRegisterReqBody,
			code:    409,
			prepare: func() {
				mockRegisterUsecase.
					EXPECT().
					Register(gomock.Any(), gomock.Eq(register_usecase.RegisterUserDTO{
						Login:    "login1",
						Password: "password1",
						UserName: "user1",
					})).
					Return(entity.Session{}, errors.NewDomainError(errors.ErrDBLoginAlredyExists, ""))
			},
		},
		{
			name:    "negative, some db err",
			reqBody: validRegisterReqBody,
			code:    500,
			prepare: func() {
				mockRegisterUsecase.
					EXPECT().
					Register(gomock.Any(), gomock.Eq(register_usecase.RegisterUserDTO{
						Login:    "login1",
						Password: "password1",
						UserName: "user1",
					})).
					Return(entity.Session{}, errors.NewDomainError(errors.ErrDB, ""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.prepare()

			resp, _ := testRequest(t, "", server, "POST", "/register", tt.reqBody)
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
