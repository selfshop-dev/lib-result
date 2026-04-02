package result

import "fmt"

// IsOk reports whether the Result holds a success value.
func (r Result[T, E]) IsOk() bool { return r.ok }

// IsErr reports whether the Result holds an error.
func (r Result[T, E]) IsErr() bool { return !r.ok }

// Value returns the success value and true, or the zero value of T and false.
//
//	if v, ok := r.Value(); ok {
//	    // use v
//	}
func (r Result[T, E]) Value() (T, bool) {
	return r.value, r.ok
}

// Err returns the error value and true, or the zero value of E and false.
//
//	if e, ok := r.Err(); ok {
//	    return e
//	}
func (r Result[T, E]) Err() (E, bool) {
	return r.err, !r.ok
}

// Unwrap returns the success value.
// Panics with the error if the Result is Err.
//
// Use only when you are certain the Result is Ok — for example in tests or
// program initialisation. For production paths prefer [Result.Value] or
// [Result.UnwrapOr].
func (r Result[T, E]) Unwrap() T {
	if r.ok {
		return r.value
	}
	panic(fmt.Sprintf("result.Unwrap called on Err: %v", r.err))
}

// UnwrapErr returns the error value.
// Panics if the Result is Ok.
func (r Result[T, E]) UnwrapErr() E {
	if r.ok {
		panic("result.UnwrapErr called on Ok")
	}
	return r.err
}

// UnwrapOr returns the success value, or fallback if the Result is Err.
//
//	name := result.Of(repo.FindUser(ctx, id)).
//	    Map(func(u User) string { return u.Name }).
//	    UnwrapOr("anonymous")
func (r Result[T, E]) UnwrapOr(fallback T) T {
	if r.ok {
		return r.value
	}
	return fallback
}

// UnwrapOrElse returns the success value, or calls fn with the error if Err.
// Prefer this over [Result.UnwrapOr] when computing the fallback is expensive.
//
//	count := r.UnwrapOrElse(func(e error) int {
//	    log.Warn("count unavailable", "err", e)
//	    return 0
//	})
func (r Result[T, E]) UnwrapOrElse(f func(E) T) T {
	if r.ok {
		return r.value
	}
	return f(r.err)
}

// UnwrapOrZero returns the success value or the zero value of T.
// Equivalent to [Result.UnwrapOr] with the zero value — saves an allocation
// when the zero value is the correct default.
func (r Result[T, E]) UnwrapOrZero() T {
	return r.value // zero value when !ok
}

// Must returns the success value or panics with a formatted message.
// Intended for program initialisation and tests — never use in request handlers.
//
//	cfg := result.Of(cfgx.Load(&AppConfig{})).Must("load config")
func (r Result[T, E]) Must(context string) T {
	if r.ok {
		return r.value
	}
	panic(fmt.Sprintf("result: %s: %v", context, r.err))
}

// String returns a human-readable representation of the Result.
// Format: "Ok(<value>)" or "Err(<error>)".
func (r Result[T, E]) String() string {
	if r.ok {
		return fmt.Sprintf("Ok(%v)", r.value)
	}
	return fmt.Sprintf("Err(%v)", r.err)
}
