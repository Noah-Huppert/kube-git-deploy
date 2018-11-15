package jobs

import (
	"context"
	"errors"

	"github.com/Noah-Huppert/kube-git-deploy/api/models"

	etcd "go.etcd.io/etcd/client"
)

// JobRunner is responsible for running jobs.
type JobRunner struct {
	// ctx is context
	ctx context.Context

	// etcdKV is an etcd key value API client
	etcdKV etcd.KeysAPI

	// jobs holds all the currently running jobs. Keys are JobIDs.
	jobs map[models.JobID]models.Job

	// jobsChan accepts Jobs to run
	jobsChan chan models.Job
}

// NewJobRunner creates a new JobRunner
func NewJobRunner(ctx context.Context) *JobRunner {
	return &JobRunner{
		jobs:     map[models.JobID]models.Job{},
		jobsChan: make(chan models.Job),
	}
}

// Run starts the JobRunner main logic loop
func (r *JobRunner) Run() error {
	// Wait for job to be submitted
	for true {
		select {
		case job := <-r.jobsChan:
			// Check job isn't already running
			_, ok := r.jobs[job.ID]
			if ok {
				// Already running
				break
			}

			// Add to jobs map
			r.jobs[job.ID] = job

			// Execute job
			r.executeJob(job)

		case <-r.ctx.Done():
		}
	}
}

// executeJob runs the logic for a job. Should be started in a Go routine as
// it will block execution until the job finishes.
func (r *JobRunner) executeJob(job models.Job) {
	// Initialize
	// Prepare
}
