package models

import (
	"context"
	"fmt"

	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"

	etcd "go.etcd.io/etcd/client"
)

// Repository contains tracked GitHub repository information
type Repository struct {
	// Owner holds the GitHub username of the repository owner
	Owner string `json:"owner"`

	// Name holds the name of the repository
	Name string `json:"name"`

	// WebHookID holds the ID of the created GitHub repository web hook
	WebHookID int `json:"web_hook_id"`
}

// key returns the Etcd key the repository should be stored in
func (r *Repository) key() string {
	return fmt.Sprintf("%s/%s/%s/information", libetcd.KeyDirTrackedGHRepos, r.Owner,
		r.Name)
}

// Set stores a repository in Etcd
func (r Repository) Set(ctx context.Context, etcdKV etcd.KeysAPI) error {
	return libetcd.SetJSON(ctx, etcdKV, r.key(), r)
}

// Get retrieves a repository from Etcd. The `Owner` and `Name` fields must be
// set for this method to work properly
func (r *Repository) Get(ctx context.Context, etcdKV etcd.KeysAPI) error {
	return libetcd.GetJSON(ctx, etcdKV, r.key(), r)
}
