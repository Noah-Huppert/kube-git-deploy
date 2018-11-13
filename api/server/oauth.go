package server

import (
	"context"
	"net/http"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"
	"github.com/Noah-Huppert/kube-git-deploy/api/github"
	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"

	"github.com/Noah-Huppert/golog"
	etcd "go.etcd.io/etcd/client"
)

// GHOAuthHandler exchanges a temporary GitHub code for an OAuth token
type GHOAuthHandler struct {
	// ctx is the context
	ctx context.Context

	// logger prints debug information
	logger golog.Logger

	// cfg is application configuration
	cfg *config.Config

	// etcdKV is the Etcd key value API client
	etcdKV etcd.KeysAPI
}

// ServeHTTP implements http.Handler
func (h GHOAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create JSON responder
	responder := NewJSONResponder(h.logger, w)

	// Try getting code from URL
	codeVals, ok := r.URL.Query()["code"]
	if !ok || len(codeVals) != 1 {
		responder.Respond(http.StatusBadRequest,
			map[string]interface{}{
				"ok":    false,
				"error": "\"code\" URL query parameter required",
			})
		return
	}

	code := codeVals[0]

	// Exchange with GitHub API
	exchangeReq := github.NewExchangeGitHubCodeReq(h.cfg, code)

	authToken, err := exchangeReq.Exchange()
	if err != nil {
		h.logger.Errorf("failed to exchange code with GitHub API: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "failed to exchange code with GitHub API",
			})
		return
	}

	// Save to Etcd
	_, err = h.etcdKV.Set(h.ctx, libetcd.KeyGitHubAuthToken, authToken,
		nil)
	if err != nil {
		h.logger.Errorf("failed to save GitHub auth token: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "failed to save GitHub auth token",
			})
		return
	}

	responder.Respond(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
