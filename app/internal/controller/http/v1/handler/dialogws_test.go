package v1

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/handler/mocks"
	v1 "github.com/The-Gleb/gmessenger/app/internal/controller/http/v1/middleware"
	dialogws_usecase "github.com/The-Gleb/gmessenger/app/internal/domain/usecase/dialogws.go"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func testWsRequest(t *testing.T, server *httptest.Server, path, token string) (*websocket.Conn, error) {

	url := "ws" + strings.TrimPrefix(server.URL, "http") + path + "?token=" + token
	clientConn, response, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)

	err = response.Body.Close()
	require.NoError(t, err)

	return clientConn, nil
}

func Test_dialogWSHandler_ServeHTTP(t *testing.T) {

	r := chi.NewRouter()

	ctrl := gomock.NewController(t)
	mockDialogWSUsecase := mocks.NewMockDialogWSUsecase(ctrl)
	dialogWSHandler := NewDialogWSHandler(mockDialogWSUsecase)

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authMiddleware := v1.NewAuthMiddleware(mockAuthUsecase)
	mockAuthUsecase.
		EXPECT().
		Auth(gomock.Any(), gomock.Any()).
		Return("login1", nil).
		AnyTimes()

	dialogWSHandler.Middlewares(authMiddleware.Websocket).AddToRouter(r)
	server := httptest.NewServer(r)

	tests := []struct {
		name        string
		expectedDTO dialogws_usecase.OpenDialogDTO
	}{
		{
			name:        "positive",
			expectedDTO: dialogws_usecase.OpenDialogDTO{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clientConn, err := testWsRequest(t, server, "/dialog/login2/ws", "123")
			require.NoError(t, err)

			err = clientConn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(time.Second))
			require.NoError(t, err)

		})
	}
}
