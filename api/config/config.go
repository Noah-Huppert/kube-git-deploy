package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds application configuration
type Config struct {
	// PrivateHTTPPort is the port the API server will respond to
	// private API requests on
	PrivateHTTPPort int `envconfig:"private_http_port" default:"5000"`

	// PublicHTTPPort is the port the API server will respond to public
	// API requests on
	PublicHTTPPort int `envconfig:"public_http_port" default:"5001"`

	// PublicHTTPHost is the host name which the public API will be
	// server under
	PublicHTTPHost string `envconfig:"public_http_host" required:"true"`

	// PublicHTTPSSLEnabled indicates if the public API server has a
	// SSL certificate
	PublicHTTPSSLEnabled bool `envconfig:"public_http_ssl_enabled" default:"false"`

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
