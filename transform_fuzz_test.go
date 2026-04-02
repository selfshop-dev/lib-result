package result_test

import (
	"errors"
	"testing"

	result "github.com/selfshop-dev/lib-result"
)

// FuzzAndThen exercises AndThen chains of arbitrary depth against both Ok and
// Err starting values.
//
// Invariants checked on every corpus entry:
//  1. Ok input  → fn is called exactly depth times, result is Ok.
//  2. Ok input  → final value equals depth (each step increments by 1).
//  3. Err input → fn is never called, result is always Err.
//
// Run the fuzzer:
//
//	go test -fuzz=FuzzAndThen -fuzztime=60s
//
// Run the seed corpus only (CI):
//
//	go test -run=FuzzAndThen
func FuzzAndThen(f *testing.F) {
	f.Add(true, uint8(1))   // Ok start, shallow chain.
	f.Add(true, uint8(10))  // Ok start, deeper chain.
	f.Add(true, uint8(255)) // Ok start, maximum depth.
	f.Add(true, uint8(0))   // Ok start, zero steps — fn never called.

	// Err start, fn must never be called regardless of depth.
	f.Add(false, uint8(0))
	f.Add(false, uint8(1))
	f.Add(false, uint8(255))

	f.Fuzz(func(t *testing.T, startOk bool, depth uint8) {
		var r result.Result[int, error]
		if startOk {
			r = result.Ok[int, error](0)
		} else {
			r = result.Err[int](errors.New("initial error"))
		}

		called := 0
		for range depth {
			r = result.AndThen(r, func(n int) result.Result[int, error] {
				called++
				return result.Ok[int, error](n + 1)
			})
		}

		if startOk {
			// Invariant 1: Ok chain must remain Ok.
			if !r.IsOk() {
				t.Fatalf("Ok chain of depth %d returned Err", depth)
			}

			// Invariant 2: fn must be called exactly depth times.
			if called != int(depth) {
				t.Errorf("fn called %d times, want %d", called, depth)
			}

			// Invariant 3: each step increments by 1, so final value == depth.
			if got := r.Unwrap(); got != int(depth) {
				t.Errorf("value = %d, want %d", got, int(depth))
			}
		} else {
			// Invariant 4: Err input must short-circuit — fn never called.
			if !r.IsErr() {
				t.Fatalf("Err input, depth %d: AndThen chain returned Ok", depth)
			}

			if called != 0 {
				t.Errorf("fn called %d times on Err input, want 0", called)
			}
		}
	})
}
