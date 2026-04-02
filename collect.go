package result

// Collect turns a []Result[T, E] into a Result[[]T, E].
// Returns the first Err encountered, or Ok with all values in order.
//
//	results := []result.Result[User, error]{r1, r2, r3}
//	all := result.Collect(results)
//	// Ok([]User{…}) or the first Err
func Collect[T any, E error](rs []Result[T, E]) Result[[]T, E] {
	out := make([]T, 0, len(rs))
	for _, r := range rs {
		if !r.ok {
			return Err[[]T](r.err)
		}
		out = append(out, r.value)
	}
	return Ok[[]T, E](out)
}

// CollectAll turns a []Result[T, E] into ([]T, []E).
// Unlike [Collect], it does not short-circuit — it accumulates all errors.
// Useful for batch operations where you want all failures, not just the first.
//
//	values, errs := result.CollectAll(results)
func CollectAll[T any, E error](rs []Result[T, E]) ([]T, []E) {
	var (
		vals []T
		errs []E
	)
	for _, r := range rs {
		if !r.ok {
			errs = append(errs, r.err)
		} else {
			vals = append(vals, r.value)
		}
	}
	return vals, errs
}
