package client

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/The-Gleb/gmessenger/app/internal/domain/entity"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHub struct {
}

func (mh *mockHub) RouteEvent(event entity.Event, senderClient *Client) {

	senderClient.Message <- event

}
func (mh *mockHub) AddClient(c *Client) {

}
func (mh *mockHub) RemoveClient(c *Client) {

}

func GetWSConnections() (server *httptest.Server, clientConn, serverConn *websocket.Conn, err error) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		var err error
		serverConn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
	}))

	var response *http.Response
	clientConn, response, err = websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if err != nil {
		return nil, nil, nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, nil, nil, err
	}

	return server, clientConn, serverConn, nil
}

func Test_WriteMessage(t *testing.T) {

	validNewMessageEventJson, err := json.Marshal(entity.NewMessageEvent{
		ID:          1,
		SenderLogin: "login1",
		Status:      entity.DELIVERED,
		Text:        "hello",
		CreatedAt:   time.Now(),
	})
	require.NoError(t, err)

	// hub := Hub.AddClient()

	tests := []struct {
		name    string
		event   entity.Event
		errCode int
	}{
		{
			name: "positive",
			event: entity.Event{
				Type:    entity.NewMessage,
				Payload: validNewMessageEventJson,
			},
			errCode: 0,
		},
		{
			name:    "negative, empty event",
			event:   entity.Event{},
			errCode: 1011,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server, clientConn, serverConn, err := GetWSConnections()
			require.NoError(t, err)

			mockHub := mockHub{}

			client := Client{
				Conn:    serverConn,
				Message: make(chan entity.Event, 5),
				Hub:     &mockHub,
			}

			fmt.Println(tt.errCode)
			go client.WriteMessage()

			client.Message <- tt.event

			time.Sleep(time.Second)

			_, msg, err := clientConn.ReadMessage()

			if tt.errCode != 0 {
				log.Println(err)
				// assert.Equal(t, tt.errCode, tt.clientConn.CloseCode)
				isUnexpected := websocket.IsUnexpectedCloseError(err, tt.errCode)
				require.False(t, isUnexpected)
				return
			}
			require.NoError(t, err)

			var event entity.Event
			err = json.Unmarshal(msg, &event)
			require.NoError(t, err)
			assert.Equal(t, tt.event, event)

			server.Close()
			clientConn.Close()
			serverConn.Close()

		})
	}
}

func Test_ReadMessage(t *testing.T) {

	validSendMessageEventJson, err := json.Marshal(entity.SendMessageEvent{
		Text: "Hello",
	})
	require.NoError(t, err)
	validEvent, err := json.Marshal(entity.Event{
		Type:    entity.SendMessage,
		Payload: validSendMessageEventJson,
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		event   entity.Event
		message json.RawMessage
		errCode int
	}{
		{
			name: "positive",
			event: entity.Event{
				Type:    entity.NewMessage,
				Payload: validSendMessageEventJson,
			},
			message: validEvent,
			errCode: 0,
		},
		{
			name:    "negative, ivalid msg",
			event:   entity.Event{},
			message: []byte("some invalid message"),
			errCode: 1007,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server, clientConn, serverConn, err := GetWSConnections()
			require.NoError(t, err)

			mockHub := mockHub{}

			client := Client{
				Conn:    serverConn,
				Message: make(chan entity.Event, 5),
				Hub:     &mockHub,
			}

			go client.ReadMessage()

			time.Sleep(time.Millisecond * 10) //nolint:all

			err = clientConn.WriteMessage(websocket.TextMessage, tt.message)
			require.NoError(t, err)

			time.Sleep(time.Second)

			if tt.errCode != 0 {
				t.Log(err)
				log.Println(err)
				_, _, err := clientConn.ReadMessage()
				slog.Error("here read error", "error", err)
				isUnexpected := websocket.IsUnexpectedCloseError(err, tt.errCode)
				require.False(t, isUnexpected)
				return
			}

			var event entity.Event
			err = json.Unmarshal(tt.message, &event)
			require.NoError(t, err)

			receivedMessage := <-client.Message

			assert.Equal(t, event, receivedMessage)

			server.Close()
			clientConn.Close()
			serverConn.Close()

		})
	}
}
