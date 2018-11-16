package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"

	etcd "go.etcd.io/etcd/client"
)

// Job holds information about a job.
type Job struct {
	// ID identifies the job.
	ID JobID `json:"id"`

	// Target identifies the Git event which triggered the job.
	Target JobTarget `json:"target"`

	// WorkingDir is the directory that the repository source is located.
	WorkingDir string `json:"working_dir"`

	// State holds the job state.
	State JobState `json:"state"`

	// Config holds the job configuration. Nil if it hasn't been
	// loaded yet.
	Config *JobConfig `json:"config"`
}

// NewJob creates a new Job. Intializes all JobState.Stage fields to Queued.
func NewJob(repoID RepositoryID, target JobTarget) *Job {
	j := Job{
		ID: JobID{
			RepositoryID: repoID,
		},
		Target: target,
	}

	// Initialize PrepareState
	j.State.PrepareState = NewActionState()

	// Initialize CleanupState
	j.State.CleanupState = NewActionState()

	return &j
}

// JobTarget identifies the Git event which triggered the job.
type JobTarget struct {
	// Branch is the Git branch.
	Branch string `json:"branch"`

	// Commit is the Git Sha.
	Commit string `json:"commit"`
}

// JobID identifies a job.
type JobID struct {
	// RepositoryID is the GitHub repository.
	RepositoryID RepositoryID `json:"repository_id"`

	// ID is the unique identifying number of the job.
	ID int64 `json:"id"`
}

// key indicates the Etcd in which job data will be stored.
func (i JobID) key() string {
	return fmt.Sprintf("%s/jobs/%d", i.RepositoryID.key(), i.ID)
}

// Create stores a new job. Finds the next job ID and saves it in the Job.ID
// field. Does not work if the job has already been stored
func (j *Job) Create(ctx context.Context, etcdKV etcd.KeysAPI) error {
	// Find next ID
	jobsDir := fmt.Sprintf("%s/jobs", j.ID.RepositoryID.key())

	resp, err := etcdKV.Get(ctx, jobsDir, &etcd.GetOptions{
		Recursive: true,
		Sort:      true,
		Quorum:    true,
	})

	if err != nil {
		return fmt.Errorf("error querying job IDs: %s", err.Error())
	}

	if resp.Node == nil {
		return errors.New("while finding next Job ID, result node " +
			"was nil")
	}

	var highestJobID int64 = -1
	for _, node := range resp.Node.Nodes {
		keyParts := strings.Split(node.Key, "/")
		jobIDStr := keyParts[len(keyParts)-1]

		jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing Job ID to int, "+
				"job ID: %s, error: %s", jobIDStr, err.Error())
		}

		if jobID > highestJobID {
			highestJobID = jobID
		}
	}

	j.ID.ID = highestJobID + 1

	// Save
	err = j.Set(ctx, etcdKV)
	if err != nil {
		return fmt.Errorf("error setting job: %s", err.Error())
	}

	return nil
}

// Set stores a job in Etcd
func (j Job) Set(ctx context.Context, etcdKV etcd.KeysAPI) error {
	return libetcd.SetJSON(ctx, etcdKV, j.ID.key(), j)
}

// Get retrieves a job from Etcd. The ID, Metadata.Owner, and
// Metadata.Name fields must be set for method to work properly.
func (j Job) Get(ctx context.Context, etcdKV etcd.KeysAPI) error {
	return libetcd.GetJSON(ctx, etcdKV, j.ID.key(), j)
}
