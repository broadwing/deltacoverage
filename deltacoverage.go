package deltacoverage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Main() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please, inform a test name")
		return 1
	}
	testName := os.Args[1]
	coverage, err := GetFullCoverage()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	tests, err := GetTestList()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	testCoverage, err := GetTestCoverage(testName, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Fprintf(os.Stdout, "%s %.1f%s\n", testName, coverage-testCoverage, "%")
	return 0
}

func GetTestList() ([]string, error) {
	cmd := exec.Command("go", "test", "-list", ".")
	cmd.Stderr = os.Stderr
	goTestList, err := cmd.StdoutPipe()
	if err != nil {
		return []string{}, err
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Err starting cmd: %+v\n", err)
		return []string{}, err
	}
	tests, err := ParseTestList(goTestList)
	if err != nil {
		fmt.Printf("Err running ParseTestList: %+v\n", err)
		return []string{}, err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Err waiting cmd: %+v\n", err)
		return []string{}, err
	}
	return tests, nil
}

func GetFullCoverage() (float64, error) {
	cmd := exec.Command("go", "test", "-coverprofile", "/dev/null")
	cmd.Stderr = os.Stderr
	goTestCover, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Err starting cmd: %+v\n", err)
		return 0, err
	}
	coverage, err := ParseCoverageResult(goTestCover)
	if err != nil {
		fmt.Printf("Err running ParseCoverageResult: %+v\n", err)
		return 0, err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Err waiting cmd: %+v\n", err)
		return 0, err
	}
	return coverage, nil
}

func GetTestCoverage(testName string, allTests []string) (float64, error) {
	tests := []string{}
	for _, test := range allTests {
		if test == testName {
			continue
		}
		tests = append(tests, fmt.Sprintf("^%s$", test))
	}
	cmd := exec.Command("go", "test", "-coverprofile", "file.out", "-run", strings.Join(tests, "|"))
	cmd.Stderr = os.Stderr
	goTestCover, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Err starting cmd: %+v\n", err)
		return 0, err
	}
	coverage, err := ParseCoverageResult(goTestCover)
	if err != nil {
		fmt.Printf("Err running ParseCoverageResult: %+v\n", err)
		return 0, err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Err waiting cmd: %+v\n", err)
		return 0, err
	}
	return coverage, nil
}

// ParseCoverageResult returns a non-negative value represeting the code coverage
// in the output result. This logic assumes that total coverage comes after
// normal coverage
func ParseCoverageResult(r io.Reader) (float64, error) {
	scanner := bufio.NewScanner(r)
	var err error
	coverage := -1.0
	for scanner.Scan() {
		items := strings.Fields(scanner.Text())
		switch items[0] {
		case "coverage:":
			coverage, err = strconv.ParseFloat(strings.ReplaceAll(items[1], "%", ""), 64)
			if err != nil {
				return 0, err
			}
		case "total":
			coverage, err = strconv.ParseFloat(strings.ReplaceAll(items[2], "%", ""), 64)
			if err != nil {
				return 0, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	if coverage < 0 {
		return 0, errors.New("coverage not found")
	}
	return coverage, nil
}

func ParseTestList(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	testsNames := []string{}
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "Test") {
			testsNames = append(testsNames, text)
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}
	return testsNames, nil
}
