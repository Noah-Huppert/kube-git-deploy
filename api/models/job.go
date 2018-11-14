package models

import (
	"context"
	"fmt"

	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"

	etcd "go.etcd.io/etcd/client"
)

// Job holds repository build deploy job information
type Job struct {
	// ID holds the unique ID of the job
	ID int `json:"id"`

	// Modules holds information about modules defined in repository
	// configuration file
	Modules []Module `json:"modules"`

	// Metadata holds information about the event which triggered the job
	Metadata JobMetadata `json:"metadata"`
}

// key returns the Etcd key a job should be stored in
func (j Job) key() string {
	return fmt.Sprintf("%s/%s/%s/jobs/%d", libetcd.KeyDirTrackedGHRepos,
		j.Metadata.Owner, j.Metadata.Name, j.ID)
}

// Set stores a job in Etcd
func (j Job) Set(ctx context.Context, etcdKV etcd.KeysAPI) error {
	return libetcd.SetJSON(ctx, etcdKV, j.key(), j)
}

// Get retrieves a job from Etcd. The `ID`, `Metadata.Owner`, and
// `Metadata.Name` fields must be set for method to work properly.
func (j Job) Get(ctx context.Context, etcdKV etcd.KeysAPI) error {
	return libetcd.GetJSON(ctx, etcdKV, j.key(), j)
}

// JobMetadata holds information about the event which triggered the job
type JobMetadata struct {
	// Owner holds the GitHub username of the repository owner
	Owner string `json:"owner"`

	// Name holds the name of the repository
	Name string `json:"name"`

	// Branch is the Git branch of the commit which triggered the job
	Branch string `json:"branch"`

	// CommitSha is the Git commit sha of the commit which triggered
	// the job
	CommitSha string `json:"commit_sha"`
}

// Module holds information about an individual item which can be built
// and deployed
type Module struct {
	// Configuration holds the raw step configuration from the repository
	// configuration file
	Configuration StepConfiguration `json:"configuration"`

	// State holds the state of the module's steps
	State StepsState `json:"state"`
}

// StepConfiguration holds the step configuration from a repository
// configuration file
type StepConfiguration struct {
	// Docker holds the configuration for a Docker step, nil if
	// not included
	Docker *DockerStepConfiguration `json:"docker"`

	// Helm holds the configuration for a Helm step, nil if not included
	Helm *HelmStepConfiguration `json:"helm"`
}

// DockerStepConfiguration holds configuration for a Docker step
type DockerStepConfiguration struct {
	// Directory is the directory the Dockerfile to build is located in
	Directory string `json:"directory"`

	// Tag is the value to tag the built Docker image with
	Tag string `json:"tag"`
}

// HelmStepConfiguration holds configuration for a Helm step
type HelmStepConfiguration struct {
	// Chart is the path to the Helm chart to deploy
	Chart string `json:"chart"`

	// Repository is the name of the Helm repository to retrieve the Chart
	// from. If empty the Chart path is considered to be a path to a
	// local directory.
	Repository string `json:"repository"`
}

// StepsState holds state information for steps in a module
type StepsState struct {
	// DockerState holds the state of a Docker step, nil if not present
	DockerState *StepState `json:"docker_state"`

	// HelmState holds the state of a Helm step, nil if not present
	HelmState *StepState `json:"helm_state"`
}

// StepState holds state information about a single step
type StepState struct {
	// Status indicates the current run status of the step
	Status StepStatus `json:"status"`

	// Output holds the raw build output
	Output string `json:"output"`
}

// StepStatus is used to indicate the current status of a step
type StepStatus string

const (
	// Waiting indicates the step is set to run, but has not started yet
	Waiting StepStatus = "waiting"

	// Running indicates the step is running
	Running StepStatus = "running"

	// Success indicates the step successfully completed
	Success StepStatus = "success"

	// Error indicates the step failed to complete
	Error StepStatus = "error"
)
