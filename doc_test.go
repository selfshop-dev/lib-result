package result_test

import (
	"errors"
	"fmt"
	"strings"

	result "github.com/selfshop-dev/lib-result"
)

func ExampleOf() {
	parse := func(s string) (int, error) {
		if s == "" {
			return 0, errors.New("empty string")
		}
		return len(s), nil
	}

	r := result.Of(parse("hello"))
	fmt.Println(r.IsOk())
	fmt.Println(r.Unwrap())

	// Output:
	// true
	// 5
}

func ExampleOk() {
	r := result.Ok[int, error](42)
	fmt.Println(r.IsOk())
	fmt.Println(r.Unwrap())

	// Output:
	// true
	// 42
}

func ExampleErr() {
	r := result.Err[int](errors.New("not found"))
	fmt.Println(r.IsErr())

	e, ok := r.Err()
	fmt.Println(ok)
	fmt.Println(e)

	// Output:
	// true
	// true
	// not found
}

func ExampleResult_Value() {
	r := result.Ok[string, error]("hello")
	v, ok := r.Value()
	fmt.Println(ok)
	fmt.Println(v)

	// Output:
	// true
	// hello
}

func ExampleResult_UnwrapOr() {
	r := result.Of(0, errors.New("not found"))
	fmt.Println(r.UnwrapOr(-1))

	// Output:
	// -1
}

func ExampleResult_UnwrapOrElse() {
	r := result.Err[int](errors.New("fail"))
	v := r.UnwrapOrElse(func(_ error) int { return -1 })
	fmt.Println(v)

	// Output:
	// -1
}

func ExampleResult_Must() {
	r := result.Ok[string, error]("config loaded")
	fmt.Println(r.Must("load config"))

	// Output:
	// config loaded
}

func ExampleResult_ToGo() {
	v, err := result.Ok[int, error](7).ToGo()
	fmt.Println(v, err)

	// Output:
	// 7 <nil>
}

func ExampleResult_String() {
	fmt.Println(result.Ok[int, error](42))
	fmt.Println(result.Err[int](errors.New("fail")))

	// Output:
	// Ok(42)
	// Err(fail)
}

func ExampleMap() {
	r := result.Map(result.Ok[int, error](5), func(n int) string {
		return fmt.Sprintf("value=%d", n)
	})
	fmt.Println(r.Unwrap())

	// Output:
	// value=5
}

func ExampleMapErr() {
	r := result.MapErr(
		result.Err[int](errors.New("low-level")),
		func(e error) error { return fmt.Errorf("wrapped: %w", e) },
	)
	e, _ := r.Err()
	fmt.Println(e)

	// Output:
	// wrapped: low-level
}

func ExampleAndThen() {
	double := func(n int) result.Result[int, error] {
		return result.Ok[int, error](n * 2)
	}

	r := result.AndThen(result.Of(5, nil), double)
	fmt.Println(r.Unwrap())

	// Output:
	// 10
}

func ExampleOrElse() {
	r := result.OrElse(
		result.Err[string](errors.New("fail")),
		func(_ error) result.Result[string, error] {
			return result.Ok[string, error]("fallback")
		},
	)
	fmt.Println(r.Unwrap())

	// Output:
	// fallback
}

func ExampleCollect() {
	rs := []result.Result[int, error]{
		result.Ok[int, error](1),
		result.Ok[int, error](2),
		result.Ok[int, error](3),
	}
	fmt.Println(result.Collect(rs).Unwrap())

	// Output:
	// [1 2 3]
}

func ExampleCollectAll() {
	rs := []result.Result[string, error]{
		result.Ok[string, error]("a"),
		result.Err[string](errors.New("e1")),
		result.Ok[string, error]("b"),
		result.Err[string](errors.New("e2")),
	}
	values, errs := result.CollectAll(rs)
	fmt.Println(strings.Join(values, ","))
	fmt.Println(len(errs))

	// Output:
	// a,b
	// 2
}

func ExampleValue() {
	// Value[T] is an alias for Result[T, error] — the common case.
	r := result.Of(99, nil)
	fmt.Println(r.Unwrap())

	// Output:
	// 99
}
