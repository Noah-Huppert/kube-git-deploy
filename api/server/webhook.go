package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Noah-Huppert/kube-git-deploy/api/models"

	"github.com/Noah-Huppert/golog"
	"github.com/google/go-github/github"
	//"github.com/gorilla/mux"
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
	//vars := mux.Vars(r)
	//user := vars["user"]
	//repo := vars["repo"]

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

	// Make Job Target
	// ... Parse branch
	refParts := strings.Split(*(event.Ref), "/")

	if len(refParts) != 3 {
		h.logger.Errorf("error, ref not in \"refs/head/<branch>\" "+
			"format, ref: %s", *(event.Ref))

		responder.Respond(http.StatusBadRequest,
			map[string]interface{}{
				"ok":    false,
				"error": "failed to parse ref field",
			})
		return
	}
	branch := refParts[2]

	// ... Make struct
	jobTarget := models.JobTarget{
		Branch: branch,
		Commit: *(event.After),
	}

	h.logger.Debugf("JobTarget: %#v", jobTarget)

	// Respond with OK
	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
