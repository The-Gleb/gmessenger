package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/handler/mocks"
	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/internal/errors"

	// v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/domain/service/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type Key string

func testRequest(t *testing.T, sToken string, ts *httptest.Server, method, path string, body []byte) (*http.Response, string) {
	req, err := http.NewRequestWithContext(context.Background(), method, ts.URL+path, bytes.NewReader(body))
	require.NoError(t, err)

	if sToken != "" {
		c := http.Cookie{
			Name:  "sessionToken",
			Value: sToken,
		}
		req.AddCookie(&c)
	}

	// login, ok := req.Context().Value(v1.Key("userLogin")).(string)
	// require.True(t, ok)
	// t.Log(login)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func Test_chatsHandler_ShowChats(t *testing.T) {

	validChatList := []entity.Chat{
		{
			Type:          client.Dialog,
			ReceiverLogin: "login2",
			Name:          "User2",
			LastMessage: entity.Message{
				ID:        1,
				Sender:    "login2",
				Text:      "Hello",
				Status:    "Sent",
				Timestamp: time.Now(),
			},
			Unread: 2,
		},
		{
			Type:          client.Dialog,
			ReceiverLogin: "login3",
			Name:          "User3",
			LastMessage: entity.Message{
				ID:        2,
				Sender:    "login1",
				Text:      "Hi",
				Status:    "Delivered",
				Timestamp: time.Now(),
			},
			Unread: 0,
		},
	}

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockChatsUsecase := mocks.NewMockChatsUsecase(ctrl)
	chatsHandler := NewChatsHandler(mockChatsUsecase)

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authMiddleware := v1.NewAuthMiddleware(mockAuthUsecase)

	chatsHandler.Middlewares(authMiddleware.Http).AddToRouter(r)
	server := httptest.NewServer(r)

	type AuthUsecaseResp struct {
		login string
		err   error
	}

	type ChatsUsecaseResp struct {
		chats []entity.Chat
		err   error
	}

	tests := []struct {
		name             string
		authUsecaseResp  AuthUsecaseResp
		chatsUsecaseResp ChatsUsecaseResp
		sessionToken     string
		code             int
	}{
		{
			name:         "positive",
			sessionToken: "123",
			authUsecaseResp: AuthUsecaseResp{
				login: "login1",
				err:   nil,
			},
			chatsUsecaseResp: ChatsUsecaseResp{
				chats: validChatList,
				err:   nil,
			},
			code: 200,
		},
		{
			name:         "sesion not found",
			sessionToken: "123",
			authUsecaseResp: AuthUsecaseResp{
				login: "",
				err:   errors.NewDomainError(errors.ErrNotAuthenticated, ""),
			},
			chatsUsecaseResp: ChatsUsecaseResp{
				chats: validChatList,
				err:   nil,
			},
			code: 401,
		},
		{
			name:         "some err in usecase",
			sessionToken: "123",
			authUsecaseResp: AuthUsecaseResp{
				login: "login1",
				err:   nil,
			},
			chatsUsecaseResp: ChatsUsecaseResp{
				chats: nil,
				err:   errors.NewDomainError(errors.ErrDB, ""),
			},
			code: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockAuthUsecase.EXPECT().Auth(gomock.Any(), gomock.Eq(tt.sessionToken)).
				Return(tt.authUsecaseResp.login, tt.authUsecaseResp.err)

			mockChatsUsecase.EXPECT().ShowChats(gomock.Any(), tt.authUsecaseResp.login).
				Return(tt.chatsUsecaseResp.chats, tt.chatsUsecaseResp.err)

			resp, body := testRequest(t, tt.sessionToken, server, "GET", "/chats", nil)
			defer resp.Body.Close()

			require.Equal(t, tt.code, resp.StatusCode)

			if tt.code != 200 {
				return
			}

			data, err := json.Marshal(validChatList)
			require.NoError(t, err)

			assert.Equal(t, body, string(data))

		})
	}
}
