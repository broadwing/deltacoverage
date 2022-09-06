package deltacoverage_test

import (
	"bytes"
	"testing"

	"github.com/broadwing/deltacoverage"
	"github.com/google/go-cmp/cmp"
)

func TestInstrumentErrorsIfProgramDoesNotStartWithPackage(t *testing.T) {
	t.Parallel()
	output := &bytes.Buffer{}
	err := deltacoverage.Instrument(output, `bogus`)
	if err == nil {
		t.Error("want error but not found")
	}
}

func TestInstrumentDoNothingWithNoBodyFunction(t *testing.T) {
	t.Parallel()
	output := &bytes.Buffer{}
	err := deltacoverage.Instrument(output, `package nobody

func nobody()`)
	if err != nil {
		t.Fatal(err)
	}
	want := `package nobody

func nobody()
`
	got := output.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInstrumentEmptyFunction(t *testing.T) {
	t.Parallel()
	output := &bytes.Buffer{}
	err := deltacoverage.Instrument(output, `package empty

func empty() {
}`)
	if err != nil {
		t.Fatal(err)
	}
	want := `package empty

import "fmt"

func empty() {
	fmt.Println("instrumentation")
}
`
	got := output.String()
	t.Fatal(got)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
