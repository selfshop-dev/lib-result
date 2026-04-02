package result_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	result "github.com/selfshop-dev/lib-result"
)

func TestIsOk_IsErr(t *testing.T) {
	t.Parallel()

	ok := result.Ok[int, error](1)
	assert.True(t, ok.IsOk())
	assert.False(t, ok.IsErr())

	bad := result.Err[int](errors.New("e"))
	assert.False(t, bad.IsOk())
	assert.True(t, bad.IsErr())
}

func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("Ok returns value and true", func(t *testing.T) {
		t.Parallel()
		r := result.Ok[string, error]("hello")
		v, ok := r.Value()
		require.True(t, ok)
		assert.Equal(t, "hello", v)
	})

	t.Run("Err returns zero and false", func(t *testing.T) {
		t.Parallel()
		r := result.Err[string](errors.New("e"))
		v, ok := r.Value()
		assert.False(t, ok)
		assert.Empty(t, v)
	})
}

func TestErr_Accessor(t *testing.T) {
	t.Parallel()

	t.Run("Err returns error and true on Err result", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("sentinel")
		r := result.Err[int](sentinel)
		e, ok := r.Err()
		require.True(t, ok)
		assert.Equal(t, sentinel, e)
	})

	t.Run("Err returns zero and false on Ok result", func(t *testing.T) {
		t.Parallel()
		r := result.Ok[int, error](1)
		e, ok := r.Err()
		assert.False(t, ok)
		assert.NoError(t, e)
	})
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	t.Run("returns value on Ok", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 42, result.Ok[int, error](42).Unwrap())
	})

	t.Run("panics on Err", func(t *testing.T) {
		t.Parallel()
		r := result.Err[int](errors.New("boom"))
		assert.Panics(t, func() { r.Unwrap() })
	})
}

func TestUnwrapErr(t *testing.T) {
	t.Parallel()

	t.Run("returns error on Err", func(t *testing.T) {
		t.Parallel()
		err := errors.New("fail")
		assert.Equal(t, err, result.Err[int](err).UnwrapErr())
	})

	t.Run("panics on Ok", func(t *testing.T) {
		t.Parallel()
		r := result.Ok[int, error](1)
		assert.Panics(t, func() { r.UnwrapErr() }) //nolint:gosec // intentional panic test — no error to handle
	})
}

func TestUnwrapOr(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 10, result.Ok[int, error](10).UnwrapOr(99))
	assert.Equal(t, 99, result.Err[int](errors.New("e")).UnwrapOr(99))
}

func TestUnwrapOrElse(t *testing.T) {
	t.Parallel()

	t.Run("returns value on Ok — fn not called", func(t *testing.T) {
		t.Parallel()
		called := false
		v := result.Ok[int, error](5).UnwrapOrElse(func(error) int {
			called = true
			return 0
		})
		assert.Equal(t, 5, v)
		assert.False(t, called)
	})

	t.Run("calls fn on Err", func(t *testing.T) {
		t.Parallel()
		err := errors.New("fail")
		v := result.Err[int](err).UnwrapOrElse(func(e error) int {
			assert.Equal(t, err, e)
			return -1
		})
		assert.Equal(t, -1, v)
	})
}

func TestUnwrapOrZero(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "hello", result.Ok[string, error]("hello").UnwrapOrZero())
	assert.Empty(t, result.Err[string](errors.New("e")).UnwrapOrZero())
}

func TestMust(t *testing.T) {
	t.Parallel()

	t.Run("returns value on Ok", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 7, result.Ok[int, error](7).Must("load"))
	})

	t.Run("panics with context on Err", func(t *testing.T) {
		t.Parallel()
		r := result.Err[int](errors.New("db unavailable"))
		assert.PanicsWithValue(t,
			"result: load config: db unavailable",
			func() { r.Must("load config") },
		)
	})
}
