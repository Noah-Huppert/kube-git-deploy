package server

import (
	"context"
	"net/http"

	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"

	"github.com/Noah-Huppert/golog"
	etcd "go.etcd.io/etcd/client"
)

// GetTrackedGHReposHandler returns a list of tracked GitHub repositories
type GetTrackedGHReposHandler struct {
	// ctx is context
	ctx context.Context

	// logger prints debug information
	logger golog.Logger

	// etcdKV is an Etcd key value API client
	etcdKV etcd.KeysAPI
}

// ServeHTTP implements http.Handler
func (h GetTrackedGHReposHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create responder
	responder := NewJSONResponder(h.logger, w)

	// Get tracked repositories
	resp, err := h.etcdKV.Get(h.ctx, libetcd.KeyDirTrackedGHRepos,
		&etcd.GetOptions{
			Recursive: true,
			Sort:      true,
			Quorum:    true,
		})
	if err != nil {
		h.logger.Errorf("error getting tracked GitHub repos from "+
			"Etcd: %s", err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok": false,
				"error": "failed to retrieve tracked GitHub" +
					" repositories from Etcd",
			})
		return
	}

	repos := libetcd.TraverseDir(resp.Node, libetcd.KeyTrackedGHRepoName)

	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok":           true,
		"repositories": repos,
	})
}
