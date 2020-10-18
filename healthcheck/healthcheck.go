package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/golibs/network"
)

// Mode checks if the program is launched to run the
// internal healthcheck against another instance of the
// program on the same machine.
func Mode(args []string) bool {
	return len(args) > 1 && args[1] == "healthcheck"
}

// Query sends an HTTP request to the other instance of
// the program, and to its internal healthcheck server.
func Query(ctx context.Context) error {
	client := network.NewClient(time.Second)
	_, status, err := client.Get(ctx, "http://127.0.0.1:9999")
	if err != nil {
		return err
	} else if status != http.StatusOK {
		return fmt.Errorf("HTTP status code is %d", status)
	}
	return nil
}

// GetHandler returns a handler function to handle healthcheck queries locally.
func GetHandler(isHealthy func() error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := isHealthy()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
