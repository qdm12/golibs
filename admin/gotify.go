package admin

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
)

// Gotify is a Gotify client
type Gotify interface {
	Ping() error
	Notify(title string, priority int, content string, args ...interface{}) error
}

type gotifyImpl struct {
	client *client.GotifyREST
	token  string
}

// NewGotify creates an API client with the token for the Gotify server
func NewGotify(URL url.URL, token string, httpClient *http.Client) Gotify {
	client := gotify.NewClient(&URL, httpClient)
	return &gotifyImpl{client: client, token: token}
}

func (g *gotifyImpl) Ping() error {
	if _, err := g.client.Version.GetVersion(nil); err != nil {
		return fmt.Errorf("cannot communicate with Gotify server: %w", err)
	}
	return nil
}

// Notify sends a notification to the Gotify server
func (g *gotifyImpl) Notify(title string, priority int, content string, args ...interface{}) error {
	content = fmt.Sprintf(content, args...)
	params := message.NewCreateMessageParams()
	params.Body = &models.MessageExternal{
		Title:    title,
		Message:  content,
		Priority: priority,
	}
	_, err := g.client.Message.CreateMessage(params, auth.TokenAuth(g.token))
	if err != nil {
		return fmt.Errorf("cannot send message to Gotify: %w", err)
	}
	return nil
}
