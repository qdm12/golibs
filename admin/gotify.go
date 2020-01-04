package admin

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
)

// Gotify contains the Gotify API client and the token for the application
type Gotify struct {
	client *client.GotifyREST
	token  string
}

// InitGotify creates a Gotify client from environment variables
// and logs warnings if the parameters are not valid.
func InitGotify(httpClient *http.Client) *Gotify {
	if httpClient == nil {
		logging.Warn("Setting HTTP client to default")
		httpClient = &http.Client{Timeout: 3 * time.Second}
	}
	gotifyURL, err := params.GetGotifyURL()
	if err != nil {
		logging.Err(err)
	}
	gotifyToken, err := params.GetGotifyToken()
	if err != nil {
		logging.Werr(err)
	}
	gotify, err := NewGotify(gotifyURL, gotifyToken, httpClient)
	if err != nil {
		logging.Werr(err)
	}
	return gotify
}

// NewGotify creates an API client with the token for the Gotify server
func NewGotify(URL *url.URL, token string, httpClient *http.Client) (g *Gotify, err error) {
	if URL == nil {
		return &Gotify{}, fmt.Errorf("Gotify URL not provided")
	} else if token == "" {
		return &Gotify{}, fmt.Errorf("Gotify token not provided")
	}
	client := gotify.NewClient(URL, httpClient)
	_, err = client.Version.GetVersion(nil)
	if err != nil {
		return &Gotify{}, fmt.Errorf("cannot communicate with Gotify server: %w", err)
	}
	return &Gotify{client: client, token: token}, nil
}

// Notify sends a notification to the Gotify server
func (g *Gotify) Notify(title string, priority int, content string, args ...interface{}) {
	if err := g.notify(title, priority, content, args...); err != nil {
		logging.Errorf("Gotify error: %s", err)
	}
}

func (g *Gotify) notify(title string, priority int, content string, args ...interface{}) error {
	if g.client == nil {
		return nil
	}
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
