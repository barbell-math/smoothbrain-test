// A very simple library that helps with assertions in unit tests.
package smoothbraintest

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"testing"
)

// Formats an error and calls the `t.Fatal` to stop any further execution of the
// unit test. The error will have the following format:
//
//	Error | File <file> Line #### | <message>
//	Expected: (<type>) <value>
//	Got:      (<type>) <value>
func FormatError(
	t *testing.T,
	expected any,
	got any,
	base string,
	file string,
	line int,
) {
	t.Fatal(fmt.Sprintf(
		"Error | File %s Line %d | %s\nExpected: (%T) '%v'\nGot     : (%T) '%v'",
		file, line, base, expected, expected, got, got,
	))
}

// Tests that the expected error is present in the given error.
func ContainsError(t *testing.T, expected error, got error) {
	if !errors.Is(got, expected) {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			t,
			expected,
			got,
			"The expected error was not contained in the given error.",
			f, line,
		)
	}
}

// Tests that the supplied action results in a panic. The panic is recovered so
// all future unit tests will still run.
func Panics(t *testing.T, action func()) {
	defer func() {
		if r := recover(); r == nil {
			_, f, line, _ := runtime.Caller(1)
			FormatError(
				t,
				"panic", "",
				"The supplied funciton did not panic when it should have.",
				f, line,
			)
		}
	}()
	action()
}

// Tests that the supplied action does not result in a panic. Any panic that
// does occur is recovered so all future unit tests will still run.
func NoPanic(t *testing.T, action func()) {
	defer func() {
		if r := recover(); r != nil {
			_, f, line, _ := runtime.Caller(1)
			FormatError(
				t,
				"", "panic",
				"The supplied funciton paniced when it shouldn't have.",
				f, line,
			)
		}
	}()
	action()
}

// Tests that the supplied values are equal. For equality rules refer to the
// language reference: https://go.dev/ref/spec#Comparison_operators
func Eq[T comparable](t *testing.T, expected T, got T) {
	if expected != got {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			t, expected, got,
			"The supplied values were not equal but were expected to be.",
			f, line,
		)
	}
}

// Tests that the expected value is present in the supplied slice. For equality
// rules refer to the language reference:
// https://go.dev/ref/spec#Comparison_operators
func EqOneOf[T comparable](t *testing.T, expected T, data []T) {
	for _, rVal := range data {
		if expected == rVal {
			return
		}
	}
	_, f, line, _ := runtime.Caller(1)
	FormatError(
		t, expected, data,
		"The supplied value is not in the supplied slice.",
		f, line,
	)
}

// Tests that the given float is within +/- eps distance of the expected float.
func EqFloat[T ~float32 | float64](t *testing.T, expected T, got T, eps T) {
	if math.Abs(float64(expected-got)) > float64(eps) {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			t, expected, got,
			fmt.Sprintf(
				"The supplied float was not within the expected range of %e to be considered equal.",
				eps,
			),
			f, line,
		)
	}
}

// Tests that the given value is equal to the expected value using the supplied
// comparison function to determine equality.
func EqFunc[T any](t *testing.T, expected T, got T, cmp func(l T, r T) bool) {
	if !cmp(expected, got) {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			t, expected, got,
			"The supplied values were not equal as defined by the supplied comparison function but were expected to be.",
			f, line,
		)
	}
}

// Tests that the supplied values are not equal. For equality rules refer to the
// language reference: https://go.dev/ref/spec#Comparison_operators
func Neq[T comparable](t *testing.T, expected any, got any) {
	if expected == got {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			t, expected, got,
			"The supplied values were equal but were expected to not be.",
			f, line,
		)
	}
}

// Tests that the supplied value is true. This is useful for validating that
// expressions that evaluate to a boolean. This should not be used for equality
// comparisons such as `True(t, 5==5)`. For equality comparisons refer to one
// of the Eq* functions defined in this file.
func True(t *testing.T, v bool) {
	if v != true {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			t, true, v,
			"The supplied value was not true when it was expected to be.",
			f, line,
		)
	}
}

// Tests that the supplied value is false. This is useful for validating that
// expressions that evaluate to a boolean. This should not be used for equality
// comparisons such as `False(t, 6!=5)`. For equality comparisons refer to one
// of the Eq* functions defined in this file.
func False(t *testing.T, v bool) {
	if v != false {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			t, false, v,
			"The supplied value was not false when it was expected to be.",
			f, line,
		)
	}
}

// Tests that the supplied value is nil. `nil` slices, maps, pointers, and
// interfaces are considered to be nil and will pass this test.
func Nil(t *testing.T, v any) {
	// The actual value is nil
	if v == nil {
		return
	}

	// The value is not nil but it's underlying data is nil (i.e. slice, map, pntr)
	rv := reflect.ValueOf(v)
	if rv.IsZero() {
		return
	}

	// The value is an interface that is nil
	tv := reflect.TypeOf(v)
	if tv.Kind() == reflect.Interface && rv.Elem().IsZero() {
		return
	}

	_, f, line, _ := runtime.Caller(1)
	FormatError(
		t, nil, v,
		"The supplied value was not nil when it was expected to be.",
		f, line,
	)
}

// Tests that the supplied value is not nil. `nil` slices, maps, pointers, and
// interfaces are considered to be nil and will fail this test.
func NotNil(t *testing.T, v any) {
	var rv reflect.Value
	var tv reflect.Type

	// The actual value is nil
	if v == nil {
		goto fail
	}

	// The value is not nil but it's underlying data is nil (i.e. slice, map, pntr)
	rv = reflect.ValueOf(v)
	if rv.IsZero() {
		goto fail
	}

	// The value is an interface that is nil
	tv = reflect.TypeOf(v)
	if tv.Kind() == reflect.Interface && !rv.Elem().IsZero() {
		goto fail
	}

	return

fail:
	_, f, line, _ := runtime.Caller(1)
	FormatError(
		t, "!nil", v,
		"The supplied value was nil when it was not expected to be.",
		f, line,
	)
}

// Tests that the supplied slices match. In order for the slices to match they
// must be the same length and values in the same index must compare equal. For
// equality rules refer to the language reference:
// https://go.dev/ref/spec#Comparison_operators
func SlicesMatch[T comparable](t *testing.T, expected []T, got []T) {
	_, f, line, _ := runtime.Caller(1)
	if len(expected) != len(got) {
		FormatError(
			t, len(expected), len(got),
			"Slices do not match in length.",
			f, line,
		)
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			FormatError(
				t, expected[i], got[i],
				fmt.Sprintf("Values do not match | Index: %d", i),
				f, line,
			)
		}
	}
}

// Tests that the supplied slices match in length and content but not in order.
// For equality rules refer to the language reference:
// https://go.dev/ref/spec#Comparison_operators
func SlicesMatchUnordered[T comparable](t *testing.T, expected []T, got []T) {
	_, f, line, _ := runtime.Caller(1)
	if len(expected) != len(got) {
		FormatError(
			t, len(expected), len(got),
			"Slices do not match in length.",
			f, line,
		)
	}

	usedIndexes := map[int]struct{}{}
	for i := 0; i < len(expected); i++ {
		found := false
		for j := 0; j < len(got) && !found; j++ {
			// The values at indexes i and j matched
			if expected[i] == got[j] {
				if _, ok := usedIndexes[j]; !ok {
					found = true
					usedIndexes[j] = struct{}{}
				}
			}
		}
		if !found {
			FormatError(
				t, expected, got[i],
				fmt.Sprintf("Slice value was not accounted for | Index: %d", i),
				f, line,
			)
		}
	}

	if len(usedIndexes) != len(expected) {
		FormatError(
			t, expected, got,
			"The slices were not found to have equivalent elements.",
			f, line,
		)
	}
}

// Tests that the supplied maps match in length and content. For equality rules
// refer to the language reference: https://go.dev/ref/spec#Comparison_operators
func MapsMatch[K comparable, V any](
	t *testing.T,
	expected map[K]V,
	got map[K]V,
) {
	_, f, line, _ := runtime.Caller(1)
	if len(expected) != len(got) {
		FormatError(
			t, len(expected), len(got),
			"Maps do not match in length.",
			f, line,
		)
	}

	for k, v := range expected {
		gotV, ok := got[k]
		if !ok {
			FormatError(
				t, true, ok,
				fmt.Sprintf("A key was not found | Key: %v", k),
				f, line,
			)
		}
		if any(gotV) != any(v) {
			FormatError(
				t, gotV, v,
				"The values stored in the map did not match.",
				f, line,
			)
		}
	}
}
