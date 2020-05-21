package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseOptions struct {
	status  int
	headers map[string]string
	body    []byte
}

func newResponseOptions() *responseOptions {
	return &responseOptions{headers: make(map[string]string)}
}

// ResponseSetter is a setter for options to Respond
type ResponseSetter func(options *responseOptions) error

// JSON sets the body as the JSON encoding of a payload
func JSON(payload interface{}) ResponseSetter {
	return func(options *responseOptions) error {
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("JSON formatting failed: %w", err)
		}
		options.body = body
		options.headers["Content-Type"] = "application/json"
		return nil
	}
}

// Bytes sets the body as the raw bytes given and sets the content type
// header depending on the contentType variable
func Bytes(b []byte, contentType string) ResponseSetter {
	return func(options *responseOptions) error {
		body := make([]byte, len(b))
		copy(body, b)
		options.body = body
		options.headers["Content-Type"] = contentType
		return nil
	}
}

// Status sets the status code to the HTTP response
func Status(status int) ResponseSetter {
	return func(options *responseOptions) error {
		options.status = status
		return nil
	}
}

// Respond responds using the HTTP writer and using the setters provided.
func Respond(w http.ResponseWriter, setters ...ResponseSetter) error {
	options := newResponseOptions()
	for _, setter := range setters {
		if err := setter(options); err != nil {
			return err
		}
	}
	for k, v := range options.headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(options.status)
	if _, err := w.Write(options.body); err != nil {
		return err
	}
	return nil
}
