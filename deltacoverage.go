package deltacoverage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/errgo.v2/errors"
)

var ErrMustBeDirectory = errors.New("Must be a directory")
var ErrMustBeFile = errors.New("Must be a file")

type profileItem struct {
	Branch     string
	Statements int
	Visited    bool
}

type CoverProfile struct {
	NumberStatements int
	OutputPath       string
	PackagePath      string
	Stderr           io.Writer
	Stdout           io.Writer
	Tests            []string
	TestsBranches    map[string][]string
	UniqueBranches   map[string]int
}

// parseProfileLine returns a pointer to ProfileItem type for a given string.
// line example with headers
// branch              statements visited
// xyz/xyz.go:3.24,5.2 1          1
func (c CoverProfile) parseProfileLine(line string) (*profileItem, error) {
	if strings.HasPrefix(line, "mode:") {
		return &profileItem{}, nil
	}
	items := strings.Fields(line)
	branch := items[0]
	profItem := &profileItem{
		Branch: branch,
	}
	nrStmt, err := strconv.Atoi(items[1])
	if err != nil {
		return &profileItem{}, err
	}
	profItem.Statements = nrStmt
	timesVisited, err := strconv.Atoi(items[2])
	if err != nil {
		return &profileItem{}, err
	}
	if timesVisited > 0 {
		profItem.Visited = true
	}
	return profItem, err
}

// String prints out the deltacoverage percentage for each test
func (c CoverProfile) String() string {
	output := []string{}
	if len(c.Tests) == 0 {
		fmt.Fprintf(c.Stderr, "No tests found\n")
		return ""
	}
	for testName, ids := range c.TestsBranches {
		perc := 0.0
		testStatements := 0
		for _, id := range ids {
			statements, ok := c.UniqueBranches[id]
			if !ok {
				continue
			}
			testStatements += statements
		}
		perc = float64(testStatements) / float64(c.NumberStatements) * 100
		output = append(output, fmt.Sprintf("%s %.1f%s", testName, perc, "%"))
	}
	sort.Strings(output)
	return strings.Join(output, "\n")
}

// Parse reads all files from a directory with extension
// .coverprofile, parses it, and populate CoverProfile object
// Design tradeoff to get the total number of statements:
// 1. Parse one file just to sum the total statements
// 2. Recalculate total statements each iteration on the loop
// 3. Add a counter and an if to just calculate the total sum in the first
// iteration
// Current implementation is number 3
func (c *CoverProfile) Parse() error {
	branchesCount := map[string]int{}
	branchesStmts := map[string]int{}
	profilesRead := 0
	files, err := os.ReadDir(c.OutputPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".coverprofile" {
			profilesRead++
			testName := strings.Split(filepath.Base(file.Name()), ".")[0]
			f, err := os.Open(c.OutputPath + "/" + file.Name())
			if err != nil {
				return err
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				profileItem, err := c.parseProfileLine(line)
				if err != nil {
					return fmt.Errorf("cannot parse profile line %q: %+v", line, err)
				}
				// only sum total statements when reading the first cover profile
				if profilesRead == 1 {
					c.NumberStatements += profileItem.Statements
				}
				if !profileItem.Visited {
					continue
				}
				_, ok := branchesCount[profileItem.Branch]
				if !ok {
					branchesStmts[profileItem.Branch] = profileItem.Statements
				}
				branchesCount[profileItem.Branch]++
				c.TestsBranches[testName] = append(c.TestsBranches[testName], profileItem.Branch)
			}
			if err := scanner.Err(); err != nil {
				return err
			}
		}
	}
	for branch, times := range branchesCount {
		if times < 2 {
			c.UniqueBranches[branch] = branchesStmts[branch]
		}
	}
	return nil
}

func (c *CoverProfile) Generate() error {
	err := c.ListTests()
	if err != nil {
		return err
	}
	for _, test := range c.Tests {
		outputFile := c.OutputPath + "/" + test + ".coverprofile"
		cmd := exec.Command("go", "test", "-run", test, "-coverprofile", outputFile)
		cmd.Dir = c.PackagePath
		cmd.Stderr = c.Stderr
		cmd.Stdout = c.Stdout
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting %q: %v", strings.Join(cmd.Args, " "), err)
		}
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("error running %q: %v", strings.Join(cmd.Args, " "), err)
		}
	}
	return nil
}

func (c *CoverProfile) ListTests() error {
	goArgs := []string{"test", "-list", "."}
	cmd := exec.Command("go", goArgs...)
	cmd.Dir = c.PackagePath
	cmd.Stderr = c.Stderr
	goTestList, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting pipe for \"go %s\": %v", strings.Join(goArgs, " "), err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting \"go %s\": %v", strings.Join(goArgs, " "), err)
	}
	c.Tests, err = parseListTests(goTestList)
	if err != nil {
		return fmt.Errorf("error running parseListTests: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error running \"go %s\": %v", strings.Join(goArgs, " "), err)
	}
	return nil
}

func (c *CoverProfile) Cleanup() error {
	return os.RemoveAll(c.OutputPath)
}

func NewCoverProfile(codePath string) (*CoverProfile, error) {
	info, err := os.Stat(codePath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrMustBeDirectory
	}
	tempDir, err := os.MkdirTemp(os.TempDir(), "deltacoverage")
	if err != nil {
		return nil, err
	}
	c := &CoverProfile{
		OutputPath:     tempDir,
		PackagePath:    codePath,
		Stderr:         os.Stderr,
		Stdout:         io.Discard,
		TestsBranches:  map[string][]string{},
		UniqueBranches: map[string]int{},
	}
	return c, nil
}

func parseListTests(r io.Reader) ([]string, error) {
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

func Main() int {
	args := os.Args[1:]
	packagePath := "./"
	if len(args) > 0 {
		packagePath = os.Args[1]
	}
	c, err := NewCoverProfile(packagePath)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	err = c.Generate()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	err = c.Parse()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	output := c.String()
	if output == "" {
		return 1
	}
	fmt.Println(output)
	err = c.Cleanup()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
