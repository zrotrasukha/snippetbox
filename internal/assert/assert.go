// Package assert is dedicated to testing and shit
package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	// to tell that this equal function is a test helper

	t.Helper()

	if actual != expected {
		t.Errorf("got: %v, wanted:%v", actual, expected)
	}
}
