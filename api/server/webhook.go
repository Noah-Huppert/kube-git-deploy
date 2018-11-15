package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Noah-Huppert/kube-git-deploy/api/jobs"
	"github.com/Noah-Huppert/kube-git-deploy/api/models"

	"github.com/Noah-Huppert/golog"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	etcd "go.etcd.io/etcd/client"
)

// WebHookHandler triggers a build and deploy when GitHub sends a web
// hook request
type WebHookHandler struct {
	// ctx is context
	ctx context.Context

	// logger prints debug information
	logger golog.Logger

	// etcdKV is an Etcd key value API client
	etcdKV etcd.KeysAPI

	// jobRunner is used to run jobs
	jobRunner *jobs.JobRunner
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

	// Save job in Etcd
	job := models.NewJob(models.RepositoryID{
		Owner: user,
		Name:  repo,
	}, jobTarget)

	err = job.Create(h.ctx, h.etcdKV)
	if err != nil {
		h.logger.Errorf("error saving Job in Etcd: %s", err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "failed to save job in Etcd",
			})

		return
	}

	// Run job
	h.logger.Debugf("submitting job")
	h.jobRunner.Submit(job)
	h.logger.Debugf("submitted job")

	// Respond with OK
	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
