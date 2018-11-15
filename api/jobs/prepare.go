package jobs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Noah-Huppert/kube-git-deploy/api/libgh"
	"github.com/Noah-Huppert/kube-git-deploy/api/models"

	"github.com/Noah-Huppert/golog"
	"github.com/google/go-github/github"
	"github.com/mholt/archiver"
	etcd "go.etcd.io/etcd/client"
)

// GetJobWorkingDir returns the path to a job's working directory
func GetJobWorkingDir(job models.Job) string {
	return fmt.Sprintf("/tmp/kube-git-deploy/%s/%s/%d",
		job.ID.RepositoryID.Owner, job.ID.RepositoryID.Name, job.ID.ID)
}

// PrepareAction downloads a GitHub repository and parses the configuration
type PrepareAction struct {
	// ctx is context
	ctx context.Context

	// logger prints debug information
	logger golog.Logger

	// etcdKV is an Etcd key value API client
	etcdKV etcd.KeysAPI
}

// Run executes the prepare action
func (a PrepareAction) Run(job *models.Job, state *models.ActionState) {
	// Set JobState.PrepareState.Stage to Running
	state.Stage = models.Running

	// Get GitHub repository download URL
	state.AddOutput("Initializing GitHub API")

	// ... Initialize GH client
	ghClient, err := libgh.NewClient(a.ctx, a.etcdKV)
	if err == libgh.ErrNoAuth {
		state.SetError("Not authenticated " +
			"with GitHub")
		return
	} else if err != nil {
		state.SetErrorf("Error initializing GitHub "+
			"API: %s", err.Error())
		return
	}

	// ... Call API
	state.AddOutput("Retrieving repository download URL")

	dlURL, _, err := ghClient.Repositories.GetArchiveLink(a.ctx,
		job.ID.RepositoryID.Owner, job.ID.RepositoryID.Name, "tarball",
		&github.RepositoryContentGetOptions{
			Ref: job.Target.Commit,
		})
	if err != nil {
		state.SetErrorf("Error retrieving repository download "+
			"URL: %s", err.Error())
		return
	}

	// Create job working directory
	state.AddOutput("Setting up working directory")

	wrkDir := GetJobWorkingDir(job)

	err = os.MkdirAll(wrkDir, 0777)
	if err != nil {
		state.SetErrorf("Error creating working directory: %s",
			err.Error())
		return
	}

	// Download GitHub repository contents
	state.AddOutput("Downloading GitHub repository")

	// ... Create file
	dlPath := fmt.Sprintf("%s/download.tar", wrkDir)
	dlFile, err := os.Create(dlPath)
	if err != nil {
		state.SetErrorf("Error creating repository download "+
			"file: %s", err.Error())
		return
	}

	// ... Download file
	resp, err := http.Get(dlURL.String())
	if err != nil {
		state.SetErrorf("Error making repository download "+
			"request: %s", err.Error())
		return
	}

	// ... Copy request body to file
	_, err = io.Copy(dlFile, resp.Body)
	if err != nil {
		state.SetErrorf("Error copying repository download request "+
			"body to file: %s", err.Error())
		return
	}

	// ... Close HTTP request body
	err = resp.Body.Close()
	if err != nil {
		state.SetErrorf("Error closing repository download "+
			"request body: %s", err.Error())
		return
	}

	// ... Close file
	err = dlFile.Close()
	if err != nil {
		state.SetErrorf("Error closing repository download "+
			"file: %s", err.Error())
		return
	}

	// Extract download
	state.AddOutput("Extracting repository download file")

	// ... Open file
	rawTarFile, err := os.Open(dlPath)
	if err != nil {
		state.SetErrorf("Error opening repository download tar "+
			"file: %s", err.Error())
		return
	}

	// ... Open tar file
	tar := &archiver.Tar{}

	err = tar.Open(rawTarFile, -1)
	if err != nil {
		state.SetErrorf("Error opening repository download tar file "+
			"as tar file: %s", err.Error())
		return
	}

	// ... Unarchive tar file
	err = tar.Unarchive(".", wrkDir)
	if err != nil {
		state.SetErrorf("Error unarchiving repository download tar "+
			"file: %s", err.Error())
		return
	}

	// ... Close tar file
	err = tar.Close()
	if err != nil {
		state.SetErrorf("Error closing repository download tar "+
			"file: %s", err.Error())
		return
	}

	// ... Close file
	err = rawTarFile.Close()
	if err != nil {
		state.SetErrorf("Error closing repository download tar "+
			"file: %s", err.Error())
		return
	}

	// Done
	status.Stage = Done
}
