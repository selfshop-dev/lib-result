package result_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	result "github.com/selfshop-dev/lib-result"
)

func TestToGo(t *testing.T) {
	t.Parallel()

	t.Run("Ok returns value and nil error", func(t *testing.T) {
		t.Parallel()
		v, err := result.Ok[int, error](42).ToGo()
		require.NoError(t, err)
		assert.Equal(t, 42, v)
	})

	t.Run("Err returns zero value and error", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("fail")
		v, err := result.Err[int](sentinel).ToGo()
		assert.Equal(t, sentinel, err)
		assert.Zero(t, v)
	})
}

func TestToGoTyped(t *testing.T) {
	t.Parallel()

	t.Run("Ok returns value and zero E", func(t *testing.T) {
		t.Parallel()
		v, e := result.Ok[int, *customError](7).ToGoTyped()
		assert.Equal(t, 7, v)
		assert.Nil(t, e)
	})

	t.Run("Err returns zero T and concrete error", func(t *testing.T) {
		t.Parallel()
		ce := &customError{msg: "typed"}
		v, e := result.Err[int](ce).ToGoTyped()
		assert.Zero(t, v)
		assert.Equal(t, "typed", e.Error())
	})
}

func TestOption(t *testing.T) {
	t.Parallel()

	t.Run("Ok returns value and true", func(t *testing.T) {
		t.Parallel()
		v, ok := result.Ok[string, error]("hi").Option()
		require.True(t, ok)
		assert.Equal(t, "hi", v)
	})

	t.Run("Err returns zero and false — error discarded", func(t *testing.T) {
		t.Parallel()
		v, ok := result.Err[string](errors.New("ignored")).Option()
		assert.False(t, ok)
		assert.Empty(t, v)
	})
}
