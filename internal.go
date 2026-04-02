package result

import "reflect"

// isNilError reports whether e, despite satisfying the error interface,
// holds a nil pointer or nil interface. This catches both:
//
//   - plain nil interface:     var e error = nil     → true
//   - typed nil pointer:       var e *MyError = nil  → true
//   - valid non-nil error:     errors.New("x")       → false
func isNilError[E error](e E) bool {
	av := any(e)
	if av == nil {
		return true
	}
	// Detect typed nil: the interface box is non-nil but holds a nil pointer.
	v := reflect.ValueOf(av)
	switch v.Kind() {
	case reflect.Pointer, reflect.Chan, reflect.Map,
		reflect.Interface, reflect.Func, reflect.Slice:
		return v.IsNil()
	}
	return false
}
