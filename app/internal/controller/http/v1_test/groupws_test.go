package v1

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/handler/mocks"
	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func Test_groupWSHandler_ServeHTTP(t *testing.T) {
	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockGroupWSUsecase := mocks.NewMockGroupWSUsecase(ctrl)
	groupWSHandler := NewGroupWSHandler(mockGroupWSUsecase)

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authMiddleware := v1.NewAuthMiddleware(mockAuthUsecase)
	mockAuthUsecase.
		EXPECT().
		Auth(gomock.Any(), gomock.Any()).
		Return("login1", nil).
		AnyTimes()

	groupWSHandler.Middlewares(authMiddleware.Websocket).AddToRouter(r)
	server := httptest.NewServer(r)

	tests := []struct {
		name string
	}{
		{
			name: "positive",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// mocGroupWSUsecase.EXPECT().
			// OpenGroup(gomock.Any(), gomock.Eq(groupws_usecase.OpenGroupDTO{
			// 	Websocket: &websocket.Conn{},
			// 	SenderLogin: "login1",
			// 	SenderToken: "123",
			// 	ReceiverLogin: "login2",
			// })).DoAndReturn()

			clientConn, err := testWsRequest(t, server, "/group/login2/ws", "123")
			require.NoError(t, err)

			err = clientConn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(time.Second))
			require.NoError(t, err)

		})

	}
}
