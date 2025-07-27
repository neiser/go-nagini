package command

import (
	"testing"

	"github.com/neiser/go-nagini/flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getFlagNames(t *testing.T) {
	t.Run("panics with non-pointer target", func(t *testing.T) {
		assert.PanicsWithValue(t, "given target must be of type pointer, but is of type string (value 'invalid')", func() {
			New().getFlagNames([]any{"invalid"})
		})
	})

	t.Run("panics with unknown target", func(t *testing.T) {
		var (
			someVal string
		)
		assert.PanicsWithValue(t, "cannot find flag names for target pointer, did you register with Flag(...) first?", func() {
			New().getFlagNames([]any{&someVal})
		})
	})
}

func TestCommand_MarkFlagsRequiredTogether(t *testing.T) {
	var (
		flag1, flag2 bool
	)
	cmd := New().
		Flag(flag.Bool(&flag1), flag.RegisterOptions{Name: "flag1"}).
		Flag(flag.Bool(&flag2), flag.RegisterOptions{Name: "flag2"}).
		MarkFlagsRequiredTogether(&flag1, &flag2).
		Run(func() error {
			// some dummy to make Cobra actually parse flags
			return nil
		})
	cmd.captureCobraOutput(t) // avoid confusing test output
	t.Run("one flag missing", func(t *testing.T) {
		require.ErrorContains(t,
			cmd.Execute(WithArgs("--flag1=false"), AssertExitCode(t, 1)),
			"if any flags in the group [flag1 flag2] are set they must all be set; missing [flag2]",
		)
		require.False(t, flag1)
		require.False(t, flag2)
	})
	t.Run("both flags provided", func(t *testing.T) {
		require.NoError(t, cmd.runTest(t, []string{"--flag1", "--flag2"}, func() {
			require.True(t, flag1)
			require.True(t, flag2)
		}))
	})
}

func TestCommand_MarkFlagsOneRequired(t *testing.T) {
	var (
		flag1, flag2 bool
	)
	cmd := New().
		Flag(flag.Bool(&flag1), flag.RegisterOptions{Name: "flag1"}).
		Flag(flag.Bool(&flag2), flag.RegisterOptions{Name: "flag2"}).
		MarkFlagsOneRequired(&flag1, &flag2).
		Run(func() error {
			// some dummy to make Cobra actually parse flags
			return nil
		})
	cmd.captureCobraOutput(t) // avoid confusing test output
	t.Run("no flags provided", func(t *testing.T) {
		require.ErrorContains(t,
			cmd.Execute(WithArgs(), AssertExitCode(t, 1)),
			"at least one of the flags in the group [flag1 flag2] is required",
		)
		require.False(t, flag1)
		require.False(t, flag2)
	})
	t.Run("one flag provided", func(t *testing.T) {
		require.NoError(t, cmd.runTest(t, []string{"--flag1"}, func() {
			require.True(t, flag1)
			require.False(t, flag2)
		}))
	})
}

func TestCommand_MarkFlagsMutuallyExclusive(t *testing.T) {
	var (
		flag1, flag2 bool
	)
	cmd := New().
		Flag(flag.Bool(&flag1), flag.RegisterOptions{Name: "flag1"}).
		Flag(flag.Bool(&flag2), flag.RegisterOptions{Name: "flag2"}).
		MarkFlagsMutuallyExclusive(&flag1, &flag2).
		Run(func() error {
			// some dummy to make Cobra actually parse flags
			return nil
		})
	cmd.captureCobraOutput(t) // avoid confusing test output
	t.Run("one flag provided", func(t *testing.T) {
		require.NoError(t, cmd.runTest(t, []string{"--flag1"}, func() {
			require.True(t, flag1)
			require.False(t, flag2)
		}))
	})
	t.Run("both flags provided", func(t *testing.T) {
		require.ErrorContains(t,
			cmd.Execute(WithArgs("--flag1", "--flag2"), AssertExitCode(t, 1)),
			"if any flags in the group [flag1 flag2] are set none of the others can be; [flag1 flag2] were all set",
		)
		require.True(t, flag1)
		require.True(t, flag2)
	})
}
