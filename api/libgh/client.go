package libgh

import (
	"context"
	"fmt"

	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"

	"github.com/google/go-github/github"
	etcd "go.etcd.io/etcd/client"
	"golang.org/x/oauth2"
)

// NewClient makes a new GitHub client with authentication
func NewClient(ctx context.Context, etcdKV etcd.KeysAPI) (*github.Client, error) {
	// Get GitHub auth token from Etcd
	resp, err := etcdKV.Get(ctx, libetcd.KeyGitHubAuthToken, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GitHub auth token"+
			" from Etcd: %s", err.Error())
	}

	authToken := resp.Node.Value

	// Create GitHub client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}
