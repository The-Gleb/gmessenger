package v1

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	v1 "github.com/The-Gleb/gmessenger/app/gateway/controller/http/v1/middleware"
	dialogmsgs_usecase "github.com/The-Gleb/gmessenger/app/gateway/domain/usecase/dialogmsgs"
	"github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/handler/mocks"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_dialogMsgsHandler_ServeHTTP(t *testing.T) {

	validDialogMsgs := []entity.Message{
		{
			ID:     1,
			Sender: "login1",
			Text:   "Hello",
		},
	}

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockDialogsUsecase := mocks.NewMockDialogMsgsUsecase(ctrl)
	dialogmsgsHandler := NewDialogMsgsHandler(mockDialogsUsecase)

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authMiddleware := v1.NewAuthMiddleware(mockAuthUsecase)
	mockAuthUsecase.
		EXPECT().
		Auth(gomock.Any(), gomock.Any()).
		Return("login1", nil).
		AnyTimes()

	dialogmsgsHandler.Middlewares(authMiddleware.Http).AddToRouter(r)
	server := httptest.NewServer(r)

	type DialogmsgsUsecaseResp struct {
		messages []entity.Message
		err      error
	}

	tests := []struct {
		name                  string
		dialogmsgsUsecaseResp DialogmsgsUsecaseResp
		code                  int
	}{
		{
			name: "positive",
			dialogmsgsUsecaseResp: DialogmsgsUsecaseResp{
				messages: validDialogMsgs,
				err:      nil,
			},
			code: 200,
		},
		{
			name: "negative, invalid receiver",
			dialogmsgsUsecaseResp: DialogmsgsUsecaseResp{
				messages: nil,
				err:      errors.NewDomainError(errors.ErrReceiverNotFound, ""),
			},
			code: 404,
		},
		{
			name: "negative, some err db",
			dialogmsgsUsecaseResp: DialogmsgsUsecaseResp{
				messages: nil,
				err:      errors.NewDomainError(errors.ErrDB, ""),
			},
			code: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDialogsUsecase.
				EXPECT().
				GetDialogMessages(gomock.Any(), gomock.Eq(dialogmsgs_usecase.GetDialogMessagesDTO{
					SenderLogin:   "login1",
					ReceiverLogin: "login2",
				})).
				Return(tt.dialogmsgsUsecaseResp.messages, tt.dialogmsgsUsecaseResp.err)

			resp, body := testRequest(t, "someToken", server, "GET", "/dialog/login2", nil)
			defer resp.Body.Close()

			require.Equal(t, tt.code, resp.StatusCode)

			if tt.code != 200 {
				return
			}

			data, err := json.Marshal(validDialogMsgs)
			require.NoError(t, err)

			assert.Equal(t, body, string(data))

		})
	}
}
