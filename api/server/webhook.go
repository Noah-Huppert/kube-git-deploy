package server

import (
	"encoding/json"
	"net/http"

	"github.com/Noah-Huppert/golog"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

// WebHookHandler triggers a build and deploy when GitHub sends a web
// hook request
type WebHookHandler struct {
	// logger prints debug information
	logger golog.Logger
}

// ServerHTTP implements http.Handler
func (h WebHookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create responder
	responder := NewJSONResponder(h.logger, w)

	// Get URL parameters
	vars := mux.Vars(r)
	user := vars["user"]
	repo := vars["repo"]

	// JSON decode body
	var event github.PushEvent

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&event)
	if err != nil {
		h.logger.Errorf("error decoding event into JSON: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "failed to interpret event",
			})
		return
	}

	h.logger.Debugf("%s/%s, event: %#v", user, repo, event)

	// Respond with OK
	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
