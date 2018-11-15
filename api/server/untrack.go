package server

import (
	"context"
	"net/http"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"
	"github.com/Noah-Huppert/kube-git-deploy/api/libgh"
	"github.com/Noah-Huppert/kube-git-deploy/api/models"

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
	name := vars["repo"]

	// Create repository model
	repo := models.Repository{
		ID: models.RepositoryID{
			Owner: user,
			Name:  name,
		},
	}

	// Check repository exists
	found, err := repo.Exists(h.ctx, h.etcdKV)
	if err != nil {
		h.logger.Errorf("error determining if repository exists: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "error determining if repository is" +
					"being tracked",
			})
		return
	}

	if !found {
		responder.Respond(http.StatusNotFound, map[string]interface{}{
			"ok":    false,
			"error": "repository not being tracked",
		})
		return
	}

	// Get GitHub hook ID
	err = repo.Get(h.ctx, h.etcdKV)
	if err != nil {
		h.logger.Errorf("error retrieving repository from Etcd: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "error retrieving repository " +
					"from Etcd",
			})
		return
	}

	// Delete GitHub hook
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

	// ... Call GitHub hook API
	_, err = ghClient.Repositories.DeleteHook(h.ctx, user, name,
		repo.WebHookID)
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
	err = repo.Delete(h.ctx, h.etcdKV)
	if err != nil {
		h.logger.Errorf("error deleting repository in Etcd",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "error deleting repository in Etcd",
			})
		return
	}

	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
