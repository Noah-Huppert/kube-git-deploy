package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"

	"github.com/Noah-Huppert/golog"
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

	// handler is the handler used to respond to requests
	handler http.Handler

	// port is the port HTTP traffic will be served on
	port int
}

// Run starts the HTTP server
func (s Server) Run() error {
	// Create server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.handler,
	}

	// Create channel to return error
	errChan := make(chan error, 1)

	// Setup CTRL+C handler
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	go func() {
		<-interruptChan

		err := server.Shutdown(s.ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to shutdown HTTP "+
				"server: %s", err.Error())
			return
		}

		errChan <- nil
	}()

	// Start server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		errChan <- fmt.Errorf("error starting HTTP server: %s",
			err.Error())
	}

	// Return error
	return <-errChan
}
