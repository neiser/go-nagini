package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_anyValue_Type(t *testing.T) {
	t.Run("built-in type", func(t *testing.T) {
		var (
			target string
		)
		sut := String(&target, NotEmptyTrimmed)
		assert.Equal(t, "string", sut.Type())
	})
	t.Run("custom type", func(t *testing.T) {
		type someType string
		var (
			target someType
		)
		sut := String(&target, NotEmptyTrimmed)
		assert.Equal(t, "flag.someType", sut.Type())
	})
}

type someStringer string

func (s someStringer) String() string {
	return "prefix " + string(s)
}

func Test_convertToString(t *testing.T) {
	t.Run("use format bool", func(t *testing.T) {
		assert.Equal(t, "true", convertToString(true))
	})
	t.Run("use format string", func(t *testing.T) {
		assert.Equal(t, "foo", convertToString("foo"))
	})
	t.Run("use fmt.Stringer interface", func(t *testing.T) {
		assert.Equal(t, "prefix foo", convertToString(someStringer("foo")))
	})
}
