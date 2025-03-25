// This is a package that defines helper functions for writing unit tests.
package test

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func FormatError(
	expected any,
	got any,
	base string,
	file string,
	line int,
	t *testing.T,
) {
	t.Fatal(fmt.Sprintf(
		"Error | File %s Line %d | %s\nExpected: '%v'\nGot     : '%v'",
		file, line, base, expected, got,
	))
}

func ContainsError(expected error, got error, t *testing.T) {
	if !errors.Is(got, expected) {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			expected,
			got,
			"The expected error was not contained in the given error.",
			f, line, t,
		)
	}
}

func Panics(action func(), t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			_, f, line, _ := runtime.Caller(1)
			FormatError(
				"panic", "",
				"The supplied funciton did not panic when it should have.",
				f, line, t,
			)
		}
	}()
	action()
}

func NoPanic(action func(), t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			_, f, line, _ := runtime.Caller(1)
			FormatError(
				"", "panic",
				"The supplied funciton paniced when it shouldn't have.",
				f, line, t,
			)
		}
	}()
	action()
}

func Eq(l any, r any, t *testing.T) {
	if r != l {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			l, r,
			"The supplied values were not equal but were expected to be.",
			f, line, t,
		)
	}
}

func EqOneOf(r []any, l any, t *testing.T) {
	for _, rVal := range r {
		if l == rVal {
			return
		}
	}
	_, f, line, _ := runtime.Caller(1)
	FormatError(
		l, r,
		"The supplied value is not in the supplied slice.",
		f, line, t,
	)
}

func FloatEq[T ~float32 | float64](l T, r T, eps T, t *testing.T) {
	if l-r > eps {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			l, r,
			fmt.Sprintf(
				"The supplied float was not within the expected range of %e to be considered equal.",
				eps,
			),
			f, line, t,
		)
	}
}

func Neq(l any, r any, t *testing.T) {
	if r == l {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			l, r,
			"The supplied values were equal but were expected to not be.",
			f, line, t,
		)
	}
}

func True(v bool, t *testing.T) {
	if v != true {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			true, v,
			"The supplied value was not true when it was expected to be.",
			f, line, t,
		)
	}
}

func False(v bool, t *testing.T) {
	if v != false {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			false, v,
			"The supplied value was not false when it was expected to be.",
			f, line, t,
		)
	}
}

func Nil(v any, t *testing.T) {
	if v != nil {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			nil, v,
			"The supplied value was not nil when it was expected to be.",
			f, line, t,
		)
	}
}

func NilSlice[T any](v []T, t *testing.T) {
	if v != nil {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			nil, v,
			"The supplied slice was not nil when it was expected to be.",
			f, line, t,
		)
	}
}

func NilPntr[T any](v *T, t *testing.T) {
	if v != (*T)(nil) {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			nil, v,
			"The supplied value was not nil when it was expected to be.",
			f, line, t,
		)
	}
}

func NotNil(v any, t *testing.T) {
	if v == nil {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			"!nil", v,
			"The supplied value was nil when it was not expected to be.",
			f, line, t,
		)
	}
}

func NotNilSlice[T any](v []T, t *testing.T) {
	if v == nil {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			nil, v,
			"The supplied slice was nil when it was not expected to be.",
			f, line, t,
		)
	}
}

func NotNilPntr[T any](v *T, t *testing.T) {
	if v == (*T)(nil) {
		_, f, line, _ := runtime.Caller(1)
		FormatError(
			nil, v,
			"The supplied value was not nil when it was expected to be.",
			f, line, t,
		)
	}
}

func SlicesMatch[T any](expected []T, got []T, t *testing.T) {
	_, f, line, _ := runtime.Caller(1)
	if len(expected) != len(got) {
		FormatError(
			len(expected),
			len(got),
			"Slices do not match in length.",
			f, line, t,
		)
	}
	for i := 0; i < len(expected); i++ {
		if any(expected[i]) != any(got[i]) {
			FormatError(
				expected[i],
				got[i],
				fmt.Sprintf("Values do not match | Index: %d", i),
				f, line, t,
			)
		}
	}
}

func SlicesMatchUnordered[T any](expected []T, got []T, t *testing.T) {
	_, f, line, _ := runtime.Caller(1)
	if len(expected) != len(got) {
		FormatError(
			len(expected),
			len(got),
			"Slices do not match in length.",
			f, line, t,
		)
	}

	usedIndexes := map[int]struct{}{}
	for i := 0; i < len(expected); i++ {
		found := false
		for j := 0; j < len(got) && !found; j++ {
			// The values at indexes i and j matched
			if any(expected[i]) == any(got[j]) {
				if _, ok := usedIndexes[j]; !ok {
					found = true
					usedIndexes[j] = struct{}{}
				}
			}
		}
		if !found {
			FormatError(
				expected,
				got[i],
				fmt.Sprintf("Slice value was not accounted for | Index: %d", i),
				f, line, t,
			)
		}
	}

	if len(usedIndexes) != len(expected) {
		FormatError(
			expected,
			got,
			"The slices were not found to have equivalent elements.",
			f, line, t,
		)
	}
}

func MapsMatch[K comparable, V any](
	expected map[K]V,
	got map[K]V,
	t *testing.T,
) {
	_, f, line, _ := runtime.Caller(1)
	if len(expected) != len(got) {
		FormatError(
			len(expected),
			len(got),
			"Maps do not match in length.",
			f, line, t,
		)
	}

	for k, v := range expected {
		gotV, ok := got[k]
		if !ok {
			FormatError(
				true,
				ok,
				fmt.Sprintf("A key was not found | Key: %v", k),
				f, line, t,
			)
		}
		if any(gotV) != any(v) {
			FormatError(
				gotV,
				v,
				"The values stored in the map did not match.",
				f, line, t,
			)
		}
	}
}
