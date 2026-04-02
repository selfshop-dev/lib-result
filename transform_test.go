package result_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	result "github.com/selfshop-dev/lib-result"
)

func TestMap(t *testing.T) {
	t.Parallel()

	t.Run("transforms Ok value", func(t *testing.T) {
		t.Parallel()
		r := result.Map(result.Ok[string, error]("hello"), strings.ToUpper)
		assert.Equal(t, "HELLO", r.Unwrap())
	})

	t.Run("passes Err through unchanged", func(t *testing.T) {
		t.Parallel()
		err := errors.New("fail")
		r := result.Map(result.Err[string](err), strings.ToUpper)
		require.True(t, r.IsErr())
		e, _ := r.Err()
		assert.Equal(t, err, e)
	})

	t.Run("can change type T to U", func(t *testing.T) {
		t.Parallel()
		r := result.Map(result.Ok[string, error]("hello"), func(s string) int { return len(s) })
		assert.Equal(t, 5, r.Unwrap())
	})
}

func TestMapErr(t *testing.T) {
	t.Parallel()

	t.Run("transforms Err value", func(t *testing.T) {
		t.Parallel()
		orig := errors.New("low-level")
		r := result.MapErr(
			result.Err[int](orig),
			func(e error) error { return errors.New("wrapped: " + e.Error()) },
		)
		require.True(t, r.IsErr())
		e, _ := r.Err()
		assert.Equal(t, "wrapped: low-level", e.Error())
	})

	t.Run("passes Ok through unchanged", func(t *testing.T) {
		t.Parallel()
		r := result.MapErr(
			result.Ok[int, error](42),
			func(_ error) error { return errors.New("should not be called") },
		)
		assert.Equal(t, 42, r.Unwrap())
	})

	t.Run("can change error type E to F", func(t *testing.T) {
		t.Parallel()

		r := result.MapErr(
			result.Err[int](errors.New("low-level")),
			func(e error) *customError { return &customError{msg: e.Error()} },
		)

		e, ok := r.Err()
		require.True(t, ok)
		assert.Equal(t, "low-level", e.msg)
	})
}

func TestAndThen(t *testing.T) {
	t.Parallel()

	double := func(n int) result.Result[int, error] {
		if n < 0 {
			return result.Err[int](errors.New("negative"))
		}
		return result.Ok[int, error](n * 2)
	}

	t.Run("chains Ok through fn", func(t *testing.T) {
		t.Parallel()
		r := result.AndThen(result.Ok[int, error](5), double)
		assert.Equal(t, 10, r.Unwrap())
	})

	t.Run("short-circuits on Err — fn not called", func(t *testing.T) {
		t.Parallel()
		called := false
		r := result.AndThen(
			result.Err[int](errors.New("initial")),
			func(n int) result.Result[int, error] {
				called = true
				return double(n)
			},
		)
		assert.True(t, r.IsErr())
		assert.False(t, called)
	})

	t.Run("fn may return Err", func(t *testing.T) {
		t.Parallel()
		r := result.AndThen(result.Ok[int, error](-1), double)
		assert.True(t, r.IsErr())
	})
}

func TestOrElse(t *testing.T) {
	t.Parallel()

	fallback := func(_ error) result.Result[int, error] {
		return result.Ok[int, error](-1)
	}

	t.Run("short-circuits on Ok — fn not called", func(t *testing.T) {
		t.Parallel()
		called := false
		r := result.OrElse(
			result.Ok[int, error](5),
			func(_ error) result.Result[int, error] {
				called = true
				return result.Ok[int, error](-1)
			},
		)
		assert.Equal(t, 5, r.Unwrap())
		assert.False(t, called)
	})

	t.Run("calls fn on Err", func(t *testing.T) {
		t.Parallel()
		r := result.OrElse(result.Err[int](errors.New("fail")), fallback)
		assert.Equal(t, -1, r.Unwrap())
	})

	t.Run("can change error type E to F", func(t *testing.T) {
		t.Parallel()

		r := result.OrElse(
			result.Err[int](errors.New("fail")),
			func(_ error) result.Result[int, *customError] {
				return result.Err[int](&customError{msg: "recovered"})
			},
		)

		e, ok := r.Err()
		require.True(t, ok)
		assert.Equal(t, "recovered", e.msg)
	})
}

func TestAnd(t *testing.T) {
	t.Parallel()

	t.Run("other is evaluated eagerly before And is called", func(t *testing.T) {
		t.Parallel()

		evaluated := false
		other := func() result.Result[string, error] {
			evaluated = true
			return result.Ok[string, error]("x")
		}()

		result.And(result.Err[int](errors.New("fail")), other)

		// Documents that And does NOT short-circuit evaluation of other —
		// unlike AndThen, which only calls fn when r is Ok.
		assert.True(t, evaluated)
	})

	other := result.Ok[string, error]("other")

	t.Run("returns other when Ok", func(t *testing.T) {
		t.Parallel()
		r := result.And(result.Ok[int, error](1), other)
		assert.Equal(t, "other", r.Unwrap())
	})

	t.Run("propagates Err ignoring other", func(t *testing.T) {
		t.Parallel()
		err := errors.New("first failed")
		r := result.And(result.Err[int](err), other)
		assert.True(t, r.IsErr())
		e, _ := r.Err()
		assert.Equal(t, err, e)
	})
}

func TestOr(t *testing.T) {
	t.Parallel()

	fallback := result.Ok[int, error](99)

	t.Run("returns self when Ok", func(t *testing.T) {
		t.Parallel()
		r := result.Or(result.Ok[int, error](1), fallback)
		assert.Equal(t, 1, r.Unwrap())
	})

	t.Run("returns other when Err", func(t *testing.T) {
		t.Parallel()
		r := result.Or(result.Err[int](errors.New("fail")), fallback)
		assert.Equal(t, 99, r.Unwrap())
	})
}
