package result

// ToGo converts a Result back to the idiomatic Go (value, error) pair.
// Use at the boundary between Result-style code and stdlib / third-party code.
//
//	user, err := result.Of(repo.FindUser(ctx, id)).
//	    Map(normalise).
//	    ToGo()
func (r Result[T, E]) ToGo() (T, error) {
	if r.ok {
		return r.value, nil
	}
	return r.value, r.err
}

// ToGoTyped is [Result.ToGo] but returns E instead of error.
// Use when the caller needs the concrete error type without a type assertion.
func (r Result[T, E]) ToGoTyped() (T, E) {
	return r.value, r.err
}

// Option converts a Result to (value, bool), discarding the error.
// Useful when you only care whether there is a value, not why there is not.
func (r Result[T, E]) Option() (T, bool) {
	return r.value, r.ok
}
