package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"
	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"
	"github.com/Noah-Huppert/kube-git-deploy/api/libgh"

	"github.com/Noah-Huppert/golog"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	etcd "go.etcd.io/etcd/client"
)

// TrackGHRepoHandler marks a GitHub repository to be tracked
type TrackGHRepoHandler struct {
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
func (h TrackGHRepoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create responder
	responder := NewJSONResponder(h.logger, w)

	// Get URL parameters
	vars := mux.Vars(r)
	user := vars["user"]
	repo := vars["repo"]

	// Check doesn't already exist in Etcd
	ok := false
	_, err := h.etcdKV.Get(h.ctx,
		libetcd.GetTrackedGHRepoNameKey(user, repo), nil)

	if err != nil && !etcd.IsKeyNotFound(err) {
		h.logger.Errorf("error determining if key exists: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "error determining if repository is" +
					"already being tracked",
			})
		return
	} else if err != nil && etcd.IsKeyNotFound(err) {
		ok = true
	}

	if !ok {
		responder.Respond(http.StatusConflict, map[string]interface{}{
			"ok":    false,
			"error": "repository already being tracked",
		})
		return
	}

	// Create GitHub hook
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

	// ... Construct hook URL
	hookURL, err := url.Parse(h.cfg.PublicHTTPHost)
	if err != nil {
		h.logger.Errorf("error parsing public HTTP host into URL: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "error constructing web hook URL",
			})
		return
	}

	noSSLVerify := 1
	if h.cfg.PublicHTTPSSLEnabled {
		hookURL.Scheme = "https"
		noSSLVerify = 0
	}

	hookURL.Path = fmt.Sprintf("/api/v0/github/repositories/%s/%s/web_hook",
		user, repo)

	// ... Call GitHub hook API
	hook, _, err := ghClient.Repositories.CreateHook(h.ctx, user, repo,
		&github.Hook{
			Config: map[string]interface{}{
				"url":          hookURL.String(),
				"content_type": "json",
				"insecure_ssl": noSSLVerify,
			},
		})
	if err != nil {
		h.logger.Errorf("error creating web hook with GitHub API: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "error creating web hook with " +
					"GitHub API",
			})
		return
	}

	// Set in Etcd
	// ... Name
	_, err = h.etcdKV.Set(h.ctx,
		libetcd.GetTrackedGHRepoNameKey(user, repo),
		fmt.Sprintf("%s/%s", user, repo), nil)
	if err != nil {
		h.logger.Errorf("error saving tracked repo name in Etcd: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "failed to save repository name " +
					"in Etcd",
			})
		return
	}

	// ... Web hook ID
	_, err = h.etcdKV.Set(h.ctx,
		libetcd.GetTrackedGHRepoWebHookIDKey(user, repo),
		strconv.FormatInt(*(hook.ID), 10), nil)
	if err != nil {
		h.logger.Errorf("error saving tracked repository web hook "+
			"ID in Etcd: %s", err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "failed to save repository web " +
					"hook in Etcd",
			})
		return
	}

	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
