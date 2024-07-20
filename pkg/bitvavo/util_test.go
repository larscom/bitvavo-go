package bitvavo

import "testing"

func assert[T comparable](t *testing.T, expected T, actual T) {
	if expected != actual {
		t.Errorf("\nexpected: %v\nactual: %v\n", expected, actual)
	}
}
