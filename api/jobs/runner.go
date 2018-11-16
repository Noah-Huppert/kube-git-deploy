package jobs

import (
	"context"

	"github.com/Noah-Huppert/kube-git-deploy/api/models"

	"github.com/Noah-Huppert/golog"
	etcd "go.etcd.io/etcd/client"
)

// JobRunner is responsible for running jobs.
type JobRunner struct {
	// ctx is context
	ctx context.Context

	// logger prints debug information
	logger golog.Logger

	// etcdKV is an etcd key value API client
	etcdKV etcd.KeysAPI

	// jobs holds all the currently running jobs. Keys are JobIDs.
	jobs map[models.JobID]*models.Job

	// jobsChan accepts Jobs to run
	jobsChan chan *models.Job
}

// NewJobRunner creates a new JobRunner
func NewJobRunner(ctx context.Context, logger golog.Logger) *JobRunner {
	return &JobRunner{
		ctx:      ctx,
		jobs:     map[models.JobID]*models.Job{},
		jobsChan: make(chan *models.Job),
	}
}

// Submit sends a job to the runner main loop for future execution
func (r *JobRunner) Submit(job *models.Job) {
	r.jobsChan <- job
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
				r.logger.Debugf("already running: %#v", job)
				break
			}

			// Add to jobs map
			r.jobs[job.ID] = job

			r.logger.Debugf("received job: %#v", job)

			// Execute job
			go r.executeJob(job)

		case <-r.ctx.Done():
			r.logger.Info("Job Runner stopping")
			return nil
		}
	}

	return nil
}

// executeJob runs the logic for a job. Should be started in a Go routine as
// it will block execution until the job finishes.
func (r *JobRunner) executeJob(job *models.Job) {
	// Prepare
	// ... Run
	prepareAction := PrepareAction{
		ctx:    r.ctx,
		logger: r.logger,
		etcdKV: r.etcdKV,
	}

	prepareOK := true

	err := prepareAction.Run(job, job.State.PrepareState)
	if err != nil {
		r.logger.Errorf("error running prepare action, Job.ID: %#v "+
			", error: %s", job.ID, err.Error())

		job.State.PrepareState.SetError(err.Error())

		prepareOK = false
	}

	// ... Save
	err = job.Set(r.ctx, r.etcdKV)
	if err != nil {
		r.logger.Errorf("error saving job after prepare action, "+
			"Job.ID: %#v, error: %s", job.ID, err.Error())
		return
	}

	if !prepareOK {
		return
	}
}
