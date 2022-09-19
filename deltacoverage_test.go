package deltacoverage_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/broadwing/deltacoverage"
	"github.com/google/go-cmp/cmp"
)

func TestNewCoverProfile_ErrorsIfPathDoesNotExist(t *testing.T) {
	t.Parallel()
	nonExistentDirPath := t.TempDir() + "/bogus-directory"
	_, err := deltacoverage.NewCoverProfile(nonExistentDirPath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("want error os.ErrNotExist, got %#v", err)
	}
}

func TestNewCoverProfile_ErrorsIfPathIsNotDirectory(t *testing.T) {
	t.Parallel()
	_, err := deltacoverage.NewCoverProfile("testdata/empty-file.txt")
	if !errors.Is(err, deltacoverage.ErrMustBeDirectory) {
		t.Errorf("want error deltacoverage.ErrMustBeDirectory, got %#v", err)
	}
}

func TestParseCoverProfile_ReturnsExpectedCoverProfileGivenCoverProfileDirectory(t *testing.T) {
	t.Parallel()
	want := &deltacoverage.CoverProfile{
		DirPath: "testdata/coverprofiles",
		UniqueBranches: map[string]int{
			"a/a.go:7.30,9.2": 1,
		},
		Tests: map[string][]string{
			"TestSumOnePlusOne":        {"a/a.go:3.24,5.2"},
			"TestSumTwoPlusTwo":        {"a/a.go:3.24,5.2"},
			"TestSubstractTwoMinusTwo": {"a/a.go:7.30,9.2"},
		},
		NumberStatements: 3,
	}
	got, err := deltacoverage.NewCoverProfile("testdata/coverprofiles")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPrintDeltaCoverage_PrintsCorrectPercentDeltaCoverageGivenCoverProfile(t *testing.T) {
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
		NumberStatements: 3,
	}
	want := "TestSubstractTwoMinusTwo 33.3%\nTestSumOnePlusOne 0.0%\nTestSumTwoPlusTwo 0.0%"
	output := &bytes.Buffer{}
	_, err := fmt.Fprint(output, covProfile)
	if err != nil {
		t.Fatal(err)
	}
	got := output.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(strings.Fields(want), strings.Fields(got)))
	}
}

func TestPrintDeltaCoverage_PrintsNoTestsFoundGivenDirectoryWithNoCoverProfile(t *testing.T) {
	t.Parallel()
	c, err := deltacoverage.NewCoverProfile(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	want := "No tests found"
	output := &bytes.Buffer{}
	_, err = fmt.Fprint(output, c)
	if err != nil {
		t.Fatal(err)
	}
	got := output.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
