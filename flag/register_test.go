package flag

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestWithUsage(t *testing.T) {
	options := RegisterOptions{}.Apply(WithUsage("some usage"))
	assert.Equal(t, "some usage", options.Usage)
}

func TestPersistent(t *testing.T) {
	options := RegisterOptions{}.Apply(Persistent())
	cmd := &cobra.Command{}
	flags := options.SelectFlags(cmd)
	assert.Same(t, cmd.PersistentFlags(), flags)
}

func TestSelectFlags(t *testing.T) {
	options := RegisterOptions{}
	cmd := &cobra.Command{}
	flags := options.SelectFlags(cmd)
	assert.Same(t, cmd.Flags(), flags)
}
