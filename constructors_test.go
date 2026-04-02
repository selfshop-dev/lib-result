package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	result "github.com/selfshop-dev/lib-result"
)

func TestOk(t *testing.T) {
	t.Parallel()

	r := result.Ok[int, error](42)

	assert.True(t, r.IsOk())
	assert.False(t, r.IsErr())
	assert.Equal(t, 42, r.Unwrap())
}

func TestErr(t *testing.T) {
	t.Parallel()

	err := errors.New("something failed")
	r := result.Err[int](err)

	assert.False(t, r.IsOk())
	assert.True(t, r.IsErr())

	e, ok := r.Err()
	require.True(t, ok)
	assert.Equal(t, err, e)
}

func TestErr_PanicsOnNil(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		result.Err[int, error](nil)
	})
}

func TestErr_PanicsOnTypedNil(t *testing.T) {
	t.Parallel()

	var err *customError
	assert.Panics(t, func() {
		result.Err[int](err)
	})
}

func TestErr_NonNilableErrorType_DoesNotPanic(t *testing.T) {
	t.Parallel()

	// internalError is a value type — its reflect.Kind is Int32,
	// which is not nilable. This exercises the final return false
	// branch in isNilError.
	assert.NotPanics(t, func() {
		r := result.Err[int](internalError(1))
		assert.True(t, r.IsErr())
	})
}

func TestOf_WithValue(t *testing.T) {
	t.Parallel()

	r := result.Of(42, nil)

	assert.True(t, r.IsOk())
	assert.Equal(t, 42, r.Unwrap())
}

func TestOf_WithError(t *testing.T) {
	t.Parallel()

	err := errors.New("failed")
	r := result.Of(0, err)

	assert.True(t, r.IsErr())
	e, ok := r.Err()
	require.True(t, ok)
	assert.Equal(t, err, e)
}

func TestOf_BridgesIdiomaticGo(t *testing.T) {
	t.Parallel()

	// Of is the bridge from stdlib (value, error) — simulate os.Open pattern.
	openFile := func(name string) (string, error) {
		if name == "" {
			return "", errors.New("empty name")
		}
		return "content of " + name, nil
	}

	ok := result.Of(openFile("readme.txt"))
	assert.True(t, ok.IsOk())
	assert.Equal(t, "content of readme.txt", ok.Unwrap())

	bad := result.Of(openFile(""))
	assert.True(t, bad.IsErr())
}

func TestOfTyped_WithValue(t *testing.T) {
	t.Parallel()

	var nilErr *customError
	r := result.OfTyped(99, nilErr)

	assert.True(t, r.IsOk())
	assert.Equal(t, 99, r.Unwrap())
}

func TestOfTyped_WithError(t *testing.T) {
	t.Parallel()

	err := &customError{msg: "typed error"}
	r := result.OfTyped(0, err)

	assert.True(t, r.IsErr())
	e, ok := r.Err()
	require.True(t, ok)
	assert.Equal(t, "typed error", e.Error())
}

func TestOfTyped_NonNilableErrorType(t *testing.T) {
	t.Parallel()

	r := result.OfTyped(0, internalError(42))
	assert.True(t, r.IsErr())

	e, ok := r.Err()
	require.True(t, ok)
	assert.Equal(t, internalError(42), e)
}

// customError is a concrete error type used to test typed nil detection.
type customError struct{ msg string }

func (e *customError) Error() string { return e.msg }

// internalError implements error as a non-pointer value type.
// reflect.ValueOf(internalError(1)).Kind() == reflect.Int32 —
// falls through the switch in isNilError to the final return false.
type internalError int32

func (e internalError) Error() string { return fmt.Sprintf("internal error %d", e) }
