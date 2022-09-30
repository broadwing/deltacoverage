package deltacoverage_test

import (
	"os"
	"testing"

	"github.com/broadwing/deltacoverage"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{"deltacoverage": deltacoverage.Main}))
}

func TestTestScript_(t *testing.T) {
	runIntegrationTests := os.Getenv("DELTACOVERAGE_TESTS_INTEGRATION")
	if runIntegrationTests == "" {
		t.Skip("Set DELTACOVERAGE_TESTS_INTEGRATION=<anything> to run integration tests")
	}
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
	})
}
