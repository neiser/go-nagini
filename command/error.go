package command

import (
	"fmt"
)

// WithExitCodeError allows to set an exit code for that error.
// See Command.Execute.
type WithExitCodeError struct {
	ExitCode int
	Wrapped  error
}

func (e WithExitCodeError) Unwrap() error {
	return e.Wrapped
}

func (e WithExitCodeError) Error() string {
	if e.Wrapped != nil {
		return e.Wrapped.Error()
	}
	return fmt.Sprintf("exit code %d", e.ExitCode)
}

type fromRunCallbackError struct {
	Wrapped error
}

func (e fromRunCallbackError) Unwrap() error {
	return e.Wrapped
}

func (e fromRunCallbackError) Error() string {
	return e.Wrapped.Error()
}
