package parameter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_anySliceValue_Type(t *testing.T) {
	t.Run("built-in type", func(t *testing.T) {
		var (
			target []string
		)
		sut := NewSlice(&target, nil, ParseSliceOf[string](NotEmptyString[string]))
		assert.Equal(t, "[]string", sut.Type())
	})
	t.Run("custom type", func(t *testing.T) {
		type someType string
		type someTypes []someType
		var (
			target someTypes
		)
		sut := NewSlice[someType](&target, nil, ParseSliceOf[someType](NotEmptyString[someType]))
		assert.Equal(t, "[]someType", sut.Type())
	})
}
