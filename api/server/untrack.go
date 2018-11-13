package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"
	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"
	"github.com/Noah-Huppert/kube-git-deploy/api/libgh"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
	etcd "go.etcd.io/etcd/client"
)

// UnrackGHRepoHandler untracks a GitHub repository
type UntrackGHRepoHandler struct {
	// ctx is context
	ctx context.Context

	// logger prints debug information
	logger golog.Logger

	// cfg is configuration
	cfg *config.Config

	// etcdKV is an Etcd key value API client
	etcdKV etcd.KeysAPI
}

// ServeHTTP implements http.Handler
func (h UntrackGHRepoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create responder
	responder := NewJSONResponder(h.logger, w)

	// Get URL parameters
	vars := mux.Vars(r)
	user := vars["user"]
	repo := vars["repo"]

	// Check doesn't already exist in Etcd
	found := true

	_, err := h.etcdKV.Get(h.ctx,
		libetcd.GetTrackedGHRepoNameKey(user, repo), nil)

	if err != nil && !etcd.IsKeyNotFound(err) {
		h.logger.Errorf("error determining if key exists: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "error determining if repository is" +
					"being tracked",
			})
		return
	} else if err != nil && etcd.IsKeyNotFound(err) {
		found = false
	}

	if !found {
		responder.Respond(http.StatusNotFound, map[string]interface{}{
			"ok":    false,
			"error": "repository not being tracked",
		})
		return
	}

	// Get GitHub hook ID
	resp, err := h.etcdKV.Get(h.ctx,
		libetcd.GetTrackedGHRepoWebHookIDKey(user, repo), nil)
	if err != nil {
		h.logger.Errorf("error retrieving GitHub hook ID from "+
			"Etcd: %s", err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "error retrieving web hook ID " +
					"from Etcd",
			})
		return
	}
	webHookIDStr := resp.Node.Value

	// Parse to int
	webHookID, err := strconv.ParseInt(webHookIDStr, 10, 64)
	if err != nil {
		h.logger.Errorf("error parsing retrieved GitHub hook ID into "+
			"integer: %s", err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "failed to interpret retrieve web " +
					"hook ID",
			})
		return
	}

	// Delete GitHub hook
	ghClient, err := libgh.NewClient(h.ctx, h.etcdKV)
	if err != nil {
		h.logger.Errorf("error creating GitHub client: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "error initializing GitHub API",
			})
		return
	}

	// ... Call GitHub hook API
	_, err = ghClient.Repositories.DeleteHook(h.ctx, user, repo, webHookID)
	if err != nil {
		h.logger.Errorf("error deleting web hook with GitHub API: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "error deleting web hook with " +
					"GitHub API",
			})
		return
	}

	// Delete Etcd directory
	_, err = h.etcdKV.Delete(h.ctx,
		libetcd.GetTrackedGHRepoDirKey(user, repo),
		&etcd.DeleteOptions{
			Recursive: true,
			Dir:       true,
		})
	if err != nil {
		h.logger.Errorf("error deleting tracked GitHub repo Etcd"+
			"directory: %s", err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "failed to delete tracked GitHub " +
					"repository in Etcd",
			})
		return
	}

	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
