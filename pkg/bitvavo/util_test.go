package bitvavo

import "testing"

func assert(t *testing.T, expected any, actual any) {
	if expected != actual {
		t.Errorf("\nexpected: %v\nactual: %v\n", expected, actual)
	}
}
