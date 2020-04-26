package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Settings contains settings to launch an HTTP server
// Name is simply a unique identifiant used for logging purposes.
type Settings struct {
	Name    string
	Addr    string
	Handler http.Handler
}

// RunServers manages multiple HTTP servers in parallel and stops
// them all as soon as one of them fails. It returns one error per HTTP
// server if there is any. It stops them also if the context is canceled.
func RunServers(ctx context.Context, settings ...Settings) (errors []error) {
	count := len(settings)
	chDone := make(chan error, count)
	chShutdownErr := make(chan error, count)
	chStop := make(chan struct{})
	for _, setting := range settings {
		setting := setting
		go serve(setting.Name, setting.Addr, setting.Handler, chStop, chDone, chShutdownErr)
	}
	stopped := false
	doneErrorsLeft, shutdownErrorsLeft := count, count
	for doneErrorsLeft > 0 || shutdownErrorsLeft > 0 {
		select {
		case err := <-chDone:
			if err != nil {
				errors = append(errors, err)
			}
			if !stopped {
				stopped = true
				close(chStop)
			}
			doneErrorsLeft--
		case err := <-chShutdownErr:
			if err != nil {
				errors = append(errors, err)
			}
			shutdownErrorsLeft--
		case <-ctx.Done():
			if !stopped {
				stopped = true
				close(chStop)
			}
		}
	}
	return errors
}

// serve listens on an address with the HTTP handler provided.
// It shuts down the server when it receives a signal from stop.
func serve(name, addr string, handler http.Handler, chStop <-chan struct{}, chDone, chShutdownErr chan<- error) {
	server := http.Server{Addr: addr, Handler: handler}
	go func() {
		<-chStop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			chShutdownErr <- fmt.Errorf("server %q failed shutting down: %w", name, err)
		} else {
			chShutdownErr <- nil
		}
	}()
	if err := server.ListenAndServe(); err != nil {
		chDone <- fmt.Errorf("server %q failed: %w", name, err)
	} else {
		chDone <- nil
	}
}
