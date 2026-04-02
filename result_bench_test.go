package result_test

import (
	"errors"
	"testing"

	result "github.com/selfshop-dev/lib-result"
)

var sink int // prevents compiler from optimizing away the benchmark body

func BenchmarkResultChain(b *testing.B) {
	b.Run("result.Of.Map.UnwrapOr", func(b *testing.B) {
		for range b.N {
			sink = result.Of(compute()).UnwrapOr(0)
		}
	})

	b.Run("idiomatic if err != nil", func(b *testing.B) {
		for range b.N {
			v, err := compute()
			if err != nil {
				v = 0
			}
			sink = v
		}
	})
}

func compute() (int, error) {
	return 42, errors.New("intentional")
}
