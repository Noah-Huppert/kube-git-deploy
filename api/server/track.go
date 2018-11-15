package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"
	"github.com/Noah-Huppert/kube-git-deploy/api/libgh"
	"github.com/Noah-Huppert/kube-git-deploy/api/models"

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
	name := vars["repo"]

	// Create repo model
	repo := models.Repository{
		ID: models.RepositoryID{
			Owner: user,
			Name:  name,
		},
	}

	// Check doesn't already exist in Etcd
	exists, err := repo.Exists(h.ctx, h.etcdKV)
	if err != nil {
		h.logger.Errorf("error determining if repository exists: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "error determining if repository " +
					"is already being tracked",
			})

		return
	}

	if exists {
		responder.Respond(http.StatusConflict, map[string]interface{}{
			"ok":    false,
			"error": "repository already being tracked",
		})
		return
	}

	// Create GitHub hook
	ghClient, err := libgh.NewClient(h.ctx, h.etcdKV)
	if err == libgh.ErrNoAuth {
		responder.Respond(http.StatusUnauthorized,
			map[string]interface{}{
				"ok":    false,
				"error": libgh.ErrNoAuth.Error(),
			})
		return

	} else if err != nil {
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
		user, name)

	// ... Call GitHub hook API
	hook, _, err := ghClient.Repositories.CreateHook(h.ctx, user, name,
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

	// ... Save web hook ID in repository
	repo.WebHookID = *(hook.ID)

	// Save repository
	err = repo.Create(h.ctx, h.etcdKV)
	if err != nil {
		h.logger.Errorf("error saving repository to Etcd: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "error saving repository to Etcd",
			})

		return
	}

	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
