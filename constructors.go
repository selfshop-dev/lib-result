package result

// Ok constructs a successful Result containing value.
func Ok[T any, E error](value T) Result[T, E] {
	return Result[T, E]{value: value, ok: true}
}

// Err constructs a failed Result containing err.
// Panics if err is nil — a nil error is success, use [Ok] instead.
func Err[T any, E error](err E) Result[T, E] {
	if isNilError(err) {
		panic("result.Err: err must not be nil — use result.Ok for success")
	}
	return Result[T, E]{err: err}
}

// Of converts the idiomatic Go (value, error) pair into a [Result].
// This is the primary bridge between stdlib / existing code and Result:
//
//	r := result.Of(os.Open("file.txt"))
//	r := result.Of(repo.FindUser(ctx, id))
func Of[T any](value T, err error) Result[T, error] {
	if err != nil {
		return Result[T, error]{err: err}
	}
	return Result[T, error]{value: value, ok: true}
}

// OfTyped is [Of] for callers that want a concrete E rather than error.
// Useful when the error type is a known domain type such as *apperr.Error.
//
//	r := result.OfTyped[User, *apperr.Error](repo.FindUser(ctx, id))
func OfTyped[T any, E error](value T, err E) Result[T, E] {
	if isNilError(err) {
		return Result[T, E]{value: value, ok: true}
	}
	return Result[T, E]{err: err}
}
