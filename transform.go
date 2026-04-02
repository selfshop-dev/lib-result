package result

// Map transforms the success value with fn, leaving any Err unchanged.
//
// Because Go generics do not support additional type parameters on methods,
// Map is a package-level function rather than a method on Result.
//
//	upper := result.Map(r, strings.ToUpper)
func Map[T, U any, E error](r Result[T, E], f func(T) U) Result[U, E] {
	if r.ok {
		return Ok[U, E](f(r.value))
	}
	return Err[U](r.err)
}

// MapErr transforms the error with fn, leaving any Ok unchanged.
//
//	// Wrap a low-level error in a domain error:
//	r2 := result.MapErr(r, func(e error) *apperr.Error {
//	    return apperr.Wrap(e, apperr.KindNotFound, "user not found")
//	})
func MapErr[T any, E, F error](r Result[T, E], f func(E) F) Result[T, F] {
	if r.ok {
		return Ok[T, F](r.value)
	}
	return Err[T](f(r.err))
}

// AndThen (flatMap / bind) calls fn with the Ok value and returns its Result.
// Short-circuits on Err — fn is never called.
//
//	r := result.Of(repo.FindUser(ctx, id)).
//	    AndThen(func(u User) result.Result[Order, error] {
//	        return result.Of(repo.LatestOrder(ctx, u.ID))
//	    })
func AndThen[T, U any, E error](r Result[T, E], f func(T) Result[U, E]) Result[U, E] {
	if r.ok {
		return f(r.value)
	}
	return Err[U](r.err)
}

// OrElse calls fn with the Err value and returns its Result.
// Short-circuits on Ok — fn is never called.
//
//	// Fall back to a cached value on not-found:
//	r := result.Of(repo.FindUser(ctx, id)).
//	    OrElse(func(e error) result.Result[User, error] {
//	        if apperr.IsKind(e, apperr.KindNotFound) {
//	            return result.Ok[User, error](guestUser)
//	        }
//	        return result.Err[User, error](e)
//	    })
func OrElse[T any, E, F error](r Result[T, E], f func(E) Result[T, F]) Result[T, F] {
	if r.ok {
		return Ok[T, F](r.value)
	}
	return f(r.err)
}

// And returns other if r is Ok, otherwise propagates r's Err retyped to U.
//
//	validated := result.And(parsed, validate(parsed.Unwrap()))
func And[T, U any, E error](r Result[T, E], other Result[U, E]) Result[U, E] {
	if r.ok {
		return other
	}
	return Err[U](r.err)
}

// Or returns r if Ok, otherwise returns other.
//
//	final := result.Or(primary, fallback)
func Or[T any, E, F error](r Result[T, E], other Result[T, F]) Result[T, F] {
	if r.ok {
		return Ok[T, F](r.value)
	}
	return other
}
