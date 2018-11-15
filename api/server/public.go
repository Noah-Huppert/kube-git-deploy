package server

import (
	"context"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
	etcd "go.etcd.io/etcd/client"
)

// NewPublicServer creates a new server for public API endpoints
func NewPublicServer(ctx context.Context, logger golog.Logger,
	cfg *config.Config, etcdKV etcd.KeysAPI) Server {
	logger = logger.GetChild("http.public")

	// Setup routes
	router := mux.NewRouter()

	router.Handle("/healthz", HealthHandler{
		logger: logger.GetChild("health"),
		server: "public",
	}).Methods("GET")

	router.Handle("/api/v0/github/repositories/{user}/{repo}/web_hook",
		WebHookHandler{
			logger: logger.GetChild("github.webhook"),
		}).Methods("POST")

	// Create server
	return Server{
		ctx:     ctx,
		logger:  logger,
		cfg:     cfg,
		etcdKV:  etcdKV,
		handler: router,
		port:    cfg.PublicHTTPPort,
	}
}
