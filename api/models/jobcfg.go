package models

// JobConfig holds information about the config of a job. Data
// sourced from a file in the Git repository root.
type JobConfig struct {
	// Units holds the config for the units in a file. Keys are
	// UnitConfig.ID values.
	Units map[string]UnitConfig `json:"units"`
}

// NewJobConfig creates a new JobConfig
func NewJobConfig() JobConfig {
	return JobConfig{
		Units: map[string]UnitConfig{},
	}
}

// UnitConfig holds the config for a unit
type UnitConfig struct {
	// ID holds the name of the unit
	ID string `json:"id"`

	// Docker holds Docker unit config. Nil if not present.
	Docker *DockerActionConfig `json:"docker"`

	// Helm holds Helm unit config. Nil if not present.
	Helm *HelmActionConfig `json:"helm"`
}

// DockerActionConfig holds the config for a Docker action.
type DockerActionConfig struct {
	// Directory indicates the directory where the Dockerfile to build
	// is located.
	Directory string `json:"directory"`

	// Tag indicates the value of the Docker image tag to apply.
	Tag string `json:"tag"`
}

// HelmActionConfig holds the config for a Helm action.
type HelmActionConfig struct {
	// Chart is the local path to a Helm chart to deploy, or if the
	// Repository field is set it holds the name of a Helm chart to deploy.
	Chart string `json:"chart"`

	// Repository is the name of the repository where the Chart is located.
	// If empty the Chart field is treated as a local path to a Helm chart.
	Repository string `json:"repository"`
}
