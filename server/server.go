package server

import (
	"context"
	"fmt"
	"net/http"
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
// server if there is any.
func RunServers(settings ...Settings) (errors []error) {
	count := len(settings)
	chDone := make(chan error, count)
	chShutdownErr := make(chan error, count)
	chStop := make(chan struct{})
	for _, setting := range settings {
		setting := setting
		go serve(setting.Name, setting.Addr, setting.Handler, chStop, chDone, chShutdownErr)
	}
	var stopped bool
	i := 0
	for i < count {
		select {
		case err := <-chDone:
			i++
			if err != nil {
				errors = append(errors, err)
			}
			if !stopped {
				stopped = true
				close(chStop)
			}
		case err := <-chShutdownErr: // best effort to collect shutdown errors
			errors = append(errors, err)
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
		err := server.Shutdown(context.Background())
		if err != nil {
			chShutdownErr <- fmt.Errorf("server %q failed shutting down: %w", name, err)
		}
	}()
	err := server.ListenAndServe()
	if err != nil {
		chDone <- fmt.Errorf("server %q failed: %w", name, err)
	} else {
		chDone <- nil
	}
}
