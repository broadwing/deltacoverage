package deltacoverage_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/broadwing/deltacoverage"
	"github.com/google/go-cmp/cmp"
)

func TestParseCoverProfile_ErrorsIfFileDoesNotExist(t *testing.T) {
	t.Parallel()
	_, err := deltacoverage.ParseCoverProfile("bogus-filename.txt.bogus")
	if err == nil {
		t.Error("want error not nil")
	}
}
func TestParseCoverProfile_(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc            string
		filePath        string
		extraFilesPaths []string
		want            deltacoverage.CoverProfile
	}{
		{
			desc:     "Unique test returns CoverProfile representing one hundred percent delta coverage",
			filePath: "testdata/sample1/TestSum.coverprofile",
			want: deltacoverage.CoverProfile{
				Branches: map[string]int{
					"xyz/xyz.go:3.24,5.2": 1,
				},
				Tests:           map[string][]string{"TestSum": {"xyz/xyz.go:3.24,5.2"}},
				TotalStatements: 1,
			},
		},
		{
			desc:            "Two tests with same code path returns CoverProfile representing zero percent delta coverage",
			extraFilesPaths: []string{"testdata/sample2/TestSumTwoPlusTwo.coverprofile"},
			filePath:        "testdata/sample2/TestSumOnePlusOne.coverprofile",
			want: deltacoverage.CoverProfile{
				Branches: map[string]int{
					"xyz/xyz.go:3.24,5.2": 1,
				},
				Tests: map[string][]string{
					"TestSumOnePlusOne": {},
					"TestSumTwoPlusTwo": {},
				},
				TotalStatements: 1,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := deltacoverage.ParseCoverProfile(tC.filePath, tC.extraFilesPaths...)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tC.want, got) {
				t.Error(cmp.Diff(tC.want, got))
			}
		})
	}
}

func TestPrintDeltaCoverage_(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc       string
		want       string
		covProfile deltacoverage.CoverProfile
	}{
		{
			covProfile: deltacoverage.CoverProfile{
				Branches: map[string]int{
					"xyz/xyz.go:3.24,5.2": 1,
				},
				Tests: map[string][]string{
					"TestSum": {"xyz/xyz.go:3.24,5.2"},
				},
				TotalStatements: 1,
			},
			desc: "Unique test prints one hundred percent of delta coverage",
			want: "TestSum 100.0%",
		},
		{
			covProfile: deltacoverage.CoverProfile{
				Branches: map[string]int{
					"xyz/xyz.go:3.24,5.2": 1,
				},
				Tests: map[string][]string{
					"TestSumOnePlusOne": {},
					"TestSumTwoPlusTwo": {},
				},
				TotalStatements: 1,
			},
			desc: "Two tests with same code path prints zero percent of delta coverage",
			want: "TestSumOnePlusOne 0.0%\nTestSumTwoPlusTwo 0.0%",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output := &bytes.Buffer{}
			_, err := fmt.Fprint(output, tC.covProfile)
			if err != nil {
				t.Fatal(err)
			}
			got := output.String()
			if !cmp.Equal(tC.want, got) {
				t.Error(cmp.Diff(tC.want, got))
			}
		})
	}
}
