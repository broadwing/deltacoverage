package deltacoverage_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/broadwing/deltacoverage"
	"github.com/google/go-cmp/cmp"
)

func TestParseCoverProfile_ErrorsIfPathDoesNotExist(t *testing.T) {
	t.Parallel()
	nonExistentFilePath := t.TempDir() + "/bogus-directory.bogus"
	_, err := deltacoverage.ParseCoverProfile(nonExistentFilePath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("want error os.ErrNotExist, got %#v", err)
	}
}

func TestParseCoverProfile_ErrorsIfPathIsNotDirectory(t *testing.T) {
	t.Parallel()
	_, err := deltacoverage.ParseCoverProfile("testdata/empty-file.txt")
	if !errors.Is(err, deltacoverage.ErrMustBeDirectory) {
		t.Errorf("want error deltacoverage.ErrMustBeDirectory, got %#v", err)
	}
}

func TestParseCoverProfile_ReturnsCorrectCoverProfileGivenCoverProfileDirectory(t *testing.T) {
	t.Parallel()
	want := deltacoverage.CoverProfile{
		Branches: map[string]int{
			"xyz/xyz.go:3.24,5.2": 1,
		},
		Tests: map[string][]string{
			"TestSumOnePlusOne": {},
			"TestSumTwoPlusTwo": {},
		},
		TotalStatements: 1,
	}
	got, err := deltacoverage.ParseCoverProfile("testdata/simple-sample")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPrintDeltaCoverage_PrintsZeroDeltaCoverageGivenNoCoverageForTests(t *testing.T) {
	t.Parallel()
	covProfile := deltacoverage.CoverProfile{
		Tests: map[string][]string{
			"TestSumOnePlusOne": {},
			"TestSumTwoPlusTwo": {},
		},
	}
	want := "TestSumOnePlusOne 0.0%\nTestSumTwoPlusTwo 0.0%"
	output := &bytes.Buffer{}
	_, err := fmt.Fprint(output, covProfile)
	if err != nil {
		t.Fatal(err)
	}
	got := output.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
