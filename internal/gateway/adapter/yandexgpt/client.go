package yandexgpt

import (
	"context"
	"github.com/sheeiavellie/go-yandexgpt"
	"log/slog"
)

type yandexGPTClient struct {
	catalogID string
	apiKey    string
	client    *yandexgpt.YandexGPTClient
}

func NewYandexGPTClient(catalogID, apiKey string) *yandexGPTClient {
	return &yandexGPTClient{
		catalogID: catalogID,
		apiKey:    apiKey,
		client:    yandexgpt.NewYandexGPTClientWithAPIKey(apiKey),
	}
}

func (c *yandexGPTClient) SendMessage(ctx context.Context, message string) (string, error) {
	request := yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(c.catalogID, yandexgpt.YandexGPTModelLite),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: 0.7,
			MaxTokens:   2000,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: message,
			},
		},
	}

	response, err := c.client.CreateRequest(context.Background(), request)
	if err != nil {
		slog.Info(err.Error())
		return "", err
	}

	return response.Result.Alternatives[0].Message.Text, nil

}
