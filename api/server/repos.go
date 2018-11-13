package server

import (
	"net/http"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"

	"github.com/Noah-Huppert/golog"
	//"github.com/google/go-github/github"
	//"golang.org/x/oauth2"
)

// GetReposHandler provides a list of GitHub repositories
type GetReposHandler struct {
	// logger is used to print debug information
	logger golog.Logger

	// cfg is configuration
	cfg *config.Config
}

// ServeHTTP implements http.Handler
func (h GetReposHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//
}
