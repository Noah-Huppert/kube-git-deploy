package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"

	"github.com/gorilla/mux"
)

// RunServer starts the HTTP server
func RunServer(ctx context.Context, cfg *config.Config) error {
	// Setup routes
	router := mux.NewRouter()

	// Create server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: router,
	}

	// Create channel to return error
	errChan := make(chan error, 1)

	// Setup CTRL+C handler
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	go func() {
		<-interruptChan

		err := server.Shutdown(ctx)
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
