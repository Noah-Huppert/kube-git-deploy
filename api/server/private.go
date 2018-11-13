package server

import (
	"context"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
	etcd "go.etcd.io/etcd/client"
)

// NewPrivateServer creates a new server for private API endpoints
func NewPrivateServer(ctx context.Context, logger golog.Logger,
	cfg *config.Config, etcdKV etcd.KeysAPI) Server {
	// Setup routes
	router := mux.NewRouter()

	router.Handle("/api/v0/github/oauth_callback",
		GHOAuthHandler{
			ctx:    ctx,
			logger: logger.GetChild("github.oauth_callback"),
			cfg:    cfg,
			etcdKV: etcdKV,
		}).Methods("GET")

	router.Handle("/api/v0/github/login_url",
		GHLoginURLHandler{
			logger: logger.GetChild("github.login_url"),
			cfg:    cfg,
		}).Methods("GET")

	router.Handle("/api/v0/github/repositories/tracked",
		GetTrackedGHReposHandler{
			ctx:    ctx,
			logger: logger,
			etcdKV: etcdKV,
		}).Methods("GET")

	router.Handle("/api/v0/github/repositories/{user}/{repo}",
		TrackGHRepoHandler{
			ctx:    ctx,
			logger: logger,
			cfg:    cfg,
			etcdKV: etcdKV,
		}).Methods("POST")

	router.Handle("/api/v0/github/repositories/{user}/{repo}",
		UntrackGHRepoHandler{
			ctx:    ctx,
			logger: logger,
			cfg:    cfg,
			etcdKV: etcdKV,
		}).Methods("DELETE")

	return Server{
		ctx:     ctx,
		logger:  logger.GetChild("http.private"),
		cfg:     cfg,
		etcdKV:  etcdKV,
		handler: router,
		port:    cfg.PrivateHTTPPort,
	}
}
