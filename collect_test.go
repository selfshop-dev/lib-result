package result_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	result "github.com/selfshop-dev/lib-result"
)

func TestCollect(t *testing.T) {
	t.Parallel()

	t.Run("all Ok returns Ok slice", func(t *testing.T) {
		t.Parallel()
		rs := []result.Result[int, error]{
			result.Ok[int, error](1),
			result.Ok[int, error](2),
			result.Ok[int, error](3),
		}
		r := result.Collect(rs)
		require.True(t, r.IsOk())
		assert.Equal(t, []int{1, 2, 3}, r.Unwrap())
	})

	t.Run("returns first Err and stops", func(t *testing.T) {
		t.Parallel()
		firstErr := errors.New("first")
		rs := []result.Result[int, error]{
			result.Ok[int, error](1),
			result.Err[int](firstErr),
			result.Err[int](errors.New("second")),
		}
		r := result.Collect(rs)
		require.True(t, r.IsErr())
		e, _ := r.Err()
		assert.Equal(t, firstErr, e)
	})

	t.Run("empty slice returns Ok empty slice", func(t *testing.T) {
		t.Parallel()
		r := result.Collect([]result.Result[int, error]{})
		require.True(t, r.IsOk())
		assert.Empty(t, r.Unwrap())
	})

	t.Run("preserves order", func(t *testing.T) {
		t.Parallel()
		rs := []result.Result[string, error]{
			result.Ok[string, error]("a"),
			result.Ok[string, error]("b"),
			result.Ok[string, error]("c"),
		}
		r := result.Collect(rs)
		assert.Equal(t, []string{"a", "b", "c"}, r.Unwrap())
	})

	t.Run("result slice has exact length", func(t *testing.T) {
		t.Parallel()

		rs := []result.Result[int, error]{
			result.Ok[int, error](1),
			result.Ok[int, error](2),
		}
		r := result.Collect(rs)
		assert.Len(t, r.Unwrap(), 2)
	})
}

func TestCollectAll(t *testing.T) {
	t.Parallel()

	t.Run("separates values and errors", func(t *testing.T) {
		t.Parallel()
		err1 := errors.New("e1")
		err2 := errors.New("e2")
		rs := []result.Result[int, error]{
			result.Ok[int, error](1),
			result.Err[int](err1),
			result.Ok[int, error](2),
			result.Err[int](err2),
		}
		values, errs := result.CollectAll(rs)
		assert.Equal(t, []int{1, 2}, values)
		assert.Equal(t, []error{err1, err2}, errs)
	})

	t.Run("all Ok returns no errors", func(t *testing.T) {
		t.Parallel()
		rs := []result.Result[int, error]{
			result.Ok[int, error](10),
			result.Ok[int, error](20),
		}
		values, errs := result.CollectAll(rs)
		assert.Equal(t, []int{10, 20}, values)
		assert.Empty(t, errs)
	})

	t.Run("all Err returns no values", func(t *testing.T) {
		t.Parallel()
		rs := []result.Result[int, error]{
			result.Err[int](errors.New("a")),
			result.Err[int](errors.New("b")),
		}
		values, errs := result.CollectAll(rs)
		assert.Empty(t, values)
		assert.Len(t, errs, 2)
	})

	t.Run("empty slice returns nil slices", func(t *testing.T) {
		t.Parallel()
		values, errs := result.CollectAll([]result.Result[int, error]{})
		assert.Nil(t, values)
		assert.Nil(t, errs)
	})
}
