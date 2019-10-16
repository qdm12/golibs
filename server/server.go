package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// RespondJSON marshals a payload into JSON and writes the JSON data
// to the response writer with a HTTP status code.
func RespondJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON formatting failed: %w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

// ServerSettings contains settings to launch an HTTP server
// Name is simply a unique identifiant used for logging purposes.
type ServerSettings struct {
	Name    string
	Addr    string
	Handler http.Handler
}

type namedError struct {
	name string
	err  error
}

// RunServers manages multiple HTTP servers in parallel and stops
// them all as soon as one of them fails. It returns one error per HTTP
// server if there is any.
func RunServers(settings ...ServerSettings) (errs map[string]error) {
	errs = make(map[string]error)
	for i := range settings {
		for j := range settings {
			if settings[i].Name == settings[j].Name {
				errs[settings[i].Name] = fmt.Errorf(
					"server settings at indexes %d and %d have the same name %s",
					i, j, settings[i].Name,
				)
				return errs
			}
		}
	}
	chDone := make(chan namedError, len(settings))
	chStop := make(chan struct{})
	for _, setting := range settings {
		setting := setting
		go serve(setting.Name, setting.Addr, setting.Handler, chDone, chStop)
	}
	var stopped bool
	for i := 0; i < cap(chDone); i++ {
		namedErr := <-chDone
		if namedErr.err != nil {
			errs[namedErr.name] = namedErr.err
		}
		if !stopped {
			stopped = true
			close(chStop)
		}
	}
	return errs
}

// serve listens on an address with the HTTP handler provided.
// It shuts down the server when it receives a signal from stop.
func serve(name, addr string, handler http.Handler, chDone chan namedError, chStop <-chan struct{}) {
	server := http.Server{Addr: addr, Handler: handler}
	go func() {
		<-chStop
		err := server.Shutdown(context.Background())
		if err != nil {
			zap.L().Error("server "+name+" shutdown error", zap.Error(err))
			zap.L().Sync()
		}
	}()
	zap.L().Info("HTTP server listening", zap.String("name", name), zap.String("addr", addr))
	err := server.ListenAndServe()
	chDone <- namedError{name, err}
}
