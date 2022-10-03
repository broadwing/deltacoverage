package deltacoverage_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
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
	_, err := deltacoverage.NewCoverProfile("testdata/emptyfile.txt")
	if !errors.Is(err, deltacoverage.ErrMustBeDirectory) {
		t.Errorf("want error deltacoverage.ErrMustBeDirectory, got %#v", err)
	}
}

func TestNewCoverProfile_SetsCorrectPackagePathGivenString(t *testing.T) {
	t.Parallel()
	packagePath := t.TempDir()
	c, err := deltacoverage.NewCoverProfile(packagePath)
	if err != nil {
		t.Error(err)
	}
	want := packagePath
	got := c.PackagePath
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestNewCoverProfile_SetsNonEmptyOutputPath(t *testing.T) {
	t.Parallel()
	packagePath := t.TempDir()
	c, err := deltacoverage.NewCoverProfile(packagePath)
	if err != nil {
		t.Error(err)
	}
	if c.OutputPath == "" {
		t.Error("want non empty output path")
	}
}

func TestParse_ReturnsExpectedUniqueBranchesGivenOutputPath(t *testing.T) {
	testCases := []struct {
		desc       string
		outputPath string
		want       map[string]int
	}{
		{
			desc:       "With simple cover profiles",
			outputPath: "testdata/coverprofiles/simple",
			want:       map[string]int{"a/a.go:7.30,9.2": 1},
		},
		{
			desc:       "With cover profiles sample",
			outputPath: "testdata/coverprofiles/sample",
			want: map[string]int{
				"github.com/broadwing/deltacoverage/deltacoverage.go:104.4,106.23":  3,
				"github.com/broadwing/deltacoverage/deltacoverage.go:106.23,109.19": 3,
				"github.com/broadwing/deltacoverage/deltacoverage.go:113.26,115.6":  1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:113.5,113.26":  1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:116.29,117.14": 1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:116.5,116.29":  1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:119.5,120.12":  2,
				"github.com/broadwing/deltacoverage/deltacoverage.go:120.12,122.6":  1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:123.5,124.70":  2,
				"github.com/broadwing/deltacoverage/deltacoverage.go:126.4,126.40":  1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:131.43,132.16": 1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:132.16,134.4":  1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:141.16,143.3":  1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:144.19,146.3":  1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:35.75,36.38":   1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:36.38,38.3":    1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:39.2,45.16":    5,
				"github.com/broadwing/deltacoverage/deltacoverage.go:48.2,50.16":    3,
				"github.com/broadwing/deltacoverage/deltacoverage.go:53.2,53.22":    1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:53.22,55.3":    1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:56.2,56.22":    1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:62.23,64.3":    1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:65.2,65.37":    1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:65.37,67.26":   2,
				"github.com/broadwing/deltacoverage/deltacoverage.go:67.26,69.14":   2,
				"github.com/broadwing/deltacoverage/deltacoverage.go:69.14,70.13":   1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:72.4,72.66":    1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:74.3,74.73":    1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:76.2,77.35":    2,
				"github.com/broadwing/deltacoverage/deltacoverage.go:96.29,97.51":   1,
				"github.com/broadwing/deltacoverage/deltacoverage.go:97.51,101.18":  4,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c, err := deltacoverage.NewCoverProfile(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			c.OutputPath = tC.outputPath
			err = c.Parse()
			if err != nil {
				t.Fatal(err)
			}
			got := c.UniqueBranches
			if !cmp.Equal(tC.want, got) {
				t.Error(cmp.Diff(tC.want, got))
			}
		})
	}
}

func TestParse_ReturnsExpectedTestsBranchesGivenOutputPath(t *testing.T) {
	testCases := []struct {
		desc       string
		outputPath string
		want       map[string][]string
	}{
		{
			desc:       "With simple cover profiles",
			outputPath: "testdata/coverprofiles/simple",
			want: map[string][]string{
				"TestSubstractTwoMinusTwo": {"a/a.go:7.30,9.2"},
				"TestSumOnePlusOne":        {"a/a.go:3.24,5.2"},
				"TestSumTwoPlusTwo":        {"a/a.go:3.24,5.2"},
			},
		},
		{
			desc:       "With cover profiles sample",
			outputPath: "testdata/coverprofiles/sample",
			want: map[string][]string{
				"TestNewCoverProfile_ErrorsIfPathDoesNotExist": {
					"github.com/broadwing/deltacoverage/deltacoverage.go:139.61,141.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:141.16,143.3",
				},
				"TestNewCoverProfile_ErrorsIfPathIsNotDirectory": {
					"github.com/broadwing/deltacoverage/deltacoverage.go:139.61,141.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:144.2,144.19",
					"github.com/broadwing/deltacoverage/deltacoverage.go:144.19,146.3",
				},
				"TestParseCoverProfile_ReturnsExpectedCoverProfileGivenCoverProfileDirectory": {
					"github.com/broadwing/deltacoverage/deltacoverage.go:35.75,36.38",
					"github.com/broadwing/deltacoverage/deltacoverage.go:39.2,45.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:48.2,50.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:53.2,53.22",
					"github.com/broadwing/deltacoverage/deltacoverage.go:56.2,56.22",
					"github.com/broadwing/deltacoverage/deltacoverage.go:36.38,38.3",
					"github.com/broadwing/deltacoverage/deltacoverage.go:53.22,55.3",
					"github.com/broadwing/deltacoverage/deltacoverage.go:88.50,93.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:96.2,96.29",
					"github.com/broadwing/deltacoverage/deltacoverage.go:131.2,131.43",
					"github.com/broadwing/deltacoverage/deltacoverage.go:136.2,136.12",
					"github.com/broadwing/deltacoverage/deltacoverage.go:96.29,97.51",
					"github.com/broadwing/deltacoverage/deltacoverage.go:97.51,101.18",
					"github.com/broadwing/deltacoverage/deltacoverage.go:104.4,106.23",
					"github.com/broadwing/deltacoverage/deltacoverage.go:126.4,126.40",
					"github.com/broadwing/deltacoverage/deltacoverage.go:106.23,109.19",
					"github.com/broadwing/deltacoverage/deltacoverage.go:113.5,113.26",
					"github.com/broadwing/deltacoverage/deltacoverage.go:116.5,116.29",
					"github.com/broadwing/deltacoverage/deltacoverage.go:119.5,120.12",
					"github.com/broadwing/deltacoverage/deltacoverage.go:123.5,124.70",
					"github.com/broadwing/deltacoverage/deltacoverage.go:113.26,115.6",
					"github.com/broadwing/deltacoverage/deltacoverage.go:116.29,117.14",
					"github.com/broadwing/deltacoverage/deltacoverage.go:120.12,122.6",
					"github.com/broadwing/deltacoverage/deltacoverage.go:131.43,132.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:132.16,134.4",
					"github.com/broadwing/deltacoverage/deltacoverage.go:139.61,141.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:144.2,144.19",
					"github.com/broadwing/deltacoverage/deltacoverage.go:147.2,153.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:156.2,156.21",
				},
				"TestPrintDeltaCoverage_PrintsCorrectPercentDeltaCoverageGivenCoverProfile": {
					"github.com/broadwing/deltacoverage/deltacoverage.go:60.39,62.23",
					"github.com/broadwing/deltacoverage/deltacoverage.go:65.2,65.37",
					"github.com/broadwing/deltacoverage/deltacoverage.go:76.2,77.35",
					"github.com/broadwing/deltacoverage/deltacoverage.go:65.37,67.26",
					"github.com/broadwing/deltacoverage/deltacoverage.go:74.3,74.73",
					"github.com/broadwing/deltacoverage/deltacoverage.go:67.26,69.14",
					"github.com/broadwing/deltacoverage/deltacoverage.go:72.4,72.66",
					"github.com/broadwing/deltacoverage/deltacoverage.go:69.14,70.13",
				},
				"TestPrintDeltaCoverage_PrintsNoTestsFoundGivenDirectoryWithNoCoverProfile": {
					"github.com/broadwing/deltacoverage/deltacoverage.go:60.39,62.23",
					"github.com/broadwing/deltacoverage/deltacoverage.go:62.23,64.3",
					"github.com/broadwing/deltacoverage/deltacoverage.go:88.50,93.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:96.2,96.29",
					"github.com/broadwing/deltacoverage/deltacoverage.go:131.2,131.43",
					"github.com/broadwing/deltacoverage/deltacoverage.go:136.2,136.12",
					"github.com/broadwing/deltacoverage/deltacoverage.go:139.61,141.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:144.2,144.19",
					"github.com/broadwing/deltacoverage/deltacoverage.go:147.2,153.16",
					"github.com/broadwing/deltacoverage/deltacoverage.go:156.2,156.21",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c, err := deltacoverage.NewCoverProfile(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			c.OutputPath = tC.outputPath
			err = c.Parse()
			if err != nil {
				t.Fatal(err)
			}
			got := c.TestsBranches
			if !cmp.Equal(tC.want, got) {
				t.Error(cmp.Diff(tC.want, got))
			}
		})
	}
}

func TestParse_ReturnsExpectedNumberStatementsGivenOutputPath(t *testing.T) {
	t.Parallel()
	c, err := deltacoverage.NewCoverProfile(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.OutputPath = "testdata/coverprofiles/sample"
	err = c.Parse()
	if err != nil {
		t.Fatal(err)
	}
	want := 73
	got := c.NumberStatements
	if want != got {
		t.Errorf("want number of statements %d, got %d", want, got)
	}
}

func TestParse_ErrorsGivenOutputPathWithCoverProfileContainsInvalidNumberStatements(t *testing.T) {
	t.Parallel()
	c, err := deltacoverage.NewCoverProfile(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.OutputPath = "testdata/coverprofiles/invalidNumberStatements"
	err = c.Parse()
	if err == nil {
		t.Error("want error but got nil")
	}
}

func TestParse_ErrorsGivenOutputPathWithCoverProfileContainsInvalidVisited(t *testing.T) {
	t.Parallel()
	c, err := deltacoverage.NewCoverProfile(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.OutputPath = "testdata/coverprofiles/invalidVisited"
	err = c.Parse()
	if err == nil {
		t.Error("want error but got nil")
	}
}

func TestString_ReturnsExpectedDeltaCoverageGivenCoverProfile(t *testing.T) {
	t.Parallel()
	c := deltacoverage.CoverProfile{
		UniqueBranches: map[string]int{
			"github.com/broadwing/deltacoverage/deltacoverage.go:104.4,106.23":  3,
			"github.com/broadwing/deltacoverage/deltacoverage.go:106.23,109.19": 3,
			"github.com/broadwing/deltacoverage/deltacoverage.go:113.26,115.6":  1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:113.5,113.26":  1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:116.29,117.14": 1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:116.5,116.29":  1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:119.5,120.12":  2,
			"github.com/broadwing/deltacoverage/deltacoverage.go:120.12,122.6":  1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:123.5,124.70":  2,
			"github.com/broadwing/deltacoverage/deltacoverage.go:126.4,126.40":  1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:131.43,132.16": 1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:132.16,134.4":  1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:141.16,143.3":  1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:144.19,146.3":  1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:35.75,36.38":   1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:36.38,38.3":    1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:39.2,45.16":    5,
			"github.com/broadwing/deltacoverage/deltacoverage.go:48.2,50.16":    3,
			"github.com/broadwing/deltacoverage/deltacoverage.go:53.2,53.22":    1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:53.22,55.3":    1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:56.2,56.22":    1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:62.23,64.3":    1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:65.2,65.37":    1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:65.37,67.26":   2,
			"github.com/broadwing/deltacoverage/deltacoverage.go:67.26,69.14":   2,
			"github.com/broadwing/deltacoverage/deltacoverage.go:69.14,70.13":   1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:72.4,72.66":    1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:74.3,74.73":    1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:76.2,77.35":    2,
			"github.com/broadwing/deltacoverage/deltacoverage.go:96.29,97.51":   1,
			"github.com/broadwing/deltacoverage/deltacoverage.go:97.51,101.18":  4,
		},
		Tests: []string{
			"TestNewCoverProfile_ErrorsIfPathDoesNotExist",
			"TestNewCoverProfile_ErrorsIfPathIsNotDirectory",
			"TestParseCoverProfile_ReturnsExpectedCoverProfileGivenCoverProfileDirectory",
			"TestPrintDeltaCoverage_PrintsNoTestsFoundGivenDirectoryWithNoCoverProfile",
			"TestPrintDeltaCoverage_PrintsNoTestsFoundGivenDirectoryWithNoCoverProfile",
		},
		TestsBranches: map[string][]string{
			"TestNewCoverProfile_ErrorsIfPathDoesNotExist": {
				"github.com/broadwing/deltacoverage/deltacoverage.go:139.61,141.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:141.16,143.3",
			},
			"TestNewCoverProfile_ErrorsIfPathIsNotDirectory": {
				"github.com/broadwing/deltacoverage/deltacoverage.go:139.61,141.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:144.2,144.19",
				"github.com/broadwing/deltacoverage/deltacoverage.go:144.19,146.3",
			},
			"TestParseCoverProfile_ReturnsExpectedCoverProfileGivenCoverProfileDirectory": {
				"github.com/broadwing/deltacoverage/deltacoverage.go:35.75,36.38",
				"github.com/broadwing/deltacoverage/deltacoverage.go:39.2,45.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:48.2,50.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:53.2,53.22",
				"github.com/broadwing/deltacoverage/deltacoverage.go:56.2,56.22",
				"github.com/broadwing/deltacoverage/deltacoverage.go:36.38,38.3",
				"github.com/broadwing/deltacoverage/deltacoverage.go:53.22,55.3",
				"github.com/broadwing/deltacoverage/deltacoverage.go:88.50,93.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:96.2,96.29",
				"github.com/broadwing/deltacoverage/deltacoverage.go:131.2,131.43",
				"github.com/broadwing/deltacoverage/deltacoverage.go:136.2,136.12",
				"github.com/broadwing/deltacoverage/deltacoverage.go:96.29,97.51",
				"github.com/broadwing/deltacoverage/deltacoverage.go:97.51,101.18",
				"github.com/broadwing/deltacoverage/deltacoverage.go:104.4,106.23",
				"github.com/broadwing/deltacoverage/deltacoverage.go:126.4,126.40",
				"github.com/broadwing/deltacoverage/deltacoverage.go:106.23,109.19",
				"github.com/broadwing/deltacoverage/deltacoverage.go:113.5,113.26",
				"github.com/broadwing/deltacoverage/deltacoverage.go:116.5,116.29",
				"github.com/broadwing/deltacoverage/deltacoverage.go:119.5,120.12",
				"github.com/broadwing/deltacoverage/deltacoverage.go:123.5,124.70",
				"github.com/broadwing/deltacoverage/deltacoverage.go:113.26,115.6",
				"github.com/broadwing/deltacoverage/deltacoverage.go:116.29,117.14",
				"github.com/broadwing/deltacoverage/deltacoverage.go:120.12,122.6",
				"github.com/broadwing/deltacoverage/deltacoverage.go:131.43,132.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:132.16,134.4",
				"github.com/broadwing/deltacoverage/deltacoverage.go:139.61,141.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:144.2,144.19",
				"github.com/broadwing/deltacoverage/deltacoverage.go:147.2,153.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:156.2,156.21",
			},
			"TestPrintDeltaCoverage_PrintsCorrectPercentDeltaCoverageGivenCoverProfile": {
				"github.com/broadwing/deltacoverage/deltacoverage.go:60.39,62.23",
				"github.com/broadwing/deltacoverage/deltacoverage.go:65.2,65.37",
				"github.com/broadwing/deltacoverage/deltacoverage.go:76.2,77.35",
				"github.com/broadwing/deltacoverage/deltacoverage.go:65.37,67.26",
				"github.com/broadwing/deltacoverage/deltacoverage.go:74.3,74.73",
				"github.com/broadwing/deltacoverage/deltacoverage.go:67.26,69.14",
				"github.com/broadwing/deltacoverage/deltacoverage.go:72.4,72.66",
				"github.com/broadwing/deltacoverage/deltacoverage.go:69.14,70.13",
			},
			"TestPrintDeltaCoverage_PrintsNoTestsFoundGivenDirectoryWithNoCoverProfile": {
				"github.com/broadwing/deltacoverage/deltacoverage.go:60.39,62.23",
				"github.com/broadwing/deltacoverage/deltacoverage.go:62.23,64.3",
				"github.com/broadwing/deltacoverage/deltacoverage.go:88.50,93.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:96.2,96.29",
				"github.com/broadwing/deltacoverage/deltacoverage.go:131.2,131.43",
				"github.com/broadwing/deltacoverage/deltacoverage.go:136.2,136.12",
				"github.com/broadwing/deltacoverage/deltacoverage.go:139.61,141.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:144.2,144.19",
				"github.com/broadwing/deltacoverage/deltacoverage.go:147.2,153.16",
				"github.com/broadwing/deltacoverage/deltacoverage.go:156.2,156.21",
			},
		},
		NumberStatements: 73,
	}
	want := `TestNewCoverProfile_ErrorsIfPathDoesNotExist 1.4%
TestNewCoverProfile_ErrorsIfPathIsNotDirectory 1.4%
TestParseCoverProfile_ReturnsExpectedCoverProfileGivenCoverProfileDirectory 49.3%
TestPrintDeltaCoverage_PrintsCorrectPercentDeltaCoverageGivenCoverProfile 13.7%
TestPrintDeltaCoverage_PrintsNoTestsFoundGivenDirectoryWithNoCoverProfile 1.4%`
	got := c.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(strings.Fields(want), strings.Fields(got)))
	}
}

func TestGenerate_ErrorsGivenCodeWithNoTests(t *testing.T) {
	t.Parallel()
	c, err := deltacoverage.NewCoverProfile(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	output := &bytes.Buffer{}
	c.Stderr = output
	err = c.Generate()
	if err == nil {
		t.Error("want error but got nil")
	}
	want := "go: go.mod file not found in current directory or any parent directory; see 'go help modules'\n"
	got := output.String()
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestGenerate_CreatesExpectedCoverProfilesGivenCodeWithThreeTests(t *testing.T) {
	t.Parallel()
	c, err := deltacoverage.NewCoverProfile("testdata/code/threeTests")
	if err != nil {
		t.Fatal(err)
	}
	tempDir := t.TempDir()
	c.OutputPath = tempDir
	err = c.Generate()
	if err != nil {
		t.Fatal(err)
	}
	wants := map[string][]byte{
		"TestSumOnePlusOne":        []byte("mode: set\ncalc/calc.go:3.24,5.2 1 1\ncalc/calc.go:7.30,9.2 1 0\ncalc/calc.go:11.29,13.2 1 0\n"),
		"TestSumTwoPlusTwo":        []byte("mode: set\ncalc/calc.go:3.24,5.2 1 1\ncalc/calc.go:7.30,9.2 1 0\ncalc/calc.go:11.29,13.2 1 0\n"),
		"TestSubstractTwoMinusTwo": []byte("mode: set\ncalc/calc.go:3.24,5.2 1 0\ncalc/calc.go:7.30,9.2 1 1\ncalc/calc.go:11.29,13.2 1 0\n"),
	}
	for _, testName := range c.Tests {
		f, err := os.Open(filepath.Join(tempDir, testName+".coverprofile"))
		if err != nil {
			t.Fatal(err)
		}
		want, ok := wants[testName]
		if !ok {
			t.Errorf("missing key %q in wants map", testName)
			continue
		}
		got, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(strings.Fields(string(want)), strings.Fields(string(got))))
		}
	}
}

func TestListTests_SetsExpectedTestsNamesGivenCodeWithOneTest(t *testing.T) {
	t.Parallel()
	c, err := deltacoverage.NewCoverProfile("testdata/code/oneTest")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"TestSumOnePlusOne"}
	err = c.ListTests()
	if err != nil {
		t.Fatal(err)
	}
	got := c.Tests
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCleanup_RemovesOutputPath(t *testing.T) {
	t.Parallel()
	c, err := deltacoverage.NewCoverProfile(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tempDir := t.TempDir()
	c.OutputPath = tempDir
	err = c.Cleanup()
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(tempDir)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("want os.ErrNotExist but got %#v", err)
	}
}
