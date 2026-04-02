package result

import "testing"

// FuzzIsNilError guards against panics in the reflect-based isNilError when
// given arbitrary byte content. Branch coverage for isNilError is provided by
// the unit tests in constructors_test.go; this fuzzer is a panic safety net.
//
// Run the fuzzer:
//
//	go test -fuzz=FuzzIsNilError -fuzztime=60s
//
// Run the seed corpus only (CI):
//
//	go test -run=FuzzIsNilError
func FuzzIsNilError(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte(""))
	f.Add([]byte("error"))
	f.Add([]byte("\x00"))
	f.Add([]byte("\xff\xfe"))  // invalid UTF-8
	f.Add(make([]byte, 1<<16)) // large input

	f.Fuzz(func(t *testing.T, msg []byte) {
		// Non-nil pointer must never be reported as nil.
		nonNil := &errorStub{msg: string(msg)}
		if isNilError(nonNil) {
			t.Errorf("isNilError(&errorStub{%q}) = true, want false", msg)
		}

		// Typed nil pointer must always be reported as nil.
		var nilPtr *errorStub
		if !isNilError(nilPtr) {
			t.Error("isNilError(typed nil *errorStub) = false, want true")
		}
	})
}

// errorStub is a minimal pointer-receiver error used as the fuzz target type.
type errorStub struct{ msg string }

func (e *errorStub) Error() string { return e.msg }
