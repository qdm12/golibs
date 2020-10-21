package admin

import (
	"fmt"
	"net/http"
	liburl "net/url"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	gotifylib "github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/qdm12/golibs/format"
	"github.com/qdm12/golibs/logging"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Gotify

// Gotify is a Gotify client.
type Gotify interface {
	// Ping obtains silently the version from the Gotify server and returns an error on failure
	Ping() error
	// Notify formats and sends a message to the Gotify server
	Notify(title string, priority int, args ...interface{}) error
	NotifyAndLog(title string, level logging.Level, logger logging.Logger, args ...interface{})
}

type gotify struct {
	client *client.GotifyREST
	token  string
}

// NewGotify creates an API client with the token for the Gotify server.
func NewGotify(url liburl.URL, token string, httpClient *http.Client) Gotify {
	client := gotifylib.NewClient(&url, httpClient)
	return &gotify{client: client, token: token}
}

// Ping obtains silently the version from the Gotify server and returns an error on failure.
func (g *gotify) Ping() error {
	if _, err := g.client.Version.GetVersion(nil); err != nil {
		return fmt.Errorf("cannot communicate with Gotify server: %w", err)
	}
	return nil
}

func (g *gotify) notify(title, content string, priority int) error {
	if g == nil {
		return nil
	}
	if content == "" {
		content = " " // content is required
	}
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

// Notify formats and sends a message to the Gotify server.
// TODO custom HTTP request with context.
func (g *gotify) Notify(title string, priority int, args ...interface{}) error {
	s := format.ArgsToString(args...)
	return g.notify(title, s, priority)
}

func (g *gotify) NotifyAndLog(title string, level logging.Level, logger logging.Logger, args ...interface{}) {
	s := format.ArgsToString(args...)
	toLog := fmt.Sprintf("%s: %s", title, s)
	var priority int
	switch level {
	case logging.DebugLevel:
		logger.Debug(toLog)
	case logging.InfoLevel:
		priority = 1
		logger.Info(toLog)
	case logging.WarnLevel:
		priority = 2
		logger.Warn(toLog)
	case logging.ErrorLevel:
		priority = 3
		logger.Error(toLog)
	default:
		logger.Debug(toLog)
	}
	if err := g.notify(title, s, priority); err != nil {
		logger.Error(err)
	}
}
