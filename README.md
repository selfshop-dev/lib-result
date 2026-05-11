# lib-result

[![CI](https://github.com/selfshop-dev/lib-result/actions/workflows/ci.yml/badge.svg)](https://github.com/selfshop-dev/lib-result/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/selfshop-dev/lib-result/branch/main/graph/badge.svg)](https://codecov.io/gh/selfshop-dev/lib-result)
[![Go Report Card](https://goreportcard.com/badge/github.com/selfshop-dev/lib-result)](https://goreportcard.com/report/github.com/selfshop-dev/lib-result)
[![Go version](https://img.shields.io/github/go-mod/go-version/selfshop-dev/lib-result)](go.mod)
[![License](https://img.shields.io/github/license/selfshop-dev/lib-result)](LICENSE)

`Result[T, E]` for Go — an explicit type for operations that return either a value or an error. No external dependencies. A project by [selfshop-dev](https://github.com/selfshop-dev).

### Installation

```bash
go get -u github.com/selfshop-dev/lib-result
```

## Overview

`lib-result` implements the `Result[T, E]` type, which makes the error path explicit in the signature and allows building transformation chains without repetitive `if err != nil` checks. The idiomatic `(T, error)` pair works well for most cases, but when you need to transform, collect, or chain multiple operations that can each fail, the checks accumulate quickly.

```go
// Without Result
user, err := repo.FindUser(ctx, id)
if err != nil {
    return "", err
}
order, err := repo.LatestOrder(ctx, user.ID)
if err != nil {
    return "", err
}
return order.Reference, nil

// With Result
return result.AndThen(
    result.Of(repo.FindUser(ctx, id)),
    func(u User) result.Result[Order, error] {
        return result.Of(repo.LatestOrder(ctx, u.ID))
    },
).Map(func(o Order) string { return o.Reference }).ToGo()
```

For straightforward service code, the idiomatic `(T, error)` pair is preferred. Use `Result` where chaining genuinely improves readability.

### Quick Start

```go
import "github.com/selfshop-dev/lib-result"

// From idiomatic Go
r := result.Of(repo.FindUser(ctx, id))

// Transform
name := r.Map(func(u User) string { return u.Name }).UnwrapOr("anonymous")

// Back to idiomatic Go
user, err := r.ToGo()
```

## Value[T]

For the common case where the error type is `error`, use the [`Value[T]`](result.go) alias:

```go
func FindUser(ctx context.Context, id int64) result.Value[User]

r := result.Of(repo.FindUser(ctx, id)) // inferred as result.Value[User]
```

`Value[T]` is simply `Result[T, error]`. When you need a concrete error type, use the full form `Result[T, *apperr.Error]`.

## Constructors

A `Result` can be created explicitly or bridged from idiomatic Go:

```go
result.Ok[User, error](user)      // explicit Ok
result.Err[User, error](err)      // explicit Err — panics if err == nil
result.Of(repo.FindUser(ctx, id)) // from an idiomatic Go (value, error) pair

result.OfTyped[User, *apperr.Error](repo.FindUser(ctx, id)) // with a concrete error type
```

`Err` panics on `nil` — including typed nil (`var e *MyError = nil`). A nil error means success; use `Ok`.

## Accessing the Value

Safe methods never panic; unsafe methods are intended only for tests and program initialization.

```go
// Safe — never panic
v, ok := r.Value()          // (T, bool)
e, ok := r.Err()            // (E, bool)
r.UnwrapOr(fallback)        // T or fallback
r.UnwrapOrElse(func(e E) T) // T or result of fn
r.UnwrapOrZero()            // T or zero value

// Unsafe — tests and initialization only
r.Unwrap()            // T or panic
r.Must("load config") // T or panic with context
```

## Transformations

All transformations are package-level functions, since Go does not support additional type parameters on methods.

```go
// Transform the Ok value
result.Map(r, strings.ToUpper)

// Transform the error
result.MapErr(r, func(e error) *apperr.Error {
    return apperr.Wrap(e, apperr.KindNotFound, "user not found")
})

// Chain operations — short-circuits on Err
result.AndThen(r, func(u User) result.Result[Order, error] {
    return result.Of(repo.LatestOrder(ctx, u.ID))
})

// Fallback on error — short-circuits on Ok
result.OrElse(r, func(e error) result.Result[User, error] {
    if apperr.IsKind(e, apperr.KindNotFound) {
        return result.Ok[User, error](guestUser)
    }
    return result.Err[User, error](e)
})

// Combining
result.And(r, other) // other if r is Ok, otherwise Err from r
result.Or(r, other)  // r if Ok, otherwise other
```

## Collections

Two options for working with a slice of results: stop on the first error, or collect all failures.

```go
// Stop on first error
all := result.Collect(results) // Result[[]T, E]

// Collect everything — do not stop on errors
values, errs := result.CollectAll(results) // ([]T, []E)
```

`Collect` is appropriate when the entire batch is meaningless if any error occurs. `CollectAll` is for when you need to process all elements and report every failure.

## Converting Back to Go

For compatibility with the standard library and third-party code, a `Result` can be converted back to idiomatic Go.

```go
user, err := r.ToGo()      // (T, error)
user, err := r.ToGoTyped() // (T, E) — no type assertion needed
value, ok := r.Option()    // (T, bool) — error is discarded
```

## License

[`MIT`](LICENSE) © 2026-present [`selfshop-dev`](https://github.com/selfshop-dev)