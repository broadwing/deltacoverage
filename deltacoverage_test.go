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

func TestParseTotalStatements_ErrorsGivenNotFile(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc string
		path string
		err  error
	}{
		{
			desc: "Given directory",
			path: t.TempDir(),
			err:  deltacoverage.ErrMustBeFile,
		},
		{
			desc: "Given file does not exist",
			path: t.TempDir() + "/bogus.file",
			err:  os.ErrNotExist,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := deltacoverage.ParseTotalStatements(tC.path)
			if !errors.Is(err, tC.err) {
				t.Errorf("want error %#v but got %#v", tC.err, err)
			}
		})
	}
}

func TestParseCoverProfile_ReturnsExpectedCoverProfileGivenCoverProfileDirectory(t *testing.T) {
	t.Parallel()
	want := &deltacoverage.CoverProfile{
		DirPath: "testdata/sample",
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
	got, err := deltacoverage.NewCoverProfile("testdata/sample")
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
		t.Error(cmp.Diff(want, got))
	}
}
