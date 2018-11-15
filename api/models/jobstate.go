package models

import (
	"fmt"
)

// JobState holds information about the current run status of a Job.
type JobState struct {
	// PrepareState is the state of the prepare action
	PrepareState *ActionState `json:"prepare_state"`

	// CleanupState is the state of the cleanup action
	CleanupState *ActionState `json:"cleanup_state"`

	// Units holds unit states. Keys are UnitState.ID values.
	Units map[string]UnitState `json:"units"`
}

// Done indicates if the Job has finished executing
func (s JobState) Done() bool {
	if !(s.PrepareState.Done() && s.CleanupState.Done()) {
		return false
	}

	for _, v := range s.Units {
		if !(v.DockerState.Done() && v.HelmState.Done()) {
			return false
		}
	}

	return true
}

// UnitState holds the state of a unit.
type UnitState struct {
	// ID is the name of a unit
	ID string `json:"id"`

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
	Output []ActionOutput `json:"output"`
}

// Done indicates if a state's Stage is in a done state
func (s ActionState) Done() bool {
	return s.Stage == Done || s.Stage == ErrDone
}

// NewActionState creates a new ActionState with the Stage field set to Queued
func NewActionState() *ActionState {
	return &ActionState{
		Stage: Queued,
	}
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

// ActionOutput holds a line of output from an action. Indicates if the line
// is an error or normal.
type ActionOutput struct {
	// Text is the output string
	Text string `json:"text"`

	// Error indicates if the text is error output
	Error bool
}

// SetError saves an error to in Output and sets the Stage to ErrDone
func (s *ActionState) SetError(errStr string) {
	s.Stage = ErrDone
	s.Output = append(s.Output, ActionOutput{
		Text:  errStr,
		Error: true,
	})
}

// SetErrorf acts like SetError but it provides string formatting functionality
// via fmt.Sprintf
func (s *ActionState) SetErrorf(errFormat string, v ...interface{}) {
	s.SetError(fmt.Sprintf(errFormat, v))
}

// AddOutput saves a line of output to the state
func (s *ActionState) AddOutput(txt string) {
	s.Output = append(s.Output, ActionOutput{
		Text:  txt,
		Error: false,
	})
}
