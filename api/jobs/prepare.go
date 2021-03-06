package jobs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Noah-Huppert/kube-git-deploy/api/libgh"
	"github.com/Noah-Huppert/kube-git-deploy/api/models"

	"github.com/BurntSushi/toml/cmd/tomlv"
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

// NewPrepareAction creates a new PrepareAction
func NewPrepareAction(ctx context.Context, logger golog.Logger,
	etcdKV etcd.KeysAPI) *PrepareAction {
	return &PrepareAction{
		ctx:    ctx,
		logger: logger,
		etcdKV: etcdKV,
	}
}

// Run executes the prepare action
func (a *PrepareAction) Run(job *models.Job, state *models.ActionState) error {
	// Set JobState.PrepareState.Stage to Running
	state.Stage = models.Running

	// { Download repository

	// Get GitHub repository download URL
	state.AddOutput("Initializing GitHub API")

	// ... Initialize GH client
	ghClient, err := libgh.NewClient(a.ctx, a.etcdKV)
	if err == libgh.ErrNoAuth {
		return errors.New("Not authenticated with GitHub")
	} else if err != nil {
		return fmt.Errorf("Error initializing GitHub API: %s",
			err.Error())
	}

	// ... Call API
	state.AddOutput("Retrieving repository download URL")

	dlURL, _, err := ghClient.Repositories.GetArchiveLink(a.ctx,
		job.ID.RepositoryID.Owner, job.ID.RepositoryID.Name, "tarball",
		&github.RepositoryContentGetOptions{
			Ref: job.Target.Commit,
		})
	if err != nil {
		return fmt.Errorf("Error retrieving repository download "+
			"URL: %s", err.Error())
	}

	// Create job working directory
	state.AddOutput("Setting up working directory")

	wrkDir := GetJobWorkingDir(*job)

	err = os.MkdirAll(wrkDir, 0777)
	if err != nil {
		return fmt.Errorf("Error creating working directory: %s",
			err.Error())
	}

	// Download GitHub repository contents
	state.AddOutput("Downloading GitHub repository")

	// ... Create file
	dlPath := fmt.Sprintf("%s/download.tar.gz", wrkDir)
	dlFile, err := os.Create(dlPath)
	if err != nil {
		return fmt.Errorf("Error creating repository download "+
			"file: %s", err.Error())
	}

	// ... Download file
	resp, err := http.Get(dlURL.String())
	if err != nil {
		return fmt.Errorf("Error making repository download "+
			"request: %s", err.Error())
	}

	// ... Copy request body to file
	_, err = io.Copy(dlFile, resp.Body)
	if err != nil {
		return fmt.Errorf("Error copying repository download request "+
			"body to file: %s", err.Error())
	}

	// ... Close HTTP request body
	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Error closing repository download "+
			"request body: %s", err.Error())
	}

	// ... Close file
	err = dlFile.Close()
	if err != nil {
		return fmt.Errorf("Error closing repository download "+
			"file: %s", err.Error())
	}

	// Extract download
	state.AddOutput("Extracting repository download file")

	// ... Unarchive tar file
	err = archiver.DefaultTarGz.Unarchive(dlPath, wrkDir)
	if err != nil {
		return fmt.Errorf("Error unarchiving repository download tar "+
			"file: %s", err.Error())
	}

	// Remove repository download file
	state.AddOutput("Removing repository download file")

	err = os.Remove(dlPath)
	if err != nil {
		return fmt.Errorf("Error removing repository download "+
			"file: %s", err.Error())
	}

	// Find extracted repository directory
	state.AddOutput("Finding working directory")

	fileInfos, err := ioutil.ReadDir(wrkDir)
	if err != nil {
		return fmt.Errorf("Error retrieving information about "+
			"files in working directory: %s", err.Error())
	}

	if len(fileInfos) != 1 {
		return fmt.Errorf("Working directory contained an unexpected "+
			"number of directories, count: %d", len(fileInfos))
	}

	if !fileInfos[0].IsDir() {
		return errors.New("Working directory does not contain the " +
			"expected sub-directory")
	}

	job.WorkingDir = fmt.Sprintf("%s/%s", wrkDir, fileInfos[0].Name())

	// }

	// { Parse configuration file
	// Read configuration file
	jobConfig := models.NewJobConfig()

	// }

	// Done
	state.Stage = models.Done

	return nil
}
