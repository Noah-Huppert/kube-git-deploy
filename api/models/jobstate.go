package models

// JobState holds information about the current run status of a Job.
type JobState struct {
	// Units holds unit states. Keys are UnitState.ID values.
	Units map[string]UnitState `json:"units"`
}

// UnitState holds the state of a unit.
type UnitState struct {
	// ID is the name of a unit
	ID string `json:"id"`

	// PrepareState is the state of the prepare action
	PrepareState ActionState `json:"prepare_state"`

	// DockerState is the state of the Docker action. Nil if the unit does
	// not contain a Docker action.
	DockerState *ActionState `json:"docker_state"`

	// HelmState is the state of the Helm action. Nil if the unit does not
	// contain a Helm action.
	HelmState *ActionState `json:"helm_state"`
}

// ActionState holds the state of an action.
type ActionState struct {
	// Stage indicates how the action is currently existing
	Stage ActionStage `json:"stage"`

	// Output holds the raw action output
	Output string `json:"output"`
}

// ActionStage indicates how the action is currently existing
type ActionStage string

const (
	// Queued indicates an action is set to be run, but hasn't started yet.
	Queued ActionStage = "queued"

	// Running indicates an action is running.
	Running ActionStage = "running"

	// Done indicates an action has finished running.
	Done ActionStage = "done"

	// ErrDone indicates that an action finished running because it
	// encountered an error.
	ErrDone ActionStage = "err_done"
)
