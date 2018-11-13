package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds application configuration
type Config struct {
	// HTTPPort is the port the API server will respond to requests on
	HTTPPort int `envconfig:"http_port" default:"5000"`

	// EtcdEndpoint is the host and port to a Etcd server
	EtcdEndpoint string `envconfig:"etcd_endpoint" default:"http://localhost:2379"`

	// GitHubClientID is the ID of the GitHub API app
	GitHubClientID string `envconfig:"github_client_id" required:"true"`

	// GitHubClientSecret is the secret value for a GitHub API app
	GitHubClientSecret string `envconfig:"github_client_secret" required:"true"`
}

// NewConfig loads configuration from the environment
func NewConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading configuration from"+
			" environment: %s", err.Error())
	}

	return &cfg, nil
}
