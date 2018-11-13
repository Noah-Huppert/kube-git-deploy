package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/golog"
)

// JSONResponder responds to HTTP requests with JSON
type JSONResponder struct {
	// logger prints debug information
	logger golog.Logger

	// w is the HTTP respond writer
	w http.ResponseWriter
}

// NewJSONResponder creates a new JSONResponder
func NewJSONResponder(logger golog.Logger, w http.ResponseWriter) JSONResponder {
	return JSONResponder{
		logger: logger,
		w:      w,
	}
}

// RespondJSON responds to an HTTP request with a JSON
func (r JSONResponder) Respond(status int, v interface{}) {
	// Set content type
	r.w.Header().Set("Content-Type", "application/json")

	// Set status
	r.w.WriteHeader(status)

	// Write JSON
	encoder := json.NewEncoder(r.w)

	err := encoder.Encode(v)
	if err != nil {
		// Write manual error response
		r.logger.Errorf("failed to respond with JSON, value: %#v, "+
			"error: %s", v, err.Error())

		status = http.StatusInternalServerError
		fmt.Fprintf(r.w, "{\"error\": \"failed to respond with JSON\"}")
	}
}
