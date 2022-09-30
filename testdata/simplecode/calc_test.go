package calc_test

import (
	"calc"
	"testing"
)

func TestSumOnePlusOne(t *testing.T) {
	t.Parallel()

	want := 2
	got := calc.Sum(1, 1)
	if want != got {
		t.Errorf("want %d got %d", want, got)
	}
}
