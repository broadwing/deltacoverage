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
	nonExistentDirPath := t.TempDir() + "/bogus-directory"
	_, err := deltacoverage.ParseCoverProfile(nonExistentDirPath)
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

func TestParseCoverProfile_ReturnsExpectedCoverProfileGivenCoverProfileDirectory(t *testing.T) {
	t.Parallel()
	want := deltacoverage.CoverProfile{
		UniqueBranches: map[string]int{
			"a/a.go:7.30,9.2": 1,
		},
		Tests: map[string][]string{
			"TestSumOnePlusOne":        {"a/a.go:3.24,5.2"},
			"TestSumTwoPlusTwo":        {"a/a.go:3.24,5.2"},
			"TestSubstractTwoMinusTwo": {"a/a.go:7.30,9.2"},
		},
		TotalStatements: 2,
	}
	got, err := deltacoverage.ParseCoverProfile("testdata/sample")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseCoverProfile_ReturnsExpectedCoverProfileGivenNewCoverProfileDirectory(t *testing.T) {
	t.Parallel()
	want := deltacoverage.CoverProfile{
		UniqueBranches: map[string]int{
			"a/a.go:7.30,9.2": 1,
		},
		Tests: map[string][]string{
			"TestSumOnePlusOne":        {"a/a.go:3.24,5.2"},
			"TestSumTwoPlusTwo":        {"a/a.go:3.24,5.2"},
			"TestSubstractTwoMinusTwo": {"a/a.go:7.30,9.2"},
		},
		TotalStatements: 3,
	}
	got, err := deltacoverage.ParseCoverProfile("testdata/new-sample")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPrintDeltaCoverage_PrintsFiftyPercentDeltaCoverageGivenThreeTestsAndOneUniqueBranch(t *testing.T) {
	t.Parallel()
	covProfile := deltacoverage.CoverProfile{
		UniqueBranches: map[string]int{
			"a/a.go:7.30,9.2": 1,
		},
		Tests: map[string][]string{
			"TestSumOnePlusOne":        {"a/a.go:3.24,5.2"},
			"TestSumTwoPlusTwo":        {"a/a.go:3.24,5.2"},
			"TestSubstractTwoMinusTwo": {"a/a.go:7.30,9.2"},
		},
		TotalStatements: 2,
	}
	want := "TestSubstractTwoMinusTwo 50.0%\nTestSumOnePlusOne 0.0%\nTestSumTwoPlusTwo 0.0%"
	output := &bytes.Buffer{}
	_, err := fmt.Fprint(output, covProfile)
	if err != nil {
		t.Fatal(err)
	}
	got := output.String()
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}
