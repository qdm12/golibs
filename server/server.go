package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Serve listens on an address with the HTTP handler provided.
// It shuts down the server when it receives a signal from stop.
func Serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	server := http.Server{Addr: addr, Handler: handler}
	go func() {
		<-stop
		server.Shutdown(context.Background())
	}()
	return server.ListenAndServe()
}

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
