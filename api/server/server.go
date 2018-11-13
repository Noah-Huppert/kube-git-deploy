package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
	etcd "go.etcd.io/etcd/client"
)

// Server is a HTTP server
type Server struct {
	// ctx is context
	ctx context.Context

	// logger prints debug information
	logger golog.Logger

	// cfg is configuration
	cfg *config.Config

	// etcdKV is a Etcd key value API client
	etcdKV etcd.KeysAPI
}

func NewServer(ctx context.Context, logger golog.Logger, cfg *config.Config,
	etcdKV etcd.KeysAPI) Server {
	return Server{
		ctx:    ctx,
		logger: logger.GetChild("http"),
		cfg:    cfg,
		etcdKV: etcdKV,
	}
}

// Run starts the HTTP server
func (s Server) Run() error {
	// Setup routes
	router := mux.NewRouter()

	router.Handle("/api/v0/github/oauth_callback",
		GHOAuthHandler{
			ctx:    s.ctx,
			logger: s.logger.GetChild("github.oauth_callback"),
			cfg:    s.cfg,
			etcdKV: s.etcdKV,
		}).Methods("GET")

	router.Handle("/api/v0/github/login_url",
		GHLoginURLHandler{
			logger: s.logger.GetChild("github.login_url"),
			cfg:    s.cfg,
		}).Methods("GET")

	router.Handle("/api/v0/github/repositories/tracked",
		GetTrackedGHReposHandler{
			ctx:    s.ctx,
			logger: s.logger,
			etcdKV: s.etcdKV,
		}).Methods("GET")

	router.Handle("/api/v0/github/repositories/{user}/{repo}",
		TrackGHRepoHandler{
			ctx:    s.ctx,
			logger: s.logger,
			etcdKV: s.etcdKV,
		}).Methods("POST")

	// Create server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.HTTPPort),
		Handler: router,
	}

	// Create channel to return error
	errChan := make(chan error, 1)

	// Setup CTRL+C handler
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	go func() {
		<-interruptChan

		s.logger.Info("Shutting down server")

		err := server.Shutdown(s.ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to shutdown HTTP "+
				"server: %s", err.Error())
			return
		}

		errChan <- nil
	}()

	// Start server
	s.logger.Infof("Starting HTTP server on :%d", s.cfg.HTTPPort)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		errChan <- fmt.Errorf("error starting HTTP server: %s",
			err.Error())
	}

	// Return error
	s.logger.Info("Done running HTTP server")

	return <-errChan
}
