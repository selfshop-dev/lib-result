// Package result provides a Result[T, E] type for Go.
//
// Result represents either a successful value ([Ok]) or a failure ([Err]).
// It makes the error path explicit in the type signature and enables
// transformation chains without repetitive if-err blocks.
//
// Zero external dependencies — stdlib only. Requires Go 1.23+.
//
// # When to use Result
//
// Result is most valuable for batch and pipeline operations where multiple
// steps may fail and you want to transform, collect, or chain them:
//
//	// Transform a (value, error) pair into a pipeline:
//	name := result.Of(repo.FindUser(ctx, id)).
//	    Map(func(u User) string { return u.Name }).
//	    UnwrapOr("anonymous")
//
//	// Collect a slice of Results, stopping at the first error:
//	all := result.Collect(results)
//
// For simple service-layer code the idiomatic Go (T, error) pair is
// preferable. Use Result at boundaries where chaining adds real clarity.
//
// # Type parameters
//
//   - T — the success value type (any).
//   - E — the error type; must satisfy the error interface.
//
// For the common case where E is error, use the [Value] alias:
//
//	func FindUser(ctx context.Context, id int64) result.Value[User]
//
//	var r result.Value[User] = result.Of(repo.FindUser(ctx, id))
//
// # Constructors
//
// Build a Result from a raw value or an (value, error) pair:
//
//	result.Ok[User, error](user)         // explicit Ok
//	result.Err[User, error](err)         // explicit Err
//	result.Of(repo.FindUser(ctx, id))    // from idiomatic Go pair
//
// # Accessing the value
//
//	if v, ok := r.Value(); ok { … use v … }
//	user := r.Unwrap()           // panics on Err — tests / init only
//	user := r.UnwrapOr(guest)    // safe fallback
//
// # Converting back to Go idiom
//
//	user, err := r.ToGo()
//
// # Concurrency
//
// Result is effectively immutable after construction — all methods
// are read-only. Safe for concurrent reads; do not share a pointer
// to a Result across goroutines if any goroutine may write.
package result
