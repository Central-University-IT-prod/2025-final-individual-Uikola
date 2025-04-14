package gigachat

import (
	"api/pkg/ai"
	"api/pkg/gigachat"
	"context"
	"fmt"

	"github.com/abadojack/whatlanggo"
)

// Client представляет клиента для GigaChat
type Client struct {
	client *gigachat.Client
}

// NewClient создает нового клиента
func NewClient(clientSecret string) (*Client, error) {
	gc, err := gigachat.NewInsecureClientWithAuthKey(clientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create gigachat client: %w", err)
	}

	return &Client{client: gc}, nil
}

// Call отправляет сообщение в GigaChat и возвращает ответ
func (c *Client) Call(ctx context.Context, req ai.CallRequest) (string, error) {
	if err := c.Auth(ctx); err != nil {
		return "", err
	}

	langInfo := whatlanggo.Detect(req.AdTitle)
	langCode := langInfo.Lang.Iso6393()

	message := fmt.Sprintf(req.Message, req.AdTitle, req.AdvertiserName, req.Context, langCode)
	request := &gigachat.ChatRequest{
		Model: gigachat.ModelLatest,
		Messages: []gigachat.Message{
			{Role: gigachat.SystemRole, Content: req.Prompt},
			{Role: gigachat.UserRole, Content: message},
		},
	}

	response, err := c.client.ChatWithContext(ctx, request)
	if err != nil {
		return "", fmt.Errorf("chat request failed: %w", err)
	}

	return response.Choices[0].Message.Content, nil
}

// Auth выполняет аутентификацию и устанавливает токен
func (c *Client) Auth(ctx context.Context) error {
	if err := c.client.AuthWithContext(ctx); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	return nil
}
