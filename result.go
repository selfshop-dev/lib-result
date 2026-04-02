package result

// Result holds either a success value of type T or an error of type E.
//
// The zero value is invalid — always construct via [Ok], [Err], or [Of].
// Result is effectively immutable after construction.
type Result[T any, E error] struct {
	value T
	err   E
	ok    bool
}

// Value is a convenience alias for the common case where E is error.
// Use it to express "this operation returns a T or fails" without
// spelling out the error type parameter every time.
//
//	func FindUser(ctx context.Context, id int64) result.Value[User]
//
//	r := result.Of(repo.FindUser(ctx, id)) // inferred as result.Value[User]
type Value[T any] = Result[T, error]
