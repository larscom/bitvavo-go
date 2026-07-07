package util

import "testing"

func TestOrZeroNil(t *testing.T) {
	var p *int
	if got := OrZero(p); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestOrZeroValue(t *testing.T) {
	v := 5
	if got := OrZero(&v); got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}
