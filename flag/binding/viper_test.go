package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViper_BindTo(t *testing.T) {
	t.Run("empty config key returns nil binder", func(t *testing.T) {
		assert.Nil(t, Viper{}.BindTo())
	})
}
