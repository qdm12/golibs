package healthcheck

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/qdm12/golibs/logging"

	"github.com/qdm12/golibs/network"
)

// Mode checks if the program is launched to run the
// internal healthcheck against another instance of the
// program on the same machine
func Mode(args []string) bool {
	return len(args) > 1 && args[1] == "healthcheck"
}

// Query sends an HTTP request to the other instance of
// the program, and to its internal healthcheck server.
func Query() error {
	client := network.NewClient(time.Second)
	_, status, err := client.GetContent("http://127.0.0.1:9999")
	if err != nil {
		return err
	} else if status != http.StatusOK {
		return fmt.Errorf("HTTP status code is %d", status)
	}
	return nil
}

// CreateRouter creates a HTTP router with the route and configuration
// to run a healthcheck server locally
func CreateRouter(isHealthy func() error) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := isHealthy()
		if err != nil {
			logging.Warnf("Unhealthy: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return router
}
