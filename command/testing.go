package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertExitCode is a ExecuteOption which uses WithExiter to assert the given exit code upon execution.
// This function is only useful for testing.
func AssertExitCode(t *testing.T, expectedExitCode int) ExecuteOption {
	t.Helper()
	hasRun := false
	t.Cleanup(func() {
		assert.Truef(t, hasRun, "exiter was not run")
	})
	return WithExiter(func(exitCode int) {
		hasRun = true
		assert.Equalf(t, expectedExitCode, exitCode, "exit code should be %d", expectedExitCode)
	})
}
