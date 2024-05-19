package v1

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	v1 "github.com/The-Gleb/gmessenger/app/gateway/controller/http/v1/middleware"
	groupmsgs_usecase "github.com/The-Gleb/gmessenger/app/gateway/domain/usecase/groupmsgs"
	"github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/handler/mocks"
	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/app/internal/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_groupMsgsHandler_ServeHTTP(t *testing.T) {
	validGroupMsgs := []entity.Message{
		{
			ID:     1,
			Sender: "login1",
			Text:   "Hello",
		},
	}

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockGroupUsecase := mocks.NewMockGroupMsgsUsecase(ctrl)
	groupmsgsHandler := NewGroupMsgsHandler(mockGroupUsecase)

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authMiddleware := v1.NewAuthMiddleware(mockAuthUsecase)
	mockAuthUsecase.
		EXPECT().
		Auth(gomock.Any(), gomock.Any()).
		Return("login1", nil).
		AnyTimes()

	groupmsgsHandler.Middlewares(authMiddleware.Http).AddToRouter(r)
	server := httptest.NewServer(r)

	type GroupmsgsUsecaseResp struct {
		messages []entity.Message
		err      error
	}

	tests := []struct {
		name                 string
		groupmsgsUsecaseResp GroupmsgsUsecaseResp
		code                 int
	}{
		{
			name: "positive",
			groupmsgsUsecaseResp: GroupmsgsUsecaseResp{
				messages: validGroupMsgs,
				err:      nil,
			},
			code: 200,
		},
		{
			name: "negative, group not found",
			groupmsgsUsecaseResp: GroupmsgsUsecaseResp{
				messages: nil,
				err:      errors.NewDomainError(errors.ErrGroupNotFound, ""),
			},
			code: 404,
		},
		{
			name: "negative, not a member",
			groupmsgsUsecaseResp: GroupmsgsUsecaseResp{
				messages: nil,
				err:      errors.NewDomainError(errors.ErrNotAMember, ""),
			},
			code: 403,
		},
		{
			name: "negative, some db err",
			groupmsgsUsecaseResp: GroupmsgsUsecaseResp{
				messages: nil,
				err:      errors.NewDomainError(errors.ErrDB, ""),
			},
			code: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockGroupUsecase.
				EXPECT().
				GetGroupMessages(gomock.Any(), gomock.Eq(groupmsgs_usecase.GetGroupMessagesDTO{
					GroupID:   1,
					UserLogin: "login1",
				})).
				Return(tt.groupmsgsUsecaseResp.messages, tt.groupmsgsUsecaseResp.err)

			resp, body := testRequest(t, "someToken", server, "GET", "/group/1", nil)
			defer resp.Body.Close()

			require.Equal(t, tt.code, resp.StatusCode)

			if tt.code != 200 {
				return
			}

			data, err := json.Marshal(validGroupMsgs)
			require.NoError(t, err)

			assert.Equal(t, body, string(data))

		})
	}
}
