package libgh

import (
	"context"
	"errors"
	"fmt"

	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"

	"github.com/google/go-github/github"
	etcd "go.etcd.io/etcd/client"
	"golang.org/x/oauth2"
)

// ErrNoAuth indicates that no user is authenticated with GitHub
var ErrNoAuth error = errors.New("not authenticated")

// NewClient makes a new GitHub client with authentication
func NewClient(ctx context.Context, etcdKV etcd.KeysAPI) (*github.Client, error) {
	// Get GitHub auth token from Etcd
	resp, err := etcdKV.Get(ctx, libetcd.KeyGitHubAuthToken,
		&etcd.GetOptions{
			Quorum: true,
		})
	if etcd.IsKeyNotFound(err) {
		return nil, ErrNoAuth
	} else if err != nil {
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
