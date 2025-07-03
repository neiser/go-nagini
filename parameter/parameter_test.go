package parameter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_anyValue_Type(t *testing.T) {
	t.Run("built-in type", func(t *testing.T) {
		var (
			target string
		)
		sut := New(&target, "", NotEmptyString)
		assert.Equal(t, "string", sut.Type())
	})
	t.Run("custom type", func(t *testing.T) {
		type someType string
		var (
			target someType
		)
		sut := New(&target, "", NotEmptyString)
		assert.Equal(t, "parameter.someType", sut.Type())
	})
}
