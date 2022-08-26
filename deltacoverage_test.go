package deltacoverage

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestParseCoverageResult_ErrorsIfGivenNoTests(t *testing.T) {
	t.Parallel()
	_, err := parseCoverageResult(strings.NewReader(`?       my_not_cool_project    [no test files]`))
	if err == nil {
		t.Error("want error but not found")
	}
}

func TestParseCoverageResult_ReturnsCorrectValueGivenTestScriptTestCoverResult(t *testing.T) {
	t.Parallel()
	want := 34.8
	got, err := parseCoverageResult(strings.NewReader(`FAIL
coverage: 20.2% of statements
total coverage: 34.8% of statements
exit status 1
FAIL    github.com/thiagonache/deltacoverage    0.832s`))
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want test coverage of %.1f got %.1f", want, got)
	}
}

func TestParseCoverageResult_ReturnsCorrectValueGivenNoTestScriptsTestCoverResult(t *testing.T) {
	t.Parallel()
	want := 16.3
	got, err := parseCoverageResult(strings.NewReader(`PASS
coverage: 16.3% of statements
ok      github.com/thiagonache/deltacoverage    0.012s`))
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want test coverage of %.1f got %.1f", want, got)
	}
}

func TestParseListTests_ReturnsZeroTestsIfGivenNoTests(t *testing.T) {
	t.Parallel()
	got, err := parseListTests(strings.NewReader(`?       my_not_cool_project    [no test files]`))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) > 0 {
		t.Error("want zero tests when given no tests")
	}
}

func TestParseListTests_ReturnsCorrectValuesGivenTestListResult(t *testing.T) {
	t.Parallel()
	want := []string{"TestParseCoverageResultReturnsCorrectValueGivenTestCoverResult", "TestParseTestListErrorsIfNoTestsFound"}
	got, err := parseListTests(strings.NewReader(`TestParseCoverageResultReturnsCorrectValueGivenTestCoverResult
TestParseTestListErrorsIfNoTestsFound
ok      github.com/thiagonache/deltacoverage    0.004s`))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{"deltacoverage": Main}))
}

func TestTestScript_(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
	})
}
