package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
	etcd "go.etcd.io/etcd/client"
)

// TrackGHRepoHandler marks a GitHub repository to be tracked
type TrackGHRepoHandler struct {
	// ctx is context
	ctx context.Context

	// logger prints debug information
	logger golog.Logger

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

	// Set in Etcd
	_, err := h.etcdKV.Set(h.ctx,
		libetcd.GetTrackedGitHubRepoKey(user, repo),
		fmt.Sprintf("%s/%s", user, repo),
		&etcd.SetOptions{
			PrevExist: etcd.PrevNoExist,
		})
	if err != nil && err.(etcd.Error).Code != etcd.ErrorCodeNodeExist {
		h.logger.Errorf("error setting tracked repo name in Etcd: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "failed track repository in Etcd",
			})
		return
	}

	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
