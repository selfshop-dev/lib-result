package result_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	result "github.com/selfshop-dev/lib-result"
)

func TestString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "Ok(42)", result.Ok[int, error](42).String())
	assert.Equal(t, "Ok(hello)", result.Ok[string, error]("hello").String())
	assert.Equal(t, "Err(something failed)", result.Err[int](errors.New("something failed")).String())
}

func TestResult_ZeroValue(t *testing.T) {
	t.Parallel()

	var r result.Result[int, error]

	// Zero value behaves as Err, not Ok.
	assert.False(t, r.IsOk())
	assert.True(t, r.IsErr())

	// Value and Err return zero values of their respective types.
	v, ok := r.Value()
	assert.False(t, ok)
	assert.Zero(t, v)

	e, ok := r.Err()
	assert.True(t, ok)
	require.NoError(t, e) // zero E is nil for error

	// Unwrap panics — zero value is not a valid Ok.
	assert.Panics(t, func() { r.Unwrap() })
}

// TestValue_Alias verifies that Value[T] is a usable alias for Result[T, error].
func TestValue_Alias(t *testing.T) {
	t.Parallel()

	r := result.Of(42, nil)
	assert.True(t, r.IsOk())
	assert.Equal(t, 42, r.Unwrap())

	bad := result.Of(0, errors.New("fail"))
	assert.True(t, bad.IsErr())
}
