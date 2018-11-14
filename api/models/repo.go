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
	WebHookID int64 `json:"web_hook_id"`
}

// GetTrackedGHRepoDirKey returns the directory key for a tracked
// GitHub repository
func GetTrackedGHRepoDirKey(user, repo string) string {
	return fmt.Sprintf("%s/%s/%s", libetcd.KeyDirTrackedGHRepos, user,
		repo)
}

// key returns the Etcd key the repository should be stored in
func (r Repository) key() string {
	return fmt.Sprintf("%s/information",
		GetTrackedGHRepoDirKey(r.Owner, r.Name))
}

// GetAll retrieves all repositories
func GetAll(ctx context.Context, etcdKV etcd.KeysAPI) ([]Repository, error) {
	// TODO
	return nil, nil
}

// Exists checks to see if repository exists in Etcd
func (r Repository) Exists(ctx context.Context,
	etcdKV etcd.KeysAPI) (bool, error) {

	_, err := etcdKV.Get(ctx, r.key(), &etcd.GetOptions{
		Quorum: true,
	})

	if err != nil && !etcd.IsKeyNotFound(err) {
		return false, fmt.Errorf("error querying Etcd for "+
			"repository: %s", err.Error())
	} else if err != nil && etcd.IsKeyNotFound(err) {
		return false, nil
	}

	return true, nil
}

// Creates creates the directory structure for a Repository in Etcd and
// Sets a directory. Does not work if Repository was previously saved.
func (r Repository) Create(ctx context.Context, etcdKV etcd.KeysAPI) error {
	// Create top level directory
	dirName := GetTrackedGHRepoDirKey(r.Owner, r.Name)

	_, err := etcdKV.Set(ctx, dirName, "", &etcd.SetOptions{
		Dir:       true,
		PrevExist: etcd.PrevNoExist,
	})

	if err != nil {
		return fmt.Errorf("error creating repository directory: %s",
			err.Error())
	}

	// Create jobs directory
	jobsDirName := fmt.Sprintf("%s/jobs", dirName)

	_, err = etcdKV.Set(ctx, jobsDirName, "", &etcd.SetOptions{
		Dir:       true,
		PrevExist: etcd.PrevNoExist,
	})

	if err != nil {
		return fmt.Errorf("error creating jobs directory: %s",
			err.Error())
	}

	// Save repository in directory
	err = r.Set(ctx, etcdKV)

	if err != nil {
		return fmt.Errorf("error saving repository: %s", err.Error())
	}

	return nil
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
