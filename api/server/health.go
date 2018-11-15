package server

import (
	"net/http"

	"github.com/Noah-Huppert/golog"
)

// HealthHandler returns an OK response for external services to verify the
// server is running
type HealthHandler struct {
	// logger prints debug information
	logger golog.Logger

	// server indicates the name of the API server
	server string
}

// ServeHTTP implements http.Handler
func (h HealthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	// Make responder
	responder := NewJSONResponder(h.logger, w)

	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok":     true,
		"server": h.server,
	})
}
