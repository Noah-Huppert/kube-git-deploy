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
	// Setup routes
	router := mux.NewRouter()

	return Server{
		ctx:     ctx,
		logger:  logger.GetChild("http.public"),
		cfg:     cfg,
		etcdKV:  etcdKV,
		handler: router,
		port:    cfg.PublicHTTPPort,
	}
}
